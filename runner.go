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

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

func main() {
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

	stdout.Write([]string{word, definition})
}

func findDefinition(html string) string {

	group := defReg.FindStringSubmatch(html)

	if len(group) == 2 {
		return group[1]
	}

	return "Not Found"
}

func getWiktionaryUrl(word string) string {
	return fmt.Sprintf("http://en.wiktionary.org/wiki/%s", strings.Replace(word, " ", "_", -1))
}
