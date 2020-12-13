package anki

import (
	"bytes"
	"log"
	"strings"

	"github.com/takkyuuplayer/go-anki/dictionary"
)

type Card struct {
	SearchWord string
	Entries    []dictionary.Entry
}

func (card Card) Front() string {
	return card.SearchWord
}

func (card Card) Back() (string, error) {
	ret := ""
	for _, entry := range card.Entries {
		content, err := ankiCard(&entry)
		if err != nil {
			log.Fatalf("execution failed: %s", err)
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
		log.Fatalf("execution failed: %s", err)
		return "", err
	}

	return strings.Join(strings.Fields(buf.String()), " "), nil
}
