package saol

import (
	"strings"
	"testing"
)

// words lists sample queries per SAOL word class.
var words = map[string][]string{
	"substantiv": {
		"bil", "hus", "bok", "katt", "hund", "äpple", "öga", "gata", "flicka",
		"pojke", "stad", "land", "vatten", "mjölk", "hjärta", "man", "mus",
		"fot", "gås", "museum", "kaktus", "nyckel", "fröken", "sax", "lag",
	},
	"verb":     nil,
	"adjektiv": nil,
}

// TestLookup checks that each word appears in the output with its word class.
func TestLookup(t *testing.T) {
	for class, queries := range words {
		t.Run(class, func(t *testing.T) {
			if class != "substantiv" {
				t.Skip("word class not supported yet")
			}
			for _, q := range queries {
				t.Run(q, func(t *testing.T) {
					entries, err := Lookup(q)
					if err != nil {
						t.Fatal(err)
					}
					if out := Format(entries); !strings.Contains(out, q) || !strings.Contains(out, class) {
						t.Errorf("output lacks %q or %q:\n%s", q, class, out)
					}
				})
			}
		})
	}
}
