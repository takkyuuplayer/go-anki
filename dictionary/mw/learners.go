package mw

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/takkyuuplayer/go-anki/dictionary"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const searchUrl = "https://www.dictionaryapi.com/api/v3/references/learners/json/%s?key=%s"

// learners is a client to access MERRIAM-WEBSTER'S LEARNER'S DICTIONARY API
// https://dictionaryapi.com/products/api-learners-dictionary
type learners struct {
	apiKey     string
	httpClient *http.Client
}

func NewLearners(apiKey string, httpClient *http.Client) *learners {
	return &learners{apiKey: apiKey, httpClient: httpClient}
}

func (dic *learners) LookUp(word string) (string, error) {
	urlToSearch := fmt.Sprintf(searchUrl, url.PathEscape(word), url.PathEscape(dic.apiKey))

	response, err := dic.httpClient.Get(urlToSearch)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	_ = response.Body.Close()

	bodyText := string(body)
	if response.StatusCode != http.StatusOK {
		return "", errors.New(bodyText)
	}

	if bodyText == "[]" {
		return "", dictionary.NotFoundError
	}

	return bodyText, nil
}

func (dic *learners) Parse(searchWord, body string) (*dictionary.Result, error) {
	var suggestion Suggestion
	err := json.Unmarshal([]byte(body), &suggestion)
	if err == nil {
		return &dictionary.Result{SearchWord: searchWord, Entries: nil, Suggestions: suggestion}, nil
	}

	var entries Entries
	err = json.Unmarshal([]byte(body), &entries)
	if err != nil {
		return nil, fmt.Errorf("unknown structure: %v", err)
	}

	isPhrasalVerb := len(strings.Fields(searchWord)) > 1

	var dictEntries []dictionary.Entry
	for _, entry := range entries {
		var lookedUp []dictionary.Entry
		if isPhrasalVerb {
			lookedUp = lookUpForPhrase(searchWord, entry)
		} else {
			lookedUp = lookUpForWord(searchWord, entry)
		}
		if lookedUp != nil {
			dictEntries = append(dictEntries, lookedUp...)
		}
	}

	if len(dictEntries) == 0 {
		return nil, dictionary.NotFoundError
	}

	return &dictionary.Result{
		SearchWord:  searchWord,
		Entries:     dictEntries,
		Suggestions: nil,
	}, nil
}

func lookUpForPhrase(searchWord string, entry Entry) []dictionary.Entry {
	var de []dictionary.Entry

	for _, definedOnRun := range entry.Dros {
		if definedOnRun.Drp != searchWord {
			continue
		}

		definitions, _ := definedOnRun.Def.convert()
		dictEntry := dictionary.Entry{
			ID:              "mw-" + definedOnRun.Drp,
			Headword:        definedOnRun.Drp,
			FunctionalLabel: definedOnRun.Gram,
			Definitions:     definitions,
		}
		de = append(de, dictEntry)
	}

	return de
}

func lookUpForWord(searchWord string, entry Entry) []dictionary.Entry {
	var de []dictionary.Entry
	var matched bool

	if entry.Hwi.Hw.Clean() == searchWord {
		matched = true
	}

	definitions, _ := entry.Def.convert()
	dictEntry := dictionary.Entry{
		ID:              "mw-" + entry.Meta.ID,
		Headword:        entry.Hwi.Hw.Clean(),
		FunctionalLabel: entry.Fl,
		Pronunciation: dictionary.Pronunciation{
			Notation: "IPA",
			Accents:  entry.Hwi.Prs.convert(),
		},
		Inflections: entry.Ins.convert(),
		Definitions: definitions,
	}
	de = append(de, dictEntry)

	for _, uro := range entry.Uros {
		if uro.Ure.Clean() == searchWord {
			matched = true
		}

		definition, _ := convertDefiningText(uro.Utxt)
		dictEntry := dictionary.Entry{
			ID:              "mw-" + entry.Meta.ID + "-" + uro.Ure.Clean(),
			Headword:        uro.Ure.Clean(),
			FunctionalLabel: uro.Fl,
			Pronunciation: dictionary.Pronunciation{
				Notation: "IPA",
				Accents:  uro.Prs.convert(),
			},
			Inflections: nil,
			Definitions: []dictionary.Definition{definition},
		}
		de = append(de, dictEntry)
	}

	if matched {
		return de
	} else {
		return nil
	}
}
