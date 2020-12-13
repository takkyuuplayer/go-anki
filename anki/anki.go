package anki

import (
	"bufio"
	"encoding/csv"
	"github.com/takkyuuplayer/go-anki/dictionary"
	"io"
	"strings"
	"sync"
)

const concurrency = 10
const errorNotFound = "error_not_found"
const errorGeneral = "error_general"

func Run(dictionaryApi dictionary.Dictionary, in io.Reader, out, outErr *csv.Writer) {
	c := make(chan bool, concurrency)
	wg := &sync.WaitGroup{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		wg.Add(1)
		c <- true
		go func(word string) {
			defer func() {
				<-c
				wg.Done()
			}()

			body, err := dictionaryApi.LookUp(word)
			if err == dictionary.NotFoundError {
				outErr.Write([]string{errorNotFound, word})
				return
			} else if err != nil {
				outErr.Write([]string{errorGeneral, word})
				return
			}
			result, err := dictionaryApi.Parse(word, body)
			if err == dictionary.NotFoundError {
				outErr.Write([]string{errorNotFound, word})
				return
			} else if err != nil {
				outErr.Write([]string{errorGeneral, word})
				return
			}
			if len(result.Suggestions) > 0 {
				outErr.Write([]string{errorNotFound, word})
				return
			} else if len(result.Entries) > 0 {
				card := Card{SearchWord: word, Entries: result.Entries}
				back, _ := card.Back()
				out.Write([]string{card.Front(), back})
			}

		}(strings.Trim(scanner.Text(), " "))
	}
	wg.Wait()

	out.Flush()
	outErr.Flush()
}
