package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestFindDefinitionForWord(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/put.html")

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

func TestFindDefinitionForIdiom(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/put_up_with.html")

	if err != nil {
		t.Fatal(err)
	}

	matched := findDefinition(string(data))

	fmt.Println(matched)

	if !strings.HasPrefix(matched, "<h3") {
		t.Errorf(`strings.HasPrefix(matched, "<h3") = %#v, want true value`, strings.HasPrefix(matched, "<div"))
	}

	if !strings.HasSuffix(matched, "</ul>") {
		t.Errorf(`string.HasSuffix(matched, "</ul>") = %#v, want true value`, strings.HasSuffix(matched, "</ul>"))
	}
}

func TestGetWiktionaryUrl(t *testing.T) {

	if getWiktionaryUrl("put") != "http://en.wiktionary.org/wiki/put" {
		t.Errorf(`getWiktionaryUrl('put') = %#v, want %#v`, getWiktionaryUrl("put"), "http://en.wiktionary.org/wiki/put")
	}

	if getWiktionaryUrl("put up with") != "http://en.wiktionary.org/wiki/put_up_with" {
		t.Errorf(`getWiktionaryUrl("put up with") = %#v, want %#v`, getWiktionaryUrl("put up with"), "http://en.wiktionary.org/wiki/put_up_with")
	}
}
