package anki

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
const generator_error_message = "anki_card_generator_error"

func (ac *Client) Run(in io.Reader, out, outErr *csv.Writer) {
	counter := 0
	ch := make(chan *Result)
	defer close(ch)
	scanner := bufio.NewScanner(in)

	for ; counter < parallel; counter++ {
		if scanner.Scan() {
			go ac.SearchDefinition(ch, scanner.Text())
		} else {
			break
		}
	}

	failures := 0
	forceStop := false

	for ; counter > 0; counter-- {
		result := <-ch
		if result.IsSuccess {
			out.Write([]string{result.Word, result.Definition})

			if failures > 0 {
				failures--
			}

		} else {
			errMsg := fmt.Sprintf("%s: %s", result.Word, result.Definition)
			outErr.Write([]string{generator_error_message, errMsg})
			log.Print(errMsg)
			failures++
		}

		if failures > 3 {
			forceStop = true
		}

		if !forceStop && scanner.Scan() {
			go ac.SearchDefinition(ch, scanner.Text())
			counter++
		}
	}

	if forceStop {
		outErr.Write([]string{generator_error_message, "Too many failures. Force Stopped"})
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
