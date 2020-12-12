package dictionary

import (
	"errors"
	"html/template"
)

var NotFoundError = errors.New("Not Found")

type Result struct {
	SearchWord  string
	Entries     []Entry
	Suggestions []string
}

type Dictionary interface {
	LookUp(string) (body string, err error)
	Parse(searchWord, body string) (*Result, error)
}

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
