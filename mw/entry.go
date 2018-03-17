package mw

import "encoding/xml"

type EntryList struct {
	XMLName     xml.Name `xml:"entry_list"`
	Version     string   `xml:"version,attr"`
	Entries     []Entry  `xml:"entry"`
	Suggestions []string `xml:"suggestion"`
}

type Entry struct {
	XMLName         xml.Name   `xml:"entry"`
	ID              string     `xml:"id,attr"`
	Pronunciation   string     `xml:"pr"`
	FunctionalLabel string     `xml:"fl"`
	Definition      Definition `xml:"def"`
}

type Definition struct {
	Gram            string           `xml:"gram"`
	DefinitionTexts []DefinitionText `xml:"dt"`
}

type DefinitionText struct {
	VerbalIllustration []string `xml:"vi"`
	InnerXML           string   `xml:",innerxml"`
}
