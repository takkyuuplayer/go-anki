package anki_test

import (
	"testing"

	. "github.com/takkyuuplayer/go-anki/anki"
	_ "github.com/takkyuuplayer/go-anki/anki/statik"
)

var testEntry = Entry{
	"test:1",
	"test",
	"noun",
	Pronunciation{
		"IPA",
		[]Accent{
			{
				"US",
				"ˈtɛst",
				"https://example.com/test.mp3",
			},
		},
	},
	[]Inflection{
		{
			"plural",
			"tests",
			nil,
		},
	},
	[]Definition{
		{
			"a set of questions or problems",
			[]string{
				"She is studying for her math/spelling/history test",
				"I passed/failed/flunked my biology test",
			},
		},
	},
}

func TestEntry_AnkiCard(t *testing.T) {
	tests := []struct {
		name   string
		fields *Entry
		want   string
	}{
		{
			"test",
			&testEntry,
			"<h2>test (noun)</h2> <h4>Pronunciation (IPA)</h4> US ˈtɛst <h4>Inflection</h4> (plural) tests <h4>Definition</h4> <ol> <li>a set of questions or problems</li> <ul> <li>She is studying for her math/spelling/history test</li> <li>I passed/failed/flunked my biology test</li> </ul> </ol>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.fields.AnkiCard(); got != tt.want {
				t.Errorf("mustParseAssets() = %v, want %v", got, tt.want)
			}
		})
	}
}
