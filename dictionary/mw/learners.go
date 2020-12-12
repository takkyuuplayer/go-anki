package mw

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/takkyuuplayer/go-anki/dictionary"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (dic *learners) Search(word string) (string, error) {
	url := fmt.Sprintf(searchUrl, url.PathEscape(word), url.PathEscape(dic.apiKey))

	response, err := dic.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

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

	var dictEntries []dictionary.Entry
	for _, entry := range entries {
		definitions, _ := entry.Def.convert()
		dictEntry := dictionary.Entry{
			ID:              "mw-" + entry.Meta.ID,
			Headword:        entry.Hwi.Hw,
			FunctionalLabel: entry.Fl,
			Pronunciation: dictionary.Pronunciation{
				Notation: "IPA",
				Accents:  entry.Hwi.Prs.convert(),
			},
			Inflections: entry.Ins.convert(),
			Definitions: definitions,
		}
		dictEntries = append(dictEntries, dictEntry)
	}

	return &dictionary.Result{
		SearchWord:  searchWord,
		Entries:     dictEntries,
		Suggestions: nil,
	}, nil
}
