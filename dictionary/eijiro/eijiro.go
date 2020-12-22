package eijiro

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

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

var mu sync.Mutex

func (dic *Eijiro) LookUp(word string) (string, error) {
	mu.Lock()
	defer func() {
		time.Sleep(2 * time.Second)
		mu.Unlock()
	}()

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
	if definitions == nil {
		ret := dictionary.Result{SearchWord: word}
		suggestions := cascadia.QueryAll(doc, cascadia.MustCompile("#sas_word > span > a"))

		if len(suggestions) > 0 {
			for _, suggestion := range suggestions {
				ret.Suggestions = append(ret.Suggestions, text(suggestion))
			}
		}
		return &ret, dictionary.ErrNotFound
	}

	headword := text(cascadia.Query(definitions, cascadia.MustCompile("span.midashi > h2 > span")))
	if headword == "" {
		return nil, dictionary.ErrNotFound
	}

	var dictEntries []dictionary.Entry
	var dictEntry dictionary.Entry

	first := cascadia.Query(definitions, cascadia.MustCompile("div > span:nth-child(1)"))
	attrRootIdx := -1
	for node := first; node != nil; node = node.NextSibling {
		switch node.Data {
		case "span":
			switch node.Attr[0].Val {
			case "wordclass":
				if dictEntry.ID != "" {
					dictEntries = append(dictEntries, dictEntry)
				}
				if attrRootIdx == -1 {
					attrRootIdx = len(dictEntries)
				}
				dictEntry = dictionary.Entry{
					ID:              fmt.Sprintf("eijiro-%s-%s", headword, node.FirstChild.Data),
					Headword:        headword,
					FunctionalLabel: strings.Trim(text(node), "【0123456789】"),
				}
			case "attr":
				dictEntries = append(dictEntries, dictEntry)
				dictEntry = dictionary.Entry{}

				parseAttr(node, attrRootIdx, dictEntries)

				attrRootIdx = -1
			}
		case "ol":
			attr := cascadia.Query(node, cascadia.MustCompile("li > span.attr"))
			if attr != nil {
				dictEntries = append(dictEntries, dictEntry)
				dictEntry = dictionary.Entry{}

				parseAttr(attr, attrRootIdx, dictEntries)

				attrRootIdx = -1
				continue
			}

			var definitions []dictionary.Definition
			if node.FirstChild.Data == "li" {
				for li := node.FirstChild; li != nil; li = li.NextSibling {
					sense, examples := parseDefinition(li.FirstChild)
					definitions = append(definitions, dictionary.Definition{
						Sense:    template.HTML(sense),
						Examples: examples,
					})
				}
			} else {
				sense, examples := parseDefinition(node.FirstChild)
				definitions = append(definitions, dictionary.Definition{
					Sense:    template.HTML(sense),
					Examples: examples,
				})
			}
			dictEntry.Definitions = definitions
		}
	}
	if dictEntry.ID != "" {
		dictEntries = append(dictEntries, dictEntry)
	}

	return &dictionary.Result{SearchWord: word, Entries: dictEntries}, nil
}

func parseAttr(node *html.Node, attrRootIdx int, dictEntries []dictionary.Entry) {
	var pronunciationText, inflectionTest string
	for node2 := node.FirstChild; node2 != nil; node2 = node2.NextSibling {
		if node2.Data == "span" && node2.Attr[0].Val == "pron" {
			pronunciationText = text(node2)
		} else if node2.Data == "span" && text(node2) == "【変化】" {
			inflectionTest = text(node2.NextSibling)
		}
	}

	if strings.HasPrefix(pronunciationText, "《") {
		for _, pronunciation := range strings.Split(pronunciationText, "《") {
			if pronunciation == "" {
				continue
			}
			functionLabel := string([]rune(pronunciation)[0:1])
			pronunciation = string([]rune(pronunciation)[2:])

			for idx := attrRootIdx; idx < len(dictEntries); idx++ {
				if strings.Contains(dictEntries[idx].FunctionalLabel, functionLabel) {
					dictEntries[idx].Pronunciation = parsePronunciation(pronunciation)
				}
			}
		}
	} else {
		dictEntries[attrRootIdx].Pronunciation = parsePronunciation(pronunciationText)
	}

	if strings.HasPrefix(inflectionTest, "《") {
		for _, inflection := range strings.Split(inflectionTest, "《") {
			if inflection == "" {
				continue
			}
			functionLabel := string([]rune(inflection)[0:1])
			inflection = string([]rune(inflection)[2:])

			for idx := attrRootIdx; idx < len(dictEntries); idx++ {
				if strings.Contains(dictEntries[idx].FunctionalLabel, functionLabel) {
					infs := strings.Split(inflection, "｜")
					for _, inf := range infs {
						dictEntries[idx].Inflections = append(
							dictEntries[idx].Inflections,
							dictionary.Inflection{InflectedForm: strings.TrimSpace(inf)},
						)
					}
				} else if dictEntries[idx].FunctionalLabel == "名" && functionLabel == "複" {
					infs := strings.Split(inflection, "｜")
					for _, inf := range infs {
						dictEntries[idx].Inflections = append(
							dictEntries[idx].Inflections,
							dictionary.Inflection{InflectedForm: strings.TrimSpace(inf)},
						)
					}
				}
			}
		}
	} else {
		dictEntries[attrRootIdx].Pronunciation = parsePronunciation(pronunciationText)
	}
}

func parseDefinition(first *html.Node) (string, []template.HTML) {
	sense := ""
	var examples []template.HTML = nil
	inExampleNode := false
	for node2 := first; node2 != nil; node2 = node2.NextSibling {
		switch node2.Data {
		case "span":
			if node2.Attr[0].Val != "kana" {
				sense += node2.FirstChild.Data
			}
		case "br":
			inExampleNode = true
		default:
			if inExampleNode {
				examples = append(examples, template.HTML(strings.TrimLeft(node2.Data, "・")))
			} else {
				sense += node2.Data
			}
		}
	}
	return sense, examples
}

func parsePronunciation(text string) *dictionary.Pronunciation {
	// [US] wə́rd ｜ [UK] wə́ːd、
	var accents []dictionary.Accent
	for _, accent := range strings.Split(text, "｜") {
		accent = strings.TrimSpace(accent)
		if strings.HasPrefix(accent, "[") {
			label := accent[strings.Index(accent, "[")+1 : strings.LastIndex(accent, "]")]
			spelling := accent[strings.LastIndex(accent, "]")+1:]
			accents = append(accents, dictionary.Accent{AccentLabel: label, Spelling: strings.TrimSpace(strings.TrimRight(spelling, "、"))})
		} else {
			accents = append(accents, dictionary.Accent{AccentLabel: "-", Spelling: strings.TrimSpace(strings.TrimRight(accent, "、"))})
		}
	}
	return &dictionary.Pronunciation{Notation: "IPA", Accents: accents}
}

func text(node *html.Node) string {
	s := ""
	if node.Type == html.TextNode {
		s += " " + node.Data
	} else {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			s += " " + text(child)
		}
	}
	return strings.TrimSpace(s)
}
