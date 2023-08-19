package types

/* The following file contains the definitions via
 * structs for the different types of words (given by
 * SAOL) */

/* By default, every word has the following:
 * Grundform: the word itself (in its basic form)
 * Ordklass: the type of word (e.g. noun, verb, etc.)
 * Meanings: an array of strings, each string representing a meaning
 * Rules: a map of rules that apply to the word (does it have a plural
 *        form? does it have other forms?)
 * Type: further grammatical information about the word's type (e.g.
 *        noun: definite/indefinite forms, etc.)
 */
type Word struct {
	BaseForm string          // Grundform
	WordType string          // Ordklass
	Meanings []string        // Betydelser
	Rules    map[string]bool // Grammatiska regler
	TypeDef  interface{}     // Ordklass-specifik information
}

// Sg = singular, Pl = plural
type Noun struct {
  Suffix            string      // Böjningsändelse
  SgIndefinite      string      // Singular, obestämd form
  SgIndefiniteGen   string      // Singular, obestämd form genitiv
  SgDefinite        string      // Singular, bestämd form
  SgDefiniteGen     string      // Singular, bestämd form genitiv
  PlIndefinite      string      // Plural, obestämd form
  PlIndefiniteGen   string      // Plural, obestämd form genitiv
  PlDefinite        string      // Plural, bestämd form
  PlDefiniteGen     string      // Plural, bestämd form genitiv
  OtherForms        [][]string  // Övriga former
}

// TODO: add more definitions of grammar
