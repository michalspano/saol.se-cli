package utils

// This package contains utility functions used by the saol package.

import (
  "fmt"
  "strings"
  "github.com/michalspano/saol.se-cli/types"
)

/*
 * This is a helper function that is used to format the "övrig(a) form(er)" 
 * section of a word. It takes a slice of strings (an array) and returns a
 * pair (array) of strings in an array (a 2D array). This way, we can assign
 * each form an additional explanation in a reusable way.
 */
func FormatOtherForms(arr []string) [][]string {
	tmp := make([][]string, 0)
	for i := 0; i < len(arr); i += 2 {
		tmp = append(tmp, []string{arr[i], arr[i+1]})
	}
	return tmp
}

/*
 * This function is responsible for creating a `Noun struct` given the specifications
 * and is then assigned to the `Word struct`'s `TypeDef` field. The return value is a pointer
 * to the newly created `Noun struct`. Furthermore, given the `rules` map, the function
 * will also assign the `Plural` and `OtherForms` fields of the `Noun struct` if the
 * rules allow so.
 */
func CreateNoun(grammar []string, rules map[string]bool, prefix string) *types.Noun {
	noun := types.Noun{}

	noun.Suffix           = prefix
	noun.SgIndefinite     = grammar[0]
	noun.SgIndefiniteGen  = grammar[2]
	noun.SgDefinite       = grammar[4]
	noun.SgDefiniteGen    = grammar[6]

	if rules["plural"] {
		noun.PlIndefinite     = grammar[8]
		noun.PlIndefiniteGen  = grammar[10]
		noun.PlDefinite       = grammar[12]
		noun.PlDefiniteGen    = grammar[14]
	}
	if rules["ovrigt"] {
		noun.OtherForms = FormatOtherForms(grammar[16:])
	}

	return &noun
}

/*
 * This function removes all whitespaces and converts the string to lowercase.
 * It is a helper function used to ensure that the word types stay consistent.
 */
func FormatWordType(T string) string {
  tmp := strings.Replace(T, " ", "", -1)
  return strings.ToLower(tmp)
}

/*
 * The `formatNoun` function takes a Word struct (assumed to be a noun) and returns a string.
 * The string contains all the information in a formatted, human-readable
 * way, based on the formatting per SAOL's website.
 */
func FormatNoun(w types.Word) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("%s %s %s\n", w.BaseForm, w.WordType, w.TypeDef.(*types.Noun).Suffix))

  // If there's only one meaning, add a bullet point (instead of a number).
	if len(w.Meanings) == 1 {
		output.WriteString(fmt.Sprintf("•%s\n", w.Meanings[0]))
	} else {
		for i, meaning := range w.Meanings {
			output.WriteString(fmt.Sprintf("%d. %s\n", i+1, meaning))
		}
	}

	output.WriteString("\n")
	output.WriteString("Singular\n")
	output.WriteString(fmt.Sprintf("%s\t obestämd form\n", w.TypeDef.(*types.Noun).SgIndefinite))
	output.WriteString(fmt.Sprintf("%s\t obestämd form genitiv\n", w.TypeDef.(*types.Noun).SgIndefiniteGen))
	output.WriteString(fmt.Sprintf("%s\t bestämd form\n", w.TypeDef.(*types.Noun).SgDefinite))
	output.WriteString(fmt.Sprintf("%s\t bestämd form genitiv\n", w.TypeDef.(*types.Noun).SgDefiniteGen))

	if w.Rules["plural"] {
		output.WriteString("\n")
		output.WriteString("Plural\n")
		output.WriteString(fmt.Sprintf("%s\t obestämd form\n", w.TypeDef.(*types.Noun).PlIndefinite))
		output.WriteString(fmt.Sprintf("%s\t obestämd form genitiv\n", w.TypeDef.(*types.Noun).PlIndefiniteGen))
		output.WriteString(fmt.Sprintf("%s\t bestämd form\n", w.TypeDef.(*types.Noun).PlDefinite))
		output.WriteString(fmt.Sprintf("%s\t bestämd form genitiv\n", w.TypeDef.(*types.Noun).PlDefiniteGen))
	}

	if w.Rules["ovrigt"] {
		output.WriteString("\n")
		output.WriteString("Övrig(a) form(er)\n")

    newLine := "\n"
		for i, form := range w.TypeDef.(*types.Noun).OtherForms {
      // Rule to not add a new line at the last iteration.
      if i == len(w.TypeDef.(*types.Noun).OtherForms) - 1 { 
        newLine = "" 
      } 
			output.WriteString(fmt.Sprintf("%s\t %s%s", form[0], form[1], newLine))
		}
	}

	return output.String()
}

