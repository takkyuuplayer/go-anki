package anki

import (
	"bytes"
	"html/template"
	"log"
	"strings"
)

type Entry struct {
	ID              string
	Headword        string
	FunctionalLabel string
	Pronunciation   Pronunciation
	Inflections     []Inflection
	Definitions     []Definition
}

type Definition struct {
	Sense    string
	Examples []string
}

type Inflection struct {
	FormLabel     string
	InflectedForm string
	Pronunciation Pronunciation
}

type Pronunciation struct {
	Notation string
	Accents  []Accent
}

type Accent struct {
	AccentLabel string
	Spelling    string
	Audio       template.URL
}

func (entry Entry) AnkiCard() (string, error) {
	buf := bytes.NewBufferString("")

	if err := tmpl.Lookup("entry").Execute(buf, entry); err != nil {
		log.Fatalf("execution failed: %s", err)
		return "", err
	}

	return strings.Join(strings.Fields(buf.String()), " "), nil
}
