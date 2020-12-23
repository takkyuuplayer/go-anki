package anki

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/takkyuuplayer/go-anki/dictionary"
)

const concurrency = 10
const errorNotFound = "Error: Not Found"
const errorUnknown = "Error: Unknown"
const errorLimit = 10

// Run generates anki card tsv file
func Run(dic dictionary.Dictionary, in io.Reader, out, outErr *csv.Writer) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	scanner := bufio.NewScanner(in)
	unkownErrCount := 0
	c := make(chan bool, concurrency)

	for scanner.Scan() {
		word := strings.Trim(scanner.Text(), " ")
		if word == "" {
			continue
		}
		if errorLimit <= unkownErrCount {
			break
		}

		wg.Add(1)
		c <- true
		go func(word string) {
			defer func() {
				<-c
				wg.Done()
			}()

			res, err := RunWord(dic, word)
			if err == dictionary.ErrNotFound {
				if unkownErrCount > 0 {
					unkownErrCount -= 1
				}
			} else if err != nil {
				unkownErrCount += 1
			} else {
				if unkownErrCount > 0 {
					unkownErrCount -= 1
				}
			}

			mu.Lock()
			defer mu.Unlock()

			if err == nil {
				out.Write(res)
			} else {
				outErr.Write(res)
			}
		}(word)
	}
	wg.Wait()

	if unkownErrCount >= errorLimit {
		outErr.Write([]string{errorUnknown, "Stopped because of too many unknown errors"})
	}
	out.Flush()
	outErr.Flush()
}

func RunWord(dic dictionary.Dictionary, word string) ([]string, error) {
	body, err := dic.LookUp(word)
	if err == dictionary.ErrNotFound {
		return []string{errorNotFound, word}, err
	} else if err != nil {
		return []string{errorUnknown, fmt.Sprintf("%s: %s", word, err)}, err
	}

	result, err := dic.Parse(word, body)
	if err == dictionary.ErrNotFound {
		return []string{errorNotFound, word}, err
	} else if err != nil {
		return []string{errorUnknown, fmt.Sprintf("%s: %s", word, err)}, err
	}
	if len(result.Suggestions) > 0 {
		return []string{errorNotFound, word}, dictionary.ErrNotFound
	}

	card := Card{SearchWord: word, Entries: result.Entries}
	back, err := card.Back()
	if err != nil {
		log.Println(err)
		return []string{errorUnknown, fmt.Sprintf("%s: %s", word, err)}, err
	}
	return []string{card.Front(), back}, nil
}
