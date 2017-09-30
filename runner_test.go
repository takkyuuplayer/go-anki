package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestFindDefinitionForWord(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/put.html")

	initialize()

	if err != nil {
		t.Fatal(err)
	}

	matched := findDefinition(string(data))

	if !strings.HasPrefix(matched, "<h4") {
		t.Errorf(`strings.HasPrefix(matched, "<h4") = %#v, want true`, strings.HasPrefix(matched, "<div"))
	}
}

func TestFindDefinitionForWordHittingOnlyEnglish(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/subtle.html")

	initialize()

	if err != nil {
		t.Fatal(err)
	}

	matched := findDefinition(string(data))

	if !strings.HasPrefix(matched, "<h3") {
		t.Errorf(`strings.HasPrefix(matched, "<h3") = %#v, want true`, strings.HasPrefix(matched, "<h3"))
	}
}

func TestFindDefinitionForIdiom(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/put_up_with.html")

	initialize()

	if err != nil {
		t.Fatal(err)
	}

	matched := findDefinition(string(data))

	if !strings.HasPrefix(matched, "<h3") {
		t.Errorf(`strings.HasPrefix(matched, "<h3") = %#v, want true value`, strings.HasPrefix(matched, "<div"))
	}

	if !strings.HasSuffix(matched, "</ul>") {
		t.Errorf(`string.HasSuffix(matched, "</ul>") = %#v, want true value`, strings.HasSuffix(matched, "</ul>"))
	}
}

func TestFindDefinitionNotFound(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/put_up_on.html")

	initialize()

	if err != nil {
		t.Fatal(err)
	}

	matched := findDefinition(string(data))

	if matched != "Not Found" {
		t.Errorf(`matched = %#v, want %#v`, matched, "Not Found")
	}
	fmt.Println(matched)
}

func TestGetWiktionaryUrl(t *testing.T) {

	if getWiktionaryUrl("put") != "https://en.wiktionary.org/wiki/put" {
		t.Errorf(`getWiktionaryUrl('put') = %#v, want %#v`, getWiktionaryUrl("put"), "https://en.wiktionary.org/wiki/put")
	}

	if getWiktionaryUrl("put up with") != "https://en.wiktionary.org/wiki/put_up_with" {
		t.Errorf(`getWiktionaryUrl("put up with") = %#v, want %#v`, getWiktionaryUrl("put up with"), "https://en.wiktionary.org/wiki/put_up_with")
	}
}
