package wiktionary

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/takkyuuplayer/go-anki"
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

type Client struct {
	HttpClient *http.Client
}

func Init() {
	for idx, val := range ignoreParagraphs {
		ignoreRegexps[idx] = regexp.MustCompile(`(?s)<h[3-5]><span class="mw-headline" id="` + val + `(?:.+?)(<h.>|\z)`)
	}
}

func FindDefinition(html string) string {
	group := defReg.FindStringSubmatch(html)

	if len(group) != 2 {
		return "Not Found"
	}

	definition := group[1]

	for _, reg := range ignoreRegexps {
		definition = reg.ReplaceAllString(definition, "$1")
	}
	for _, reg := range deleteReg {
		definition = reg.ReplaceAllString(definition, "")
	}

	return definition
}

func GetWiktionaryUrl(word string) string {
	return fmt.Sprintf("https://en.wiktionary.org/wiki/%s", strings.Replace(word, " ", "_", -1))
}

func (client *Client) SearchDefinition(ch chan<- *anki.Result, word string) {
	resp, err := client.HttpClient.Get(GetWiktionaryUrl(word))

	if err != nil {
		ch <- &anki.Result{
			Word:      word,
			IsSuccess: false,
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ch <- &anki.Result{
			Word:      word,
			IsSuccess: false,
		}
	}

	definition := FindDefinition(string(body))

	ch <- &anki.Result{
		Word:       word,
		Definition: definition,
		IsSuccess:  definition != "Not Found",
	}
}
