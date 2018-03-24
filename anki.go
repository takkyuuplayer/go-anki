package anki

import (
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
			Word:       word,
			Definition: err.Error(),
			IsSuccess:  false,
		}
		return
	}

	definition, err := ac.Dictionary.AnkiCard(string(body), word)

	if err != nil {
		ch <- &Result{
			Word:       word,
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
