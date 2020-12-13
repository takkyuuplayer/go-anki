package anki

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/takkyuuplayer/go-anki/dictionary"
)

const concurrency = 10
const errorNotFound = "Error: Not Found"
const errorUnknown = "Error: Unknown"

// Run generates anki card tsv file
func Run(dic dictionary.Dictionary, in io.Reader, out, outErr *csv.Writer) {
	c := make(chan bool, concurrency)
	wg := &sync.WaitGroup{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		word := strings.Trim(scanner.Text(), " ")
		if word == "" {
			continue
		}

		fmt.Println("Start ", word)

		wg.Add(1)
		c <- true
		go func(word string) {
			defer func() {
				fmt.Println("End ", word)
				<-c
				wg.Done()
			}()

			body, err := dic.LookUp(word)
			if err == dictionary.ErrNotFound {
				outErr.Write([]string{errorNotFound, word})
				return
			} else if err != nil {
				outErr.Write([]string{errorUnknown, fmt.Sprintf("%s: %s", word, err)})
				return
			}
			result, err := dic.Parse(word, body)
			if err == dictionary.ErrNotFound {
				outErr.Write([]string{errorNotFound, word})
				return
			} else if err != nil {
				outErr.Write([]string{errorUnknown, fmt.Sprintf("%s: %s", word, err)})
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

		}(word)
	}
	wg.Wait()

	out.Flush()
	outErr.Flush()
}
