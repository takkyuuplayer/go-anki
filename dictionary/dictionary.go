package dictionary

import (
	"errors"
	"html/template"
)

// ErrNotFound means there is no looked up result
var ErrNotFound = errors.New("Not Found")

// Result is the looked up result
type Result struct {
	Dictionary  string
	WebUrl      template.URL
	Entries     []Entry
	Suggestions []string
}

// Dictionary is dictionary API's client interface
type Dictionary interface {
	LookUp(string) (body string, err error)
	Parse(searchWord, body string) (*Result, error)
}

// Entry is an entry of looked up result
type Entry struct {
	ID              string
	Headword        string
	FunctionalLabel string
	Pronunciation   *Pronunciation
	Inflections     []Inflection
	Definitions     []Definition
}

// Definition is a meaning of the entry
type Definition struct {
	Sense    template.HTML
	Examples []template.HTML
}

// Inflection is inflection of the entry
type Inflection struct {
	FormLabel     string
	InflectedForm string
	Pronunciation *Pronunciation
}

// Pronunciation is inflection of the entry
type Pronunciation struct {
	Notation string
	Accents  []Accent
}

// Accent is a variant of the pronunciations
type Accent struct {
	AccentLabel string
	Spelling    string
	Audio       template.URL
}
