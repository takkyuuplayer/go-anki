package anki_test

import (
	"html/template"
	"testing"

	"github.com/takkyuuplayer/go-anki/anki"
	"github.com/takkyuuplayer/go-anki/dictionary"
)

func TestCard_Back(t *testing.T) {
	tests := []struct {
		name    string
		fields  anki.Card
		want    string
		wantErr bool
	}{
		{
			name: "test",
			fields: anki.Card{
				SearchWord: "test",
				Result: &dictionary.Result{
					Dictionary: "ExampleDict",
					WebURL:     "https://example.com",
					Entries: []dictionary.Entry{
						testEntry,
						testEntry,
					},
				},
			},
			want: `<h2>test (noun)</h2> <h4>Pronunciation (IPA)</h4> [US] ˈtɛst <h4>Inflection</h4> (plural) tests <h4>Definition</h4> <ol> <li>a set of questions or problems</li> <ul> <li>She is studying for her math/spelling/history test</li> <li>I passed/failed/flunked my biology test</li> </ul> </ol> <h2>test (noun)</h2> <h4>Pronunciation (IPA)</h4> [US] ˈtɛst <h4>Inflection</h4> (plural) tests <h4>Definition</h4> <ol> <li>a set of questions or problems</li> <ul> <li>She is studying for her math/spelling/history test</li> <li>I passed/failed/flunked my biology test</li> </ul> </ol><hr><a href="https://example.com">test - ExampleDict</a>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Back()
			if (err != nil) != tt.wantErr {
				t.Errorf("Card.Back() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Card.Back() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testEntry dictionary.Entry = dictionary.Entry{
	ID:              "test:1",
	Headword:        "test",
	FunctionalLabel: "noun",
	Pronunciation: &dictionary.Pronunciation{
		Notation: "IPA",
		Accents: []dictionary.Accent{
			{
				AccentLabel: "US",
				Spelling:    "ˈtɛst",
				Audio:       "https://example.com/test.mp3",
			},
		},
	},
	Inflections: []dictionary.Inflection{
		{
			FormLabel:     "plural",
			InflectedForm: "tests",
		},
	},
	Definitions: []dictionary.Definition{
		{
			Sense: "a set of questions or problems",
			Examples: []template.HTML{
				"She is studying for her math/spelling/history test",
				"I passed/failed/flunked my biology test",
			},
		},
	},
}
