package mw_test

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takkyuuplayer/go-anki/mw"
)

func TestEntryListToParseDefinition(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/merriam-webster/test.xml")

	ret := mw.EntryList{}
	xml.Unmarshal([]byte(data), &ret)

	assert.Equal(t, "1.0", ret.Version)
	assert.Equal(t, 10, len(ret.Entries))

	entry := ret.Entries[0]

	assert.Equal(t, "test[1]", entry.ID)
	assert.Equal(t, "test", entry.HeadWord)
	assert.Equal(t, "ˈtɛst", entry.Pronunciation)
	assert.Equal(t, "noun", entry.FunctionalLabel)

	assert.Equal(t, "count", entry.Definition.Gram)
	assert.Equal(t, 6, len(entry.Definition.DefinitionTexts))

	definitionText := entry.Definition.DefinitionTexts[0]

	assert.Equal(t, 9, len(definitionText.VerbalIllustrations))
	assert.Equal(t,
		"She is studying for her math/spelling/history <it>test</it>.",
		string(definitionText.VerbalIllustrations[0].Text),
	)
	assert.NotNil(t, definitionText.InnerXML)

	senseNumber := entry.Definition.SenseNumber

	assert.Equal(t, 6, len(senseNumber))
	assert.Equal(t, "1", senseNumber[0])
}

func TestEntryListToParseSuggestion(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/merriam-webster/furnitura.xml")

	ret := mw.EntryList{}
	xml.Unmarshal([]byte(data), &ret)

	assert.Nil(t, ret.Entries)
	assert.Equal(t, 2, len(ret.Suggestions))
	assert.Equal(t, "furniture", ret.Suggestions[0])
	assert.Equal(t, "fornicator", ret.Suggestions[1])
}

func TestDefinitionText(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/merriam-webster/test.xml")

	ret := mw.EntryList{}
	xml.Unmarshal([]byte(data), &ret)
	entry := ret.Entries[0]
	definitionText := entry.Definition.DefinitionTexts[0]

	assert.Equal(t, "a set of questions or problems that are designed to measure a person's knowledge, skills, or abilities ", definitionText.Def())
}

func TestEntryToAnkiCard(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/merriam-webster/test.xml")
	ret := mw.EntryList{}
	xml.Unmarshal([]byte(data), &ret)

	assert.Equal(t, true, strings.Contains(ret.Entries[0].AnkiCard(), "<h3>noun</h3>"))
}

func TestEntryListToAnkiCard(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/merriam-webster/test.xml")
	ret := mw.EntryList{}
	xml.Unmarshal([]byte(data), &ret)

	fmt.Println(ret.AnkiCard("test"))

	assert.Equal(t, false, strings.Contains(ret.AnkiCard("test"), "<h3>noun</h3>"))
	assert.Equal(t, false, strings.Contains(ret.AnkiCard("wrong"), "<h3>noun</h3>"))
}