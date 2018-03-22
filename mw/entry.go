package mw

import (
	"bytes"
	"encoding/xml"
	"log"
	"regexp"
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
	Gram       string     `xml:"gram"`
	Definition Definition `xml:"def"`
}

type Inflection struct {
	InnerXML string `xml:",innerxml"`
}

type Definition struct {
	Gram            string           `xml:"gram"`
	PhrasalVerbForm []string         `xml:"phrasev>pva"`
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

var extractDef = regexp.MustCompile(`<vi>(.+)?</vi>`)

func (dt DefinitionText) Def() string {
	ret := extractDef.ReplaceAllString(dt.InnerXML, "")

	if strings.HasPrefix(ret, ":") {
		return ret[1:]
	}

	return ret
}

var word = template.Must(template.ParseFiles("mw/word.tmpl.html", "mw/definition.tmpl.html"))
var phrase = template.Must(template.ParseFiles("mw/phrase.tmpl.html", "mw/definition.tmpl.html"))

func render(tpl *template.Template, e interface{}) string {
	buf := bytes.NewBufferString("")

	if err := tpl.Execute(buf, e); err != nil {
		log.Fatalf("execution failed: %s", err)
	}

	return buf.String()
}

func (e *Entry) AnkiCard(headWord string) string {

	if strings.Replace(e.HeadWord, "*", "", -1) == headWord {
		return strings.Replace(render(word, e), "\n", "", -1)
	}

	if e.DefinedRunOn != nil {
		for _, dro := range e.DefinedRunOn {
			if dro.Phrase == headWord {
				return strings.Replace(render(phrase, dro), "\n", "", -1)
			}
		}
	}

	return ""
}

func (el *EntryList) AnkiCard(headWord string) string {
	ret := ""

	for _, entry := range el.Entries {
		ret += entry.AnkiCard(headWord)
	}

	return ret
}
