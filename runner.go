package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var client *http.Client
var defReg = regexp.MustCompile(`(?s)<h2><span class="mw-headline" id="English">English</span>(?:.+?)</h2>.*?(<h3>.+?)\n+(?:<hr />|<!-- \nNewPP limit report)`)
var stdout = csv.NewWriter(os.Stdout)

var ignoreParagraphs = []string{
	"Etymology",
	"Derived_terms",
	"Translations",
	"Further_reading",
	"References",
	"Anagrams",
}

var ignoreRegexps = make([]*regexp.Regexp, len(ignoreParagraphs))
var deleteReg = regexp.MustCompile(`<span class="mw-editsection">.+</span></span>`)

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

func main() {
	for idx, val := range ignoreParagraphs {
		ignoreRegexps[idx] = regexp.MustCompile(`(?s)<h[3-5]><span class="mw-headline" id="` + val + `(?:.+?)(<h.>|\z)`)
	}

	tbProxyURL, err := url.Parse("socks5://proxy:9050")

	if err != nil {
		fatalf("Failed to parse proxy URL: %v\n", err)
	}

	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		fatalf("Failed to obtain proxy dialer: %v\n", err)
	}

	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client = &http.Client{Transport: tbTransport}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		searchDefinition(client, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	stdout.Flush()

}

func searchDefinition(client *http.Client, word string) {
	resp, err := client.Get(getWiktionaryUrl(word))

	if err != nil {
		fatalf("Failed to issue GET request: %v\n", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fatalf("Failed to read the body: %v\n", err)
	}

	definition := findDefinition(string(body))

	if definition == "Not Found" {
		fmt.Fprintf(os.Stderr, "%s,%s\n", word, definition)
	} else {
		stdout.Write([]string{word, definition})
	}

}

func findDefinition(html string) string {

	group := defReg.FindStringSubmatch(html)

	if len(group) != 2 {
		return "Not Found"
	}

	definition := group[1]

	for _, reg := range ignoreRegexps {
		definition = reg.ReplaceAllString(definition, "$1")
	}
	definition = deleteReg.ReplaceAllString(definition, "")

	return definition
}

func getWiktionaryUrl(word string) string {
	return fmt.Sprintf("http://en.wiktionary.org/wiki/%s", strings.Replace(word, " ", "_", -1))
}
