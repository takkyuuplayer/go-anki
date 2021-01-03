package mw

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/takkyuuplayer/go-anki/dictionary"
)

const dictionaryName = "Merriam-Webster's Learner's Dictionary"
const searchURL = "https://www.dictionaryapi.com/api/v3/references/learners/json/%s?key=%s"
const webURL = "https://learnersdictionary.com/definition/%s"

// Learners is a client to access MERRIAM-WEBSTER'S LEARNER'S DICTIONARY API
// https://dictionaryapi.com/products/api-learners-dictionary
type Learners struct {
	apiKey     string
	httpClient *http.Client
}

// NewLearners returns an instance of learner's dictionary API
func NewLearners(apiKey string, httpClient *http.Client) *Learners {
	return &Learners{apiKey: apiKey, httpClient: httpClient}
}

// LookUp looks up the word on the dictionary
func (dic *Learners) LookUp(word string) (string, error) {
	urlToSearch := fmt.Sprintf(searchURL, url.PathEscape(word), url.PathEscape(dic.apiKey))

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
		return "", dictionary.ErrNotFound
	}

	return bodyText, nil
}

// Parse parses the response body of the looked up result
func (dic *Learners) Parse(searchWord, body string) (*dictionary.Result, error) {
	var suggestion suggestion
	err := json.Unmarshal([]byte(body), &suggestion)
	if err == nil {
		return &dictionary.Result{
			Dictionary:  dictionaryName,
			Entries:     nil,
			Suggestions: suggestion,
			WebURL:      template.URL(fmt.Sprintf(webURL, url.PathEscape(searchWord))),
		}, nil
	}

	var entries entries
	err = json.Unmarshal([]byte(body), &entries)
	if err != nil {
		return nil, fmt.Errorf("unknown structure: %v", err)
	}

	var dictEntries []dictionary.Entry
	for _, entry := range entries {
		if lookedUp, err := parseEntry(searchWord, entry); err != nil {
			return nil, err
		} else if lookedUp != nil {
			dictEntries = append(dictEntries, lookedUp...)
		}
	}

	if len(dictEntries) == 0 {
		return nil, dictionary.ErrNotFound
	}

	return &dictionary.Result{
		Dictionary:  dictionaryName,
		Entries:     dictEntries,
		Suggestions: nil,
		WebURL:      template.URL(fmt.Sprintf(webURL, url.PathEscape(searchWord))),
	}, nil
}

func parseEntry(searchWord string, entry entry) ([]dictionary.Entry, error) {
	if len(entry.Shortdef) == 0 && entry.Meta.AppShortdef == nil {
		return nil, nil
	}

	var de []dictionary.Entry

	var pronunciation *dictionary.Pronunciation
	if len(entry.Hwi.Prs) > 0 {
		pronunciation = &dictionary.Pronunciation{
			Notation: "IPA",
			Accents:  entry.Hwi.Prs.convert(),
		}
	}

	if isEntryDefMatches(searchWord, entry) {
		definitions, err := entry.Def.convert()
		if err != nil {
			return nil, err
		}

		if definitions != nil {
			dictEntry := dictionary.Entry{
				ID:              "mw-" + entry.Meta.ID,
				Headword:        clean(entry.Hwi.Hw),
				FunctionalLabel: entry.Fl,
				Pronunciation:   pronunciation,
				Inflections:     entry.Ins.convert(),
				Definitions:     definitions,
			}
			de = append(de, dictEntry)
		}

		for _, uro := range entry.Uros {
			var definitions []dictionary.Definition
			if len(uro.Utxt) > 0 {
				definition, err := convertDefiningText(uro.Utxt)
				if err != nil {
					return nil, err
				}
				definitions = append(definitions, definition)
			}
			if len(uro.Prs) > 0 {
				pronunciation = &dictionary.Pronunciation{
					Notation: "IPA",
					Accents:  uro.Prs.convert(),
				}
			}

			dictEntry := dictionary.Entry{
				ID:              "mw-" + entry.Meta.ID + "-" + clean(uro.Ure),
				Headword:        clean(uro.Ure),
				FunctionalLabel: uro.Fl,
				Pronunciation:   pronunciation,
				Inflections:     uro.Ins.convert(),
				Definitions:     definitions,
			}
			de = append(de, dictEntry)
		}
	}

	for _, definedOnRun := range entry.Dros {
		if definedOnRun.Drp != searchWord {
			continue
		}

		definitions, err := definedOnRun.Def.convert()
		if err != nil {
			return nil, err
		}
		dictEntry := dictionary.Entry{
			ID:              "mw-" + definedOnRun.Drp,
			Headword:        definedOnRun.Drp,
			FunctionalLabel: definedOnRun.Gram,
			Definitions:     definitions,
		}
		de = append(de, dictEntry)
	}

	return de, nil
}

func isEntryDefMatches(searchWord string, entry entry) bool {
	if clean(entry.Hwi.Hw) == searchWord {
		return true
	}
	for _, in := range entry.Ins {
		if clean(in.If) == searchWord {
			return true
		}
	}

	if isPhrase := len(strings.Fields(searchWord)) > 1; isPhrase {
		return false
	}

	for _, stem := range entry.Meta.Stems {
		if clean(stem) == searchWord {
			return true
		}
	}

	return false
}
