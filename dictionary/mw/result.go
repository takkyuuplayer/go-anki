package mw

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/takkyuuplayer/go-anki/dictionary"
)

const audioURL = "https://media.merriam-webster.com/audio/prons/en/us/mp3/%s/%s.mp3"

// suggestion suggestive alternative search words
type suggestion []string

// entries are the entries of looked up result
type entries []entry

type entry struct {
	Meta meta               `json:"meta"`
	Hwi  hwi                `json:"hwi"`
	Hom  int                `json:"hom"`
	Fl   string             `json:"fl"`
	Lbs  []string           `json:"lbs,omitempty"`
	Ins  ins                `json:"ins,omitempty"`
	Gram string             `json:"gram,omitempty"`
	Def  definitionSections `json:"def"`
	Uros []struct {
		Ure  hw            `json:"ure"`
		Prs  prs           `json:"prs"`
		Fl   string        `json:"fl"`
		Ins  ins           `json:"ins"`
		Gram string        `json:"gram"`
		Utxt []interface{} `json:"utxt"`
	} `json:"uros,omitempty"`
	Dros []struct {
		Drp  string             `json:"drp"`
		Def  definitionSections `json:"def"`
		Gram string             `json:"gram,omitempty"`
		Vrs  []struct {
			Vl string `json:"vl"`
			Va string `json:"va"`
		} `json:"vrs,omitempty"`
	} `json:"dros,omitempty"`
	Dxnls    []string `json:"dxnls,omitempty"`
	Shortdef []string `json:"shortdef"`
}

type meta struct {
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
		Hw  hw       `json:"hw"`
		Fl  string   `json:"fl"`
		Def []string `json:"def"`
	} `json:"app-shortdef"`
	Offensive bool `json:"offensive"`
}

type hwi struct {
	Hw  hw  `json:"hw"`
	Prs prs `json:"prs"`
}

type prs []struct {
	Ipa   string `json:"ipa"`
	Sound struct {
		Audio audio `json:"audio"`
	} `json:"sound"`
}

type hw string

type audio string

// ins is Inflections
type ins []struct {
	Il  string `json:"il"`
	If  string `json:"if"`
	Ifc string `json:"ifc"`
	Prs prs    `json:"prs,omitempty"`
}

type definitionSections []struct {
	Sls  []string          `json:"sls,omitempty"`
	Sseq [][][]interface{} `json:"sseq"`
}

func (prs prs) convert() []dictionary.Accent {
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
func (audio audio) convert() template.URL {
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

	url := fmt.Sprintf(audioURL, subDir, audio)
	return template.URL(url)
}

func (ins ins) convert() []dictionary.Inflection {
	inflections := make([]dictionary.Inflection, len(ins))

	for idx, in := range ins {
		var pronunciation *dictionary.Pronunciation
		if len(in.Prs) > 0 {
			pronunciation = &dictionary.Pronunciation{
				Notation: "IPA",
				Accents:  in.Prs.convert(),
			}
		}
		inflection := dictionary.Inflection{
			FormLabel:     in.Il,
			InflectedForm: in.If,
			Pronunciation: pronunciation,
		}
		inflections[idx] = inflection
	}

	return inflections
}

// https://dictionaryapi.com/products/json#sec-2.def
func (sections definitionSections) convert() ([]dictionary.Definition, error) {
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
	var dicSenses []string
	var dicExample []template.HTML

	for _, tuple := range dt.([]interface{}) {
		tuple := tuple.([]interface{})
		switch tuple[0].(string) {
		case "text":
			txt := strings.Trim(tuple[1].(string), " ")
			if strings.HasPrefix(txt, "{bc}") {
				dicSenses = append(dicSenses, Format(txt))
			}
		case "vis":
			for _, example := range tuple[1].([]interface{}) {
				dicExample = append(dicExample, template.HTML(Format(example.(map[string]interface{})["t"].(string))))
			}
		case "uns":
			for _, dt2 := range tuple[1].([]interface{}) {
				res, _ := convertDefiningText(dt2)
				dicExample = append(dicExample, res.Examples...)
			}
		case "snote":
			if len(dicSenses) == 0 {
				dicSenses = append(dicSenses, Format(tuple[1].([]interface{})[0].([]interface{})[1].(string)))
			}
			res, _ := convertDefiningText(tuple[1].([]interface{})[1:])
			dicExample = append(dicExample, res.Examples...)
		case "wsgram", "bnw", "ri":
			// Something todo?
		default:
			panic("unknown")
		}
	}

	return dictionary.Definition{Sense: template.HTML(strings.Join(dicSenses, " / ")), Examples: dicExample}, nil
}

var formatter = regexp.MustCompile("{.+}(.+)?{/.+}")

// Format markup text with html
func Format(text string) string {
	text = strings.ReplaceAll(text, "{bc}", "")
	text = formatter.ReplaceAllString(text, "<i>$1</i>")
	return strings.Trim(text, " ")
}

func (hw hw) clean() string {
	return strings.ReplaceAll(string(hw), "*", "")
}
