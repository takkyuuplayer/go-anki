package eijiro

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/andybalholm/cascadia"

	"github.com/takkyuuplayer/go-anki/dictionary"
	"golang.org/x/net/html"
)

const searchURL = "https://eow.alc.co.jp/search?q=%s"

type Eijiro struct {
	httpClient *http.Client
}

func NewEijiro(httpClient *http.Client) *Eijiro {
	return &Eijiro{httpClient: httpClient}
}

func (dic *Eijiro) LookUp(word string) (string, error) {
	urlToSearch := fmt.Sprintf(searchURL, url.QueryEscape(word))

	response, err := dic.httpClient.Get(urlToSearch)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	_ = response.Body.Close()

	bodyText := string(body)
	if response.StatusCode != http.StatusOK {
		return "", errors.New(bodyText)
	}

	return bodyText, nil
}

func (dic *Eijiro) Parse(word, body string) (*dictionary.Result, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to pase HTML: %v", err)
	}

	definitions := cascadia.Query(doc, cascadia.MustCompile("#resultsList > ul > li:nth-child(1)"))
	if err != nil {
		return nil, dictionary.ErrNotFound
	}

	headword := cascadia.Query(definitions, cascadia.MustCompile("span.midashi > h2 > span")).
		FirstChild.Data
	if headword == "" {
		return nil, dictionary.ErrNotFound
	}

	var dictEntries []dictionary.Entry
	var dictEntry dictionary.Entry

	first := cascadia.Query(definitions, cascadia.MustCompile("div > span:nth-child(1)"))
	attrRoot := -1
	for node := first; node != nil; node = node.NextSibling {
		switch node.Data {
		case "span":
			switch node.Attr[0].Val {
			case "wordclass":
				if dictEntry.ID != "" {
					dictEntries = append(dictEntries, dictEntry)
				}
				if attrRoot == -1 {
					attrRoot = len(dictEntries) - 1
				}
				dictEntry = dictionary.Entry{
					ID:       fmt.Sprintf("eijiro-%s-%s", headword, node.FirstChild.Data),
					Headword: headword,
				}
			case "attr":
				dictEntries = append(dictEntries, dictEntry)
			}
		case "ol":
			var definitions []dictionary.Definition
			for li := node.FirstChild; li != nil; li = li.NextSibling {
				sense := ""
				var examples []template.HTML
				inExampleNode := false
				for node2 := li.FirstChild; node2 != nil; node2 = node2.NextSibling {
					switch node2.Data {
					case "span": // Nothing to do
					case "br":
						inExampleNode = true
					default:
						if inExampleNode {
							examples = append(examples, template.HTML(strings.TrimLeft("・", node2.Data)))
						} else {
							sense += node2.Data
						}
					}
				}
				definitions = append(definitions, dictionary.Definition{
					Sense:    template.HTML(sense),
					Examples: examples,
				})
			}
			dictEntry.Definitions = definitions
		}
	}

	return &dictionary.Result{SearchWord: word, Entries: dictEntries}, nil
}
