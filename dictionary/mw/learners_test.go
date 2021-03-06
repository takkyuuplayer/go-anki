package mw_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/takkyuuplayer/go-anki/dictionary/mw"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/takkyuuplayer/go-anki/dictionary"
)

func Test_learners_Parse(t *testing.T) {
	t.Parallel()

	learners := mw.NewLearners("", nil)

	t.Run("Returning for a word", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("test", load(t, "test.json"))

		assert.Len(t, result.Entries, 3)
		assert.Equal(t, "Merriam-Webster's Learner's Dictionary", result.Dictionary)
		assert.Equal(t, "test", result.Entries[0].Headword)
		assert.Equal(t, "test", result.Entries[1].Headword)
		assert.Equal(t, "testable", result.Entries[2].Headword)
		assert.NotEqual(t, "", result.Entries[0].Definitions[0].Sense)
		assert.Nil(t, err)
	})

	t.Run("Returning for a word: deathbed", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("deathbed", load(t, "deathbed.json"))

		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "Merriam-Webster's Learner's Dictionary", result.Dictionary)
		assert.Equal(t, "deathbed", result.Entries[0].Headword)
		assert.Len(t, result.Entries[0].Pronunciation.Accents, 1)
		assert.Nil(t, err)
	})

	t.Run("Returning for a word having empty definition", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("endorse", load(t, "endorse.json"))

		assert.Len(t, result.Entries, 2)
		assert.Equal(t, "endorse", result.Entries[0].Headword)
		assert.Equal(t, "endorser", result.Entries[1].Headword)
		assert.Len(t, result.Entries[1].Definitions, 0)
		assert.Nil(t, err)
	})

	t.Run("Returning for a word under uros", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("accountability", load(t, "accountability.json"))

		assert.Len(t, result.Entries, 2)
		assert.Equal(t, "accountable", result.Entries[0].Headword)
		assert.Equal(t, "accountability", result.Entries[1].Headword)
		assert.Nil(t, err)
	})

	t.Run("Returning in original form", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("muddies", load(t, "muddies.json"))

		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "muddy", result.Entries[0].Headword)
		assert.Nil(t, err)
	})

	t.Run("Returning stratum for strata", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("strata", load(t, "strata.json"))

		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "stratum", result.Entries[0].Headword)
		assert.Len(t, result.Entries[0].Definitions, 2)
		assert.Len(t, result.Entries[0].Pronunciation.Accents, 2)
		assert.Equal(t, "US", result.Entries[0].Pronunciation.Accents[0].AccentLabel)
		assert.Equal(t, "British", result.Entries[0].Pronunciation.Accents[1].AccentLabel)
		assert.Nil(t, err)
	})

	t.Run("Returning for a phrasal verb", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("go through", load(t, "go_through.json"))

		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "go through", result.Entries[0].Headword)
		assert.Nil(t, err)

		t.Run("Returning from snote", func(t *testing.T) {
			assert.NotEqual(t, "", result.Entries[0].Definitions[6].Sense)
			assert.Len(t, result.Entries[0].Definitions[6].Examples, 3)
		})
	})

	t.Run("Returning ErrNotFound for a phrasal verb", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("put up with", load(t, "put_up_with.json"))

		assert.Nil(t, result)
		assert.Equal(t, dictionary.ErrNotFound, err)
	})

	t.Run("Returning suggestions", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("furnitura", load(t, "furnitura.json"))

		assert.Equal(t, "Merriam-Webster's Learner's Dictionary", result.Dictionary)
		assert.Len(t, result.Suggestions, 16)
		assert.Nil(t, result.Entries)
		assert.Nil(t, err)
	})

	t.Run("Returning for a idiom", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("under way", load(t, "under_way.json"))

		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "under way", result.Entries[0].Headword)
		assert.Nil(t, err)
	})

	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal()
			return nil
		}
		if info.IsDir() {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			_, err := learners.Parse("dummy", load(t, filepath.Base(path)))

			assert.True(t, err == nil || errors.Is(dictionary.ErrNotFound, err))
		})
		return nil
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

func Test_learners_LookUp(t *testing.T) {
	t.Parallel()

	learners := mw.NewLearners("dummy", &http.Client{})

	t.Run("Returning response body", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		resBody := load(t, "test.json")
		httpmock.RegisterResponder("GET", "https://www.dictionaryapi.com/api/v3/references/learners/json/test?key=dummy",
			httpmock.NewStringResponder(200, resBody),
		)

		body, err := learners.LookUp("test")
		assert.Equal(t, resBody, body)
		assert.Nil(t, err)
	})

	t.Run("Returning error when status code != 200", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		resBody := load(t, "test.json")
		httpmock.RegisterResponder("GET", "https://www.dictionaryapi.com/api/v3/references/learners/json/test?key=dummy",
			httpmock.NewStringResponder(404, resBody),
		)

		body, err := learners.LookUp("test")

		assert.Equal(t, "", body)
		assert.NotNil(t, err)
		assert.NotEqual(t, dictionary.ErrNotFound, err)
	})

	t.Run("Returning ErrNotFound when response is empty array i.e. no suggestions", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "https://www.dictionaryapi.com/api/v3/references/learners/json/test?key=dummy",
			httpmock.NewStringResponder(200, `[]`),
		)

		body, err := learners.LookUp("test")

		assert.Equal(t, "", body)
		assert.Equal(t, dictionary.ErrNotFound, err)
	})
}
