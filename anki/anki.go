package anki

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/takkyuuplayer/go-anki/dictionary"
)

const concurrency = 10
const errorNotFound = "Error: Not Found"
const errorUnknown = "Error: Unknown"
const errorLimit = 10

// Run generates anki card tsv file
func Run(dic dictionary.Dictionary, in io.Reader, out, outErr *csv.Writer) {
	scanner := bufio.NewScanner(in)
	unkownErrCount := 0
	type wordResult struct {
		out []string
		err error
	}
	c := make(chan wordResult, concurrency)

	for scanner.Scan() {
		word := strings.Trim(scanner.Text(), " ")
		if word == "" {
			continue
		}
		if errorLimit <= unkownErrCount {
			break
		}

		go func(word string) {
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
			c <- wordResult{res, err}
		}(word)

		res := <-c
		if res.err == nil {
			out.Write(res.out)
		} else {
			outErr.Write(res.out)
		}
	}
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
		return []string{errorNotFound, word}, err
	}

	card := Card{SearchWord: word, Entries: result.Entries}
	back, _ := card.Back()
	return []string{card.Front(), back}, nil
}
