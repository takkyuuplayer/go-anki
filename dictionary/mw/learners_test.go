package mw

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/takkyuuplayer/go-anki/dictionary"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func Test_learners_Parse(t *testing.T) {
	t.Parallel()

	learners := NewLearners("", nil)

	t.Run("Returning for a word", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("test", load(t, "test.json"))

		assert.Len(t, result.Entries, 3)
		assert.Equal(t, "test", result.Entries[0].Headword)
		assert.Equal(t, "test", result.Entries[1].Headword)
		assert.Equal(t, "testable", result.Entries[2].Headword)
		assert.Nil(t, err)
	})

	t.Run("Returning for a word under uros", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("accountability", load(t, "accountability.json"))

		assert.Len(t, result.Entries, 2)
		assert.Equal(t, "ac*count*able", result.Entries[0].Headword)
		assert.Equal(t, "ac*count*abil*i*ty", result.Entries[1].Headword)
		assert.Nil(t, err)
	})

	t.Run("Returning for a phrasal verb", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("go through", load(t, "go_through.json"))

		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "go through", result.Entries[0].Headword)
		assert.Nil(t, err)
	})

	t.Run("Returning NotFoundError for a phrasal verb", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("put up with", load(t, "put_up_with.json"))

		assert.Nil(t, result)
		assert.Equal(t, dictionary.NotFoundError, err)
	})

	t.Run("Returning suggestions", func(t *testing.T) {
		t.Parallel()
		result, err := learners.Parse("furnitura", load(t, "furnitura.json"))

		assert.Equal(t, "furnitura", result.SearchWord)
		assert.Len(t, result.Suggestions, 16)
		assert.Nil(t, result.Entries)
		assert.Nil(t, err)
	})

	filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
			return err
		}
		if info.IsDir() {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			t.Parallel()

			_, err = learners.Parse("dummy", load(t, filepath.Base(path)))
			assert.Nil(t, err)
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

func Test_learners_Search(t *testing.T) {
	t.Parallel()

	learners := NewLearners("dummy", &http.Client{})

	t.Run("Returning response body", func(t *testing.T) {
		t.Parallel()

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		resBody := load(t, "test.json")
		httpmock.RegisterResponder("GET", "https://www.dictionaryapi.com/api/v3/references/learners/json/test?key=dummy",
			httpmock.NewStringResponder(200, resBody),
		)

		body, err := learners.Search("test")

		assert.Equal(t, resBody, body)
		assert.Nil(t, err)
	})

	t.Run("Returning error when status code != 200", func(t *testing.T) {
		t.Parallel()

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		resBody := load(t, "test.json")
		httpmock.RegisterResponder("GET", "https://www.dictionaryapi.com/api/v3/references/learners/json/test?key=dummy",
			httpmock.NewStringResponder(404, resBody),
		)

		body, err := learners.Search("test")

		assert.Equal(t, "", body)
		assert.NotNil(t, err)
		assert.NotEqual(t, dictionary.NotFoundError, err)
	})

	t.Run("Returning NotFoundError when response is empty array i.e. no suggestions", func(t *testing.T) {
		t.Parallel()

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "https://www.dictionaryapi.com/api/v3/references/learners/json/test?key=dummy",
			httpmock.NewStringResponder(200, `[]`),
		)

		body, err := learners.Search("test")

		assert.Equal(t, "", body)
		assert.Equal(t, dictionary.NotFoundError, err)
	})
}
