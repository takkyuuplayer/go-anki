package anki

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"strings"

	"github.com/rakyll/statik/fs"
	// https://github.com/rakyll/statik#usage
	_ "github.com/takkyuuplayer/go-anki/anki/statik"
)

//go:generate statik -src assets -f

var tmpl = template.New("entry")

func init() {
	mustParseAssets(tmpl, "/entry.html.tmpl")
}

func mustParseAssets(tmpl *template.Template, path string) *template.Template {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	r, err := statikFS.Open(path)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return template.Must(tmpl.Parse(string(contents)))
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
	Pronunciation *Pronunciation
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

func (entry *Entry) AnkiCard() (string, error) {
	buf := bytes.NewBufferString("")

	if err := tmpl.Execute(buf, entry); err != nil {
		log.Fatalf("execution failed: %s", err)
		return "", err
	}

	return strings.Join(strings.Fields(buf.String()), " "), nil
}
