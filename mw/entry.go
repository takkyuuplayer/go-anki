package mw

import (
	"bytes"
	"encoding/xml"
	"log"
	"strings"
	"text/template"
)

type EntryList struct {
	XMLName     xml.Name `xml:"entry_list"`
	Version     string   `xml:"version,attr"`
	Entries     []Entry  `xml:"entry"`
	Suggestions []string `xml:"suggestion"`
}

type Entry struct {
	XMLName         xml.Name       `xml:"entry"`
	ID              string         `xml:"id,attr"`
	HeadWord        string         `xml:"hw"`
	Inflection      []Inflection   `xml:"in"`
	Pronunciation   string         `xml:"pr"`
	FunctionalLabel string         `xml:"fl"`
	DefinedRunOn    []DefinedRunOn `xml:"dro"`
	Definition      Definition     `xml:"def"`
}

type DefinedRunOn struct {
	Phrase     string     `xml:"dre"`
	Definition Definition `xml:"def"`
}

type Inflection struct {
	InnerXML string `xml:",innerxml"`
}

type Definition struct {
	Gram            string           `xml:"gram"`
	SenseNumber     []string         `xml:"sn"`
	DefinitionTexts []DefinitionText `xml:"dt"`
}

type DefinitionText struct {
	VerbalIllustrations []VerbalIllustration `xml:"vi"`
	Synonyms            []Synonym            `xml:"sx"`
	InnerXML            string               `xml:",innerxml"`
}

type Synonym struct {
	InnerXML string `xml:",innerxml"`
}

type VerbalIllustration struct {
	Text string `xml:",innerxml"`
}

func (dt DefinitionText) Def() string {
	from := strings.Index(dt.InnerXML, ":")
	to := strings.Index(dt.InnerXML, "<")

	if to == -1 {
		return dt.InnerXML[from+1:]
	}

	return dt.InnerXML[from+1 : to]
}

var word = template.Must(template.ParseFiles("mw/word.tmpl.html", "mw/definition.tmpl.html"))
var phrase = template.Must(template.ParseFiles("mw/phrase.tmpl.html", "mw/definition.tmpl.html"))

func (e *Entry) AnkiCard() string {
	buf := bytes.NewBufferString("")

	if err := word.Execute(buf, e); err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	return strings.Replace(buf.String(), "\n", "", -1)
}

func (el *EntryList) AnkiCard(headWord string) string {
	ret := ""

	for _, entry := range el.Entries {
		if len(el.Entries) == 1 {
			ret += entry.AnkiCard()
		} else if strings.Replace(entry.HeadWord, "*", "", -1) == headWord {
			ret += entry.AnkiCard()
		}
	}

	return ret
}
