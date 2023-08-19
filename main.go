package main

import (
	"flag"
  "fmt"
  "github.com/michalspano/saol.se-cli/pkg/saolcli"
  "github.com/michalspano/saol.se-cli/cmd/utils"
)

// An example entry point of the SAOL package.
func main() {
	query := flag.String("query", "", "A word passed to SAOL [required]")
  isID  := flag.Bool("id", false, "Whether to use the ID of the word [default: false]")
  wordType := flag.String("type", "", "The type of the word [optional]")
	flag.Parse()

	if *query == "" {
		panic("No word provided.")
	}

  // Do some formatting of the wordType.
  // Remove whitespaces, turn to lowercase.
  fwordType := utils.FormatWordType(*wordType) 
  
  // Call the Execute function from the saol package.
	result, err := saol.Execute(*query, *isID, fwordType)
	if err != nil {
		fmt.Println(err)
		return
	}
  
  // Display the result
	fmt.Println(result)
}
