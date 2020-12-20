package eijiro_test

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/takkyuuplayer/go-anki/dictionary"

	"github.com/takkyuuplayer/go-anki/dictionary/eijiro"

	"github.com/stretchr/testify/assert"
)

func TestEijiro_Parse(t *testing.T) {
	t.Parallel()

	dic := eijiro.NewEijiro(http.DefaultClient)
	t.Run("Returning for a word", func(t *testing.T) {
		t.Parallel()

		result, err := dic.Parse("test", load(t, "test.html"))

		assert.Len(t, result.Entries, 4)

		assert.Equal(t, "eijiro-test-【1名】", result.Entries[0].ID)
		assert.Equal(t, "test", result.Entries[0].Headword)
		assert.Equal(t, "名", result.Entries[0].FunctionalLabel)
		assert.Equal(t, &dictionary.Pronunciation{
			Notation: "IPA",
			Accents:  []dictionary.Accent{{AccentLabel: "-", Spelling: "tɛst"}},
		}, result.Entries[0].Pronunciation)
		assert.Nil(t, result.Entries[0].Inflections)
		assert.Len(t, result.Entries[0].Definitions, 9)
		assert.Equal(t,
			dictionary.Definition{Sense: "〔教育の〕試験、考査、テスト", Examples: []template.HTML{`How'd the test go? : テストどうだった？◆How'dはHow didの略で、口語的表現。`}},
			result.Entries[0].Definitions[0])
		assert.Equal(t,
			dictionary.Definition{Sense: "〔機器や製法などの〕検査、試験運転、動作確認", Examples: nil},
			result.Entries[0].Definitions[1])

		assert.Len(t, result.Entries[1].Inflections, 3)
		assert.Equal(t, dictionary.Inflection{InflectedForm: "tests"}, result.Entries[1].Inflections[0])
		assert.Equal(t, dictionary.Inflection{InflectedForm: "testing"}, result.Entries[1].Inflections[1])
		assert.Equal(t, dictionary.Inflection{InflectedForm: "tested"}, result.Entries[1].Inflections[2])

		assert.Nil(t, err)
	})

	t.Run("Returning for multiple pronunciation word", func(t *testing.T) {
		t.Parallel()

		result, err := dic.Parse("protest", load(t, "protest.html"))

		assert.Len(t, result.Entries, 3)

		assert.Equal(t, &dictionary.Pronunciation{
			Notation: "IPA",
			Accents:  []dictionary.Accent{{AccentLabel: "-", Spelling: "próutest"}},
		}, result.Entries[0].Pronunciation)
		assert.Equal(t, &dictionary.Pronunciation{
			Notation: "IPA",
			Accents:  []dictionary.Accent{{AccentLabel: "-", Spelling: "prətést"}},
		}, result.Entries[1].Pronunciation)

		assert.Nil(t, err)
	})
}

func load(t *testing.T, testfile string) string {
	t.Helper()

	body, err := ioutil.ReadFile("testdata/" + testfile)
	if err != nil {
		t.Fatal(err)
	}

	return string(body)
}
