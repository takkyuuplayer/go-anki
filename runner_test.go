package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestFindDefinition(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/result.txt")

	if err != nil {
		t.Fatal(err)
	}

	matched := findDefinition(string(data))

	if !strings.HasPrefix(matched, "<div") {
		t.Errorf(`strings.HasPrefix(matched, "<div") = %#v, want true value`, strings.HasPrefix(matched, "<div"))
	}

	if !strings.HasSuffix(matched, "</ul>") {
		t.Errorf(`string.HasSuffix(matched, "</ul>") = %#v, want true value`, strings.HasSuffix(matched, "</ul>"))
	}
}
