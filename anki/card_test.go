package anki_test

import (
	"testing"

	. "github.com/takkyuuplayer/go-anki/anki"
)

func TestCard_Back(t *testing.T) {
	tests := []struct {
		name    string
		fields  Card
		want    string
		wantErr bool
	}{
		{
			"test",
			Card{
				"test",
				[]*Entry{
					&testEntry,
					&testEntry,
				},
			},
			"<h2>test (noun)</h2> <h4>Pronunciation (IPA)</h4> US ˈtɛst <h4>Inflection</h4> (plural) tests <h4>Definition</h4> <ol> <li>a set of questions or problems</li> <ul> <li>She is studying for her math/spelling/history test</li> <li>I passed/failed/flunked my biology test</li> </ul> </ol> <h2>test (noun)</h2> <h4>Pronunciation (IPA)</h4> US ˈtɛst <h4>Inflection</h4> (plural) tests <h4>Definition</h4> <ol> <li>a set of questions or problems</li> <ul> <li>She is studying for her math/spelling/history test</li> <li>I passed/failed/flunked my biology test</li> </ul> </ol>",
			false,
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
