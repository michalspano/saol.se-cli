package saol

import (
	"fmt"
	"html"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
)

var (
	supTag  = regexp.MustCompile(`<sup>\s*(\d+)\s*</sup>`)
	htmlTag = regexp.MustCompile(`<[^>]*>`)
)

// Format renders entries as plain text, in the layout used by svenska.se.
func Format(entries []Entry) string {
	return align(join("\n\n", mapf(entries, entry)...))
}

func entry(e Entry) string {
	article := slices.Concat(
		[]string{header(e)},
		senses(e.Senses),
		[]string{clean(e.Compounds), wrap("se även: ", join(", ", mapf(e.SeeAlso, clean)...), "")},
	)
	return join("\n\n", join("\n", article...), join("\n\n", mapf(e.Tables, table)...))
}

// header renders e.g. "cd [se`de långt sista e] substantiv cd:n; pl. cd:ar el. ~".
func header(e Entry) string {
	return join(" ",
		homograph(e.Homograph)+e.Headword,
		wrap("[", join(" ", join(", ", mapf(e.Pronunciation, clean)...), clean(e.PronNote)), "]"),
		e.WordClass,
		clean(e.Inflection),
		join(" ", mapf(e.Variants, func(v Variant) string {
			return join(" ", v.Kind, clean(v.Headword), clean(v.Inflection))
		})...),
	)
}

// senses numbers the meanings when there is more than one.
func senses(ss []Sense) []string {
	glosses := filter(mapf(ss, sense), nonEmpty)
	if len(glosses) < 2 {
		return glosses
	}
	return mapi(glosses, func(i int, g string) string { return fmt.Sprintf("%d. %s", i+1, g) })
}

// sense renders e.g. "⟨vard.⟩ i pl. vanl. order; militär befallning: ge order".
func sense(s Sense) string {
	gloss := join(" ",
		wrap("⟨", clean(s.Usage), "⟩"),
		join("; ", clean(join(" ", s.Form.Text, s.Form.Addition)), clean(s.Definition)),
		join("; ", mapf(s.References, func(r Reference) string { return join(" ", r.Kind, clean(r.Target)) })...),
	)
	examples := join("; ", mapf(s.Examples, func(x Example) string { return join(" ", clean(x.Text), clean(x.Paraphrase)) })...)
	return join(": ", gloss, examples)
}

func table(t Table) string {
	return join("\n", slices.Concat([]string{t.Heading}, mapf(t.Rows, row))...)
}

// row renders e.g. "en cd:ar (cd)<TAB>obestämd form"; align pads the tab.
func row(r Row) string {
	forms := mapf(r.Forms, func(f Form) string { return join(" ", f.Article, f.Form) })
	return join("\t", join(" ", variants(forms), clean(r.Addition)), r.Label)
}

// variants renders the main form followed by alternatives: "cd:ar (cd)".
func variants(forms []string) string {
	if len(forms) == 0 {
		return ""
	}
	return join(" ", forms[0], wrap("(", strings.Join(forms[1:], ", "), ")"))
}

func homograph(n int) string {
	if n == 0 {
		return ""
	}
	return superscript(strconv.Itoa(n))
}

func superscript(digits string) string {
	sup := []rune("⁰¹²³⁴⁵⁶⁷⁸⁹")
	return strings.Map(func(r rune) rune {
		if '0' <= r && r <= '9' {
			return sup[r-'0']
		}
		return r
	}, digits)
}

// clean converts SAOL's inline HTML to plain text.
func clean(s string) string {
	s = supTag.ReplaceAllStringFunc(s, func(m string) string { return superscript(supTag.FindStringSubmatch(m)[1]) })
	s = html.UnescapeString(htmlTag.ReplaceAllString(s, ""))
	return strings.Join(strings.Fields(s), " ")
}

// align pads tab-separated columns with spaces.
func align(s string) string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	fmt.Fprint(w, s)
	w.Flush()
	return b.String()
}
