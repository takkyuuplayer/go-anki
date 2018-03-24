package wiktionary_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takkyuuplayer/go-anki/wiktionary"
)

var dictionary = wiktionary.New()

func TestFindDefinitionForWord(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/wiktionary/put.html")
	definition, err := dictionary.AnkiCard(string(data), "put")

	assert.Equal(t, true, strings.HasPrefix(definition, "<h4"))
	assert.Nil(t, err)
}

func TestFindDefinitionForWordHittingOnlyEnglish(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/wiktionary/subtle.html")
	definition, err := dictionary.AnkiCard(string(data), "subtle")

	assert.Equal(t, true, strings.HasPrefix(definition, "<h3"))
	assert.Nil(t, err)
}

func TestFindDefinitionForIdiom(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/wiktionary/put_up_with.html")

	definition, err := dictionary.AnkiCard(string(data), "put up on")

	assert.Equal(t, true, strings.HasPrefix(definition, "<h3"))
	assert.Equal(t, true, strings.HasSuffix(definition, "</ul>"))
	assert.Nil(t, err)
}

func TestFindDefinitionNotFound(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/wiktionary/put_up_on.html")

	definition, err := dictionary.AnkiCard(string(data), "put up on")

	assert.Equal(t, "", definition)
	assert.NotNil(t, err)
}

func TestGetSearchUrl(t *testing.T) {
	assert.Equal(t, dictionary.GetSearchUrl("put"), "https://en.wiktionary.org/wiki/put")
	assert.Equal(t, dictionary.GetSearchUrl("put up with"), "https://en.wiktionary.org/wiki/put_up_with")
}
