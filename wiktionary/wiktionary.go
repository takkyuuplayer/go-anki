package wiktionary

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var defReg = regexp.MustCompile(`(?s)<h2><span class="mw-headline" id="English">English</span>.*?</h2>.*?(<h3>.+?)\n+(?:<hr />|<!-- \nNewPP limit report)`)
var ignoreParagraphs = []string{
	"Etymology",
	"Derived_terms",
	"Translations",
	"Further_reading",
	"References",
	"Anagrams",
	"Descendants",
	"Hyponyms",
}
var deleteReg = []*regexp.Regexp{
	regexp.MustCompile(`<span class="mw-editsection">.+</span></span>`),
	regexp.MustCompile(`</?a ?[^>]*>`),
}
var ignoreRegexps = make([]*regexp.Regexp, len(ignoreParagraphs))

func init() {
	for idx, val := range ignoreParagraphs {
		ignoreRegexps[idx] = regexp.MustCompile(`(?s)<h[3-5]><span class="mw-headline" id="` + val + `(?:.+?)(<h.>|\z)`)
	}
}

type Wiktionary struct{}

func New() *Wiktionary {
	return &Wiktionary{}
}

func (w *Wiktionary) GetSearchUrl(word string) string {
	return fmt.Sprintf("https://en.wiktionary.org/wiki/%s", strings.Replace(word, " ", "_", -1))
}

func (w *Wiktionary) AnkiCard(body, word string) (string, error) {
	group := defReg.FindStringSubmatch(body)

	if len(group) != 2 {
		return "", errors.New("Not Found")
	}

	definition := group[1]

	for _, reg := range ignoreRegexps {
		definition = reg.ReplaceAllString(definition, "$1")
	}
	for _, reg := range deleteReg {
		definition = reg.ReplaceAllString(definition, "")
	}

	return definition, nil
}
