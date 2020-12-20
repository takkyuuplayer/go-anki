package eijiro_test

import (
	"io/ioutil"
	"net/http"
	"testing"

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
		assert.Len(t, result.Entries[0].Definitions, 9)
		assert.Len(t, result.Entries[3].Definitions, 1)
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
