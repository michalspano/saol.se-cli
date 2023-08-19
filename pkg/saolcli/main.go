package saol

import (
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/michalspano/saol.se-cli/types"
	"github.com/michalspano/saol.se-cli/cmd/utils"
)

// A global state variable that is used to filter out the words
// that are not of the desired type.
// TODO: avoid using global state variables.
var (
  found = false
)

func Execute(query string, isID bool, T string) (string, error) {
	// By default, the id is not used, hence the `sok` parameter is used.
	// If the id is used (`isID` is true), then the `id` parameter is used.
	// That is carried in the following if-construct.
	url := fmt.Sprintf("https://svenska.se/saol/?sok=%s", query)
	if isID {
		url = fmt.Sprintf("https://svenska.se/saol/?id=%s", query)
	}

	c := colly.NewCollector() // Declare a collector instance

	// Here we declare the response variable, which will be returned
	// by the function. It is of type Word, which is a struct defined
	// above.
	var response types.Word

	// A div class="lemma" contains all the information about a word,
	// thereby, we can use it as a selector.
	c.OnHTML("div.lemma", func(e *colly.HTMLElement) {
    found = true 
		baseForm := e.ChildText("span.grundform")
		wordType := e.ChildText("a.ordklass")
  
    if T != wordType && T != "" {
      return
    }

		response.BaseForm = baseForm
		response.WordType = wordType

		// The meanings are stored in a span with class="lexemid".
		// The `\u200b` character is a zero-width space, which is used
		// to separate the meanings. For our purposes, we can remove it.
		// Moreover, all the meanings are storred in a buffer array,
		// which is later assigned to the response.
		meanings := make([]string, 0)
		e.DOM.Find("span.lexemid").Each(func(_ int, s *goquery.Selection) {
			meanings = append(meanings, strings.Replace(s.Text(), "\u200b", "", -1))
		})

		response.Meanings = meanings

		// By default, these are set to false
		// The values can be later overwritten
		response.Rules = map[string]bool{
			"plural": false,
			"ovrigt": false,
		}

		// A table of class="tabell" on the SAOL's website contains
		// all the grammatical information about a word. We can use
		// it as a selector.
		table := e.DOM.Find("table.tabell")

		// Each header row contains information, namely the possible grammar
		// rules, that apply to the word. We can use it as a selector.
		// By default, these rules are set to false. In this step we decide
		// whether such fules are applicable for out selected word.
		table.Find("th.ordformth").Each(func(_ int, th *goquery.Selection) {
			if th.Text() == "Plural" {
				response.Rules["plural"] = true
			} else if th.Text() == "Övrig(a) form(er)" {
				response.Rules["ovrigt"] = true
			}
		})

		// The following iterates over the table and assigns the values.
		// This way, we extract the grammatical information about the word.
		// The information (in a string format) are then added to the `grammarBuff`
		// array.
		// However, such a traversal also includes the text labels, therefore,
		// we need to filter them out (this is done by the `switch` statement).
		grammarBuff := make([]string, 0)
		table.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			tr.Find("td").Each(func(_ int, td *goquery.Selection) {
				grammarBuff = append(grammarBuff, td.Text())
			})
		})
   
    // Decide what type of word is added to the `response.TypeDef` field of the `response` struct.
		switch wordType {
		case "substantiv":
			{
				response.TypeDef = utils.CreateNoun(grammarBuff, response.Rules, e.ChildText("span.bojning"))
			}
		default:
			{
				response.TypeDef = nil
			}
		}
	})
  
  // Visit the website and catch any errors.
  // This is the actual execution of the crawler.
	err := c.Visit(url)
	if err != nil {
		return "", err
	}
  
  // In case that the BaseForm is empty, then the word is not found.
  // This can happen for multiple reasons: the word is not in the dictionary,
  // the word can be of more than one type. We use the additional `found` boolean
  // to only access the block of code in the event of the 2nd case (multiple types).
	if response.BaseForm == "" && !found {
		tmpC := colly.NewCollector() // Declare a new temporary collector instance

    // The following code block is used to extract all the possible types of the word.
    // A div class="cshow" contains all the possible types of the word. Then, each a tag
    // holds the type and the hred value the id of the word. The `options` array is used 
    // to store the type and the id of the word.
		options := make([][]string, 0)
		tmpC.OnHTML("div.cshow", func(e *colly.HTMLElement) {
			e.DOM.Find("a").Each(func(_ int, aTag *goquery.Selection) {
        // remove redundant spaces
				wordType := strings.Replace(aTag.Find("span.wordclass").Text(), " ", "", -1)
        // additional remapping; the website (sometimes) uses abbreviations
        if wordType == "subst." {
          wordType = "substantiv"
        }
				wordID := aTag.AttrOr("href", "")
				options = append(options, []string{wordType, wordID})
			})
		})

    // Call the temporary collector instance to visit the website.
    // Handle any errors.
    err := tmpC.Visit(url)
    if err != nil {
      return "", err
    }
    
    // Providing the type is an optional argument. However, if it is provided,
    // a few steps can be skipped. The following block of code checks whether
    // the type (T) belongs to current word's types. If it does, then the function
    // is called recursively with the type as an argument and the stripped ID.
    // Giving the T paramter means skipping the scanning of the standard input,
    // so an advantage in terms of scalability is achieved.
    if (T != "") {
      for i := range options {
        if options[i][0] == T {
		      stripID := strings.Split(options[i][1], "?id=")[1]
          return Execute(stripID, true, T)
        }
      }
    }

    // Print all the possible forms, so that the user can select the desired one. 
		fmt.Printf("The word %s occurs in the following forms:\n", query)
		for i := range options {
			fmt.Printf("  %d. %s: %s\n", i+1, options[i][0], options[i][1])
		}
  
    // Prompt the user to select the desired form.
    // The cycle is only broken when the user enters a valid number.
    // Otherwise, C-z can be used to exit the program.
		var choice int
		fmt.Printf("Select a number between 1 and %d for the desired form.\n", len(options))
		for {
			fmt.Print("~> ") // a prompt cursor
			_, err := fmt.Scanf("%d", &choice)
			if choice < 1 || choice > len(options) || err != nil {
				fmt.Println("Please enter a valid number.")
				continue
			}
			break
		}
  
    // The `chosenID` variable holds the id of the selected word.
    // We need to additionally parse the id from the href value (not the whole sub-url).
    // Moreover, we update the `desiredType` variable, so that the crawler knows what type
    // of word it is dealing with (in the recursive call).
		chosenID    := options[choice-1][1]
		stripID     := strings.Split(chosenID, "?id=")[1]
    desiredType := options[choice-1][0]
    
    // Recursively call the `Execute` function with the new id.
		return Execute(stripID, true, desiredType) 
	}
  
  // Decide what formatting function to use based on the type of the word.
  // and return the string.
	switch response.WordType {
	case "substantiv":
		{
			return utils.FormatNoun(response), nil
		}
  case "verb":
    {
      return "", fmt.Errorf("TODO: add support for verbs.")
    }
  case "adjektiv":
    {
      return "", fmt.Errorf("TODO: add support for adjectives.")
    }
	default:
		{
			return "", fmt.Errorf("The word %s was not found in SAOL.", query)
		}
	}
}

