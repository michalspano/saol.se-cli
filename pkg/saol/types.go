package saol

import "encoding/json"

// Entry is a SAOL article as served by the svenska.se API. Text fields hold
// the API's raw strings, which may contain inline HTML; Format cleans them.
type Entry struct {
	ID            string        `json:"id"`
	Headword      string        `json:"ordled"`     // spelling with syllable marks, e.g. "vatten|kval·itet"
	Homograph     int           `json:"homografNr"` // 0 unless the spelling has homographs
	WordClass     string        `json:"ordklass"`   // e.g. "substantiv"
	EntryType     string        `json:"ingångstyp"` // "lemma", "variant" or "se under"
	MainLemma     string        `json:"huvudlemma"` // link to the lemma of a non-lemma entry
	Inflection    string        `json:"böjning"`    // abbreviated inflection, e.g. "~en ~ar"
	Pronunciation Pronunciation `json:"uttal"`
	PronNote      string        `json:"uttalskommentar"`
	Variants      []Variant     `json:"variantformer"`
	Senses        []Sense       `json:"huvudbetydelser"`
	Compounds     string        `json:"sammansättningskommentar"` // compound forms, e.g. "gatu-"
	SeeAlso       []string      `json:"enbartDigitalaHänvisningar"`
	Tables        []Table       `json:"böjningstabell"`
}

// Variant is an alternative form of the headword, e.g. "el. vattenkoppor".
type Variant struct {
	Kind       string `json:"typ"` // "el.", "hellre än", ...
	Headword   string `json:"ordled"`
	Inflection string `json:"böjning"`
}

// Sense is a main meaning (huvudbetydelse) of an entry.
type Sense struct {
	Definition string      `json:"definition"`
	Usage      string      `json:"bruklighetskommentar"` // e.g. "vard."
	Form       FormNote    `json:"formkommentar"`        // e.g. "mest i pl."
	References []Reference `json:"hänvisningar"`
	Examples   []Example   `json:"exempel"`
}

// Reference points to another entry, e.g. "till tjänst 4".
type Reference struct {
	Kind   string `json:"typ"`
	Target string `json:"hänvisning"`
}

// Example is a usage example with an optional paraphrase.
type Example struct {
	Text       string `json:"text"`
	Paraphrase string `json:"parafras"`
}

// Table is one block of the inflection table: "Singular", "Plural" or
// "Ordform(er)" for indeclinable nouns.
type Table struct {
	Heading string `json:"rubrik"`
	Rows    []Row  `json:"rader"`
}

// Row is one grammatical form, e.g. "en bil — obestämd form". A row holds
// several forms when SAOL lists variants, e.g. "cd:ar (cd)".
type Row struct {
	Forms    []Form `json:"böjningsvarianter"`
	Label    string `json:"postskript"` // e.g. "bestämd form genitiv"
	Addition string `json:"tilläggstext"`
}

// Form is an inflected form with its leading article ("en", "ett" or "").
type Form struct {
	Article string `json:"ledtext"`
	Form    string `json:"böjningsform"`
}

// Pronunciation lists pronunciation forms. The API sends a string, an object
// or a list of objects, each with a "form" field.
type Pronunciation []string

func (p *Pronunciation) UnmarshalJSON(b []byte) error {
	var v any
	err := json.Unmarshal(b, &v)
	*p = pronForms(v)
	return err
}

func pronForms(v any) []string {
	switch t := v.(type) {
	case string:
		return []string{t}
	case map[string]any:
		f, _ := t["form"].(string)
		return []string{f}
	case []any:
		return flatMap(t, pronForms)
	}
	return nil
}

// FormNote is a comment on a sense's form. The API sends either a string or
// an object {"text": "i pl. vanl.", "tillägg": "order"}.
type FormNote struct {
	Text     string
	Addition string
}

func (n *FormNote) UnmarshalJSON(b []byte) error {
	var v any
	err := json.Unmarshal(b, &v)
	switch t := v.(type) {
	case string:
		*n = FormNote{Text: t}
	case map[string]any:
		text, _ := t["text"].(string)
		addition, _ := t["tillägg"].(string)
		*n = FormNote{Text: text, Addition: addition}
	}
	return err
}
