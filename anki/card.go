package anki

import (
	"bytes"
	"strings"

	"github.com/takkyuuplayer/go-anki/dictionary"
)

// Card is the raw data of anki card
type Card struct {
	SearchWord string
	Entries    []dictionary.Entry
}

// Front returns the content of front
func (card Card) Front() string {
	return card.SearchWord
}

// Back returns the content of back
func (card Card) Back() (string, error) {
	ret := ""
	for _, entry := range card.Entries {
		content, err := ankiCard(&entry)
		if err != nil {
			return "", err
		} else {
			ret += " " + content
		}
	}
	ret = strings.TrimSpace(strings.Join(strings.Fields(ret), " "))

	return ret, nil
}

func ankiCard(entry *dictionary.Entry) (string, error) {
	buf := bytes.NewBufferString("")

	if err := tmpl.Lookup("entry").Execute(buf, entry); err != nil {
		return "", err
	}

	return strings.Join(strings.Fields(buf.String()), " "), nil
}
