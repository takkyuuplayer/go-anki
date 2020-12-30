package anki

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/takkyuuplayer/go-anki/dictionary"
)

// Card is the raw data of anki card
type Card struct {
	SearchWord string
	Result     *dictionary.Result
}

// Front returns the content of front
func (card Card) Front() string {
	return card.Result.Entries[0].Headword
}

// Back returns the content of back
func (card Card) Back() (string, error) {
	ret := ""
	for _, entry := range card.Result.Entries {
		content, err := ankiCard(&entry)
		if err != nil {
			return "", err
		} else {
			ret += " " + content
		}
	}
	ret += fmt.Sprintf(`<hr><a href="%s">%s - %s</a>`,
		card.Result.WebUrl,
		html.EscapeString(card.SearchWord),
		html.EscapeString(card.Result.Dictionary),
	)
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
