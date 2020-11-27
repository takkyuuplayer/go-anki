package anki

import (
	"log"
	"strings"
)

type Card struct {
	Headword string
	Entries  []*Entry
}

func (card Card) Front() string {
	return card.Headword
}

func (card Card) Back() (string, error) {
	ret := ""
	for _, entry := range card.Entries {
		content, err := entry.AnkiCard()
		if err != nil {
			log.Fatalf("execution failed: %s", err)
		} else {
			ret += " " + content
		}
	}
	ret = strings.TrimSpace(strings.Join(strings.Fields(ret), " "))

	return ret, nil
}
