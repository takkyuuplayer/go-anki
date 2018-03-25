package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/proxy"

	"github.com/takkyuuplayer/go-anki"
	"github.com/takkyuuplayer/go-anki/mw"
	"github.com/takkyuuplayer/go-anki/wiktionary"
)

var dictionaries = map[string]anki.Dictionary{
	"mw":         mw.New(os.Getenv("MW_API_KEY"), "learners"),
	"wiktionary": wiktionary.New(),
}

func main() {
	var (
		dictionary = flag.String("dictionary", "mw", "dictionary to use. (mw|wiktionary)")
	)
	flag.Parse()

	wc := &anki.Client{
		Dictionary: dictionaries[*dictionary],
		HttpClient: httpClient(),
	}

	if wc.Dictionary == nil {
		log.Printf("Unknown Dictionary: %s", *dictionary)
		return
	}

	out := csv.NewWriter(os.Stdout)
	out.Comma = '\t'

	outErr := csv.NewWriter(os.Stderr)
	outErr.Comma = '\t'

	wc.Run(os.Stdin, out, outErr)
}

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

func httpClient() *http.Client {
	tbProxyURL, err := url.Parse("socks5://proxy:9050")

	if err != nil {
		fatalf("Failed to parse proxy URL: %v\n", err)
	}

	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		fatalf("Failed to obtain proxy dialer: %v\n", err)
	}

	tbTransport := &http.Transport{Dial: tbDialer.Dial}

	return &http.Client{Transport: tbTransport}
}
