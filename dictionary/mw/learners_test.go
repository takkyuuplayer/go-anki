package mw

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

		assert.Len(t, result.Entries, 3)
		assert.Equal(t, "test", result.Entries[0].Headword)
		assert.Equal(t, "test", result.Entries[1].Headword)
		assert.Equal(t, "testable", result.Entries[2].Headword)
		assert.Nil(t, err)
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
