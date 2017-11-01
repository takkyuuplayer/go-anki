package wiktionary_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/takkyuuplayer/go-anki/wiktionary"
)

func TestFindDefinitionForWord(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/wiktionary/put.html")

	if err != nil {
		t.Fatal(err)
	}

	matched := wiktionary.FindDefinition(string(data))

	if !strings.HasPrefix(matched, "<h4") {
		t.Errorf(`strings.HasPrefix(matched, "<h4") = %#v, want true`, strings.HasPrefix(matched, "<div"))
	}
}

func TestFindDefinitionForWordHittingOnlyEnglish(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/wiktionary/subtle.html")

	if err != nil {
		t.Fatal(err)
	}

	matched := wiktionary.FindDefinition(string(data))

	if !strings.HasPrefix(matched, "<h3") {
		t.Errorf(`strings.HasPrefix(matched, "<h3") = %#v, want true`, strings.HasPrefix(matched, "<h3"))
	}
}

func TestFindDefinitionForIdiom(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/wiktionary/put_up_with.html")

	if err != nil {
		t.Fatal(err)
	}

	matched := wiktionary.FindDefinition(string(data))

	if !strings.HasPrefix(matched, "<h3") {
		t.Errorf(`strings.HasPrefix(matched, "<h3") = %#v, want true value`, strings.HasPrefix(matched, "<div"))
	}

	if !strings.HasSuffix(matched, "</ul>") {
		t.Errorf(`string.HasSuffix(matched, "</ul>") = %#v, want true value`, strings.HasSuffix(matched, "</ul>"))
	}
}

func TestFindDefinitionNotFound(t *testing.T) {
	data, err := ioutil.ReadFile("../testdata/wiktionary/put_up_on.html")

	if err != nil {
		t.Fatal(err)
	}

	matched := wiktionary.FindDefinition(string(data))

	if matched != "Not Found" {
		t.Errorf(`matched = %#v, want %#v`, matched, "Not Found")
	}
}

func TestGetWiktionaryUrl(t *testing.T) {

	if wiktionary.GetWiktionaryUrl("put") != "https://en.wiktionary.org/wiki/put" {
		t.Errorf(`GetWiktionaryUrl('put') = %#v, want %#v`, wiktionary.GetWiktionaryUrl("put"), "https://en.wiktionary.org/wiki/put")
	}

	if wiktionary.GetWiktionaryUrl("put up with") != "https://en.wiktionary.org/wiki/put_up_with" {
		t.Errorf(`GetWiktionaryUrl("put up with") = %#v, want %#v`, wiktionary.GetWiktionaryUrl("put up with"), "https://en.wiktionary.org/wiki/put_up_with")
	}
}
