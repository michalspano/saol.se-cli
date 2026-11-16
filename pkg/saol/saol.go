// Package saol looks up nouns in Svenska Akademiens ordlista (SAOL) through
// the JSON API behind svenska.se. The site renders articles client-side, so
// its HTML holds no entry data to scrape.
package saol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const baseURL = "https://svenska.se/api"

var (
	client  = &http.Client{Timeout: 10 * time.Second}
	lemmaID = regexp.MustCompile(`\bid=(?:lnr)?([\w-]+)`) // in links such as "?q=hus&id=lnr34396"
)

type hit struct {
	Source Entry `json:"_source"`
}

// Lookup returns the SAOL noun entries matching query, in SAOL's order.
// Inflected forms match their lemma ("byxor" finds "byxa"), and variant or
// "se under" entries resolve to their main lemma, as on svenska.se.
func Lookup(query string) ([]Entry, error) {
	entries, err := search(query)
	if err != nil {
		return nil, err
	}
	nouns := filter(entries, func(e Entry) bool { return e.WordClass == "substantiv" })
	if len(nouns) == 0 {
		return nil, notFound(query, entries)
	}
	lemmas, err := mapErr(nouns, resolve)
	return uniqueBy(lemmas, func(e Entry) string { return e.ID }), err
}

func notFound(query string, entries []Entry) error {
	if len(entries) == 0 {
		return fmt.Errorf("%q was not found in SAOL", query)
	}
	classes := uniqueBy(mapf(entries, func(e Entry) string { return e.WordClass }), func(s string) string { return s })
	return fmt.Errorf("%q is not a noun in SAOL (found: %s)", query, join(", ", classes...))
}

// resolve returns the main lemma of a variant or "se under" entry, or e itself.
func resolve(e Entry) (Entry, error) {
	m := lemmaID.FindStringSubmatch(e.MainLemma)
	if e.EntryType == "lemma" || m == nil {
		return e, nil
	}
	return article(m[1])
}

func search(query string) ([]Entry, error) {
	var r struct {
		Hits struct {
			Hits []hit `json:"hits"`
		} `json:"hits"`
	}
	err := get("/search/saol?"+url.Values{"q": {query}, "exactMatch": {"true"}}.Encode(), &r)
	return mapf(r.Hits.Hits, func(h hit) Entry { return h.Source }), err
}

func article(id string) (Entry, error) {
	var h hit
	err := get("/article/saol/"+url.PathEscape(id), &h)
	return h.Source, err
}

func get(path string, v any) error {
	resp, err := client.Get(baseURL + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("saol: GET %s: %s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
