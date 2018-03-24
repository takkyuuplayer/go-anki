package mw

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
)

type MerriamWebster struct {
	apiKey     string
	dictionary string
}

func New(key, dictionary string) *MerriamWebster {
	return &MerriamWebster{key, dictionary}
}

func (m *MerriamWebster) GetSearchUrl(word string) string {
	return fmt.Sprintf(
		"https://www.dictionaryapi.com/api/v1/references/%s/xml/%s?key=%s",
		url.PathEscape(m.dictionary),
		url.PathEscape(word),
		url.PathEscape(m.apiKey),
	)
}

func (c *MerriamWebster) AnkiCard(body, word string) (string, error) {
	var el EntryList
	xml.Unmarshal([]byte(body), &el)

	if len(el.Entries) == 0 {
		return "", errors.New("Not Found")
	}

	ret := el.AnkiCard(word)

	if ret == "" {
		return "", errors.New("Not Found")
	}

	return ret, nil
}
