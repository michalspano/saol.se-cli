// Command saol-cli prints the SAOL entries of a Swedish noun.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/michalspano/saol.se-cli/pkg/saol"
)

func main() {
	query := flag.String("query", "", "the noun to look up in SAOL (required)")
	flag.Parse()
	if *query == "" {
		flag.Usage()
		os.Exit(2)
	}

	entries, err := saol.Lookup(*query)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(saol.Format(entries))
}
