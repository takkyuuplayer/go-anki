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

	t.Run("When finding definition under uros", func(t *testing.T) {
		t.Parallel()
		result, _:= learners.Parse("accountability", load(t, "accountability.json"))

		t.Logf("%#v", result)
		//assert.Len(t, result.Entries, 2)
		//assert.Nil(t, err)
	})
	t.Run("Returns suggestions when available", func(t *testing.T) {
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
