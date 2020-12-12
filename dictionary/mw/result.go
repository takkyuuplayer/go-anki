package mw

import (
	"fmt"
	"github.com/takkyuuplayer/go-anki/dictionary"
	"html/template"
	"strings"
)

const audioUrl = "https://media.merriam-webster.com/audio/prons/en/us/mp3/%s/%s.mp3"

type Suggestion []string
type Entries []Entry

type Entry struct {
	Meta Meta               `json:"meta"`
	Hwi  Hwi                `json:"hwi"`
	Hom  int                `json:"hom"`
	Fl   string             `json:"fl"`
	Lbs  []string           `json:"lbs,omitempty"`
	Ins  Ins                `json:"ins,omitempty"`
	Gram string             `json:"gram,omitempty"`
	Def  DefinitionSections `json:"def"`
	Uros []struct {
		Ure  string          `json:"ure"`
		Prs  Prs             `json:"prs"`
		Fl   string          `json:"fl"`
		Gram string          `json:"gram"`
		Utxt [][]interface{} `json:"utxt"`
	} `json:"uros,omitempty"`
	Dros []struct {
		Drp  string             `json:"drp"`
		Def  DefinitionSections `json:"def"`
		Gram string             `json:"gram,omitempty"`
		Vrs  []struct {
			Vl string `json:"vl"`
			Va string `json:"va"`
		} `json:"vrs,omitempty"`
	} `json:"dros,omitempty"`
	Dxnls    []string `json:"dxnls,omitempty"`
	Shortdef []string `json:"shortdef"`
}

type Meta struct {
	ID      string `json:"id"`
	UUID    string `json:"uuid"`
	Src     string `json:"src"`
	Section string `json:"section"`
	Target  struct {
		Tuuid string `json:"tuuid"`
		Tsrc  string `json:"tsrc"`
	} `json:"target,omitempty"`
	Highlight   string   `json:"highlight,omitempty"`
	Stems       []string `json:"stems"`
	AppShortdef struct {
		Hw  string   `json:"hw"`
		Fl  string   `json:"fl"`
		Def []string `json:"def"`
	} `json:"app-shortdef"`
	Offensive bool `json:"offensive"`
}

type Hwi struct {
	Hw  string `json:"hw"`
	Prs Prs    `json:"prs"`
}

type Prs []struct {
	Ipa   string `json:"ipa"`
	Sound struct {
		Audio Audio `json:"audio"`
	} `json:"sound"`
}

type Audio string

type Ins []struct {
	Il  string `json:"il"`
	If  string `json:"if"`
	Ifc string `json:"ifc"`
	Prs Prs    `json:"prs,omitempty"`
}

type DefinitionSections []struct {
	Sls  []string          `json:"sls,omitempty"`
	Sseq [][][]interface{} `json:"sseq"`
}

func (prs Prs) convert() []dictionary.Accent {
	accents := make([]dictionary.Accent, len(prs))

	for idx, pr := range prs {
		accent := dictionary.Accent{
			AccentLabel: "US",
			Spelling:    pr.Ipa,
			Audio:       pr.Sound.Audio.convert(),
		}
		accents[idx] = accent
	}

	return accents
}

// https://dictionaryapi.com/products/json#sec-2.prs
func (audio Audio) convert() template.URL {
	var subDir string
	if strings.HasPrefix("bix", string(audio)) {
		subDir = "bix"
	} else if strings.HasPrefix("gg", string(audio)) {
		subDir = "gg"
	} else {
		firstLetter := string(audio[:1])
		switch firstLetter {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "_":
			subDir = "number"
		default:
			subDir = firstLetter
		}
	}

	url := fmt.Sprintf(audioUrl, subDir, audio)
	return template.URL(url)
}

func (ins Ins) convert() []dictionary.Inflection {
	inflections := make([]dictionary.Inflection, len(ins))

	for idx, in := range ins {
		inflection := dictionary.Inflection{
			FormLabel:     in.Il,
			InflectedForm: in.If,
			Pronunciation: dictionary.Pronunciation{
				Notation: "IPA",
				Accents:  in.Prs.convert(),
			},
		}
		inflections[idx] = inflection
	}

	return inflections
}

// https://dictionaryapi.com/products/json#sec-2.def
func (sections DefinitionSections) convert() ([]dictionary.Definition, error) {
	var definitions []dictionary.Definition

	for _, section := range sections {
		for _, senses := range section.Sseq {
			for _, sense := range senses {
				if sense[0].(string) == "sense" {
					sense := sense[1].(map[string]interface{})
					definition, _ := convertDefiningText(sense["dt"])
					definitions = append(definitions, definition)
				}
			}
		}
	}

	return definitions, nil
}

func convertDefiningText(dt interface{}) (dictionary.Definition, error) {
	var dicSense string
	var dicExample []string

	for _, tuple := range dt.([]interface{}) {
		tuple := tuple.([]interface{})
		switch tuple[0].(string) {
		case "text":
			if strings.HasPrefix("{bc}", tuple[0].(string)) {
				dicSense += format(tuple[1].(string)) + " "
			}
		case "vis":
			for _, example := range tuple[1].([]interface{}) {
				dicExample = append(dicExample, format(example.(map[string]interface{})["t"].(string)))
			}
		case "uns":
			for _, dt2 := range tuple[1].([]interface{}) {
				res, _ := convertDefiningText(dt2)
				dicExample = append(dicExample, res.Examples...)
			}
		case "wsgram", "bnw", "ri", "snote":
			// Something todo?
		default:
			panic("unknown")
		}
	}

	return dictionary.Definition{Sense: dicSense, Examples: dicExample}, nil
}

func format(text string) string {
	return text
}
