package anki

import (
	"bufio"
	"encoding/csv"
	"io"
	"io/ioutil"
	"net/http"
)

type Result struct {
	Word       string
	Definition string
	IsSuccess  bool
}

type Dictionary interface {
	GetSearchUrl(string) string
	AnkiCard(string, string) (string, error)
}

type Client struct {
	Dictionary Dictionary
	HttpClient *http.Client
}

const parallel = 10

func (ac *Client) Run(in io.Reader, out, outErr *csv.Writer) {
	counter := 0
	ch := make(chan *Result)
	scanner := bufio.NewScanner(in)

	for ; counter < parallel; counter++ {
		if scanner.Scan() {
			go ac.SearchDefinition(ch, scanner.Text())
		} else {
			break
		}
	}

	for i := 0; i < counter; i++ {
		result := <-ch
		if result.IsSuccess {
			out.Write([]string{result.Word, result.Definition})
		} else {
			outErr.Write([]string{result.Word, result.Definition})
		}

		if scanner.Scan() {
			go ac.SearchDefinition(ch, scanner.Text())
			counter++
		}
	}

	out.Flush()
	outErr.Flush()
}

func (ac *Client) SearchDefinition(ch chan<- *Result, word string) {
	resp, err := ac.HttpClient.Get(ac.Dictionary.GetSearchUrl(word))

	if err != nil {
		ch <- &Result{
			Word:      word,
			IsSuccess: false,
		}
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ch <- &Result{
			Word:       "anki card generator error",
			Definition: err.Error(),
			IsSuccess:  false,
		}
		return
	}

	definition, err := ac.Dictionary.AnkiCard(string(body), word)

	if err != nil {
		ch <- &Result{
			Word:       "anki card generator error",
			Definition: err.Error(),
			IsSuccess:  false,
		}
		return
	}

	ch <- &Result{
		Word:       word,
		Definition: definition,
		IsSuccess:  true,
	}
}
