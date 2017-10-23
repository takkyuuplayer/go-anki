package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"os"

	"github.com/takkyuuplayer/go-anki"
	"github.com/takkyuuplayer/go-anki/wiktionary"
)

var stdout = csv.NewWriter(os.Stdout)

const parallel = 10

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

	counter := 0
	scanner := bufio.NewScanner(os.Stdin)
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client := &http.Client{Transport: tbTransport}

	ch := make(chan *anki.Result)
	wc := wiktionary.Client{
		HttpClient: client,
	}

	for ; counter < parallel; counter++ {
		if scanner.Scan() {
			go wc.SearchDefinition(ch, scanner.Text())
		} else {
			break
		}
	}

	for i := 0; i < counter; i++ {
		result := <-ch
		if result.IsSuccess {
			stdout.Write([]string{result.Word, result.Definition})
		} else {
			fmt.Fprintf(os.Stderr, "%s,%s\n", result.Word, result.Definition)
		}

		if scanner.Scan() {
			go wc.SearchDefinition(ch, scanner.Text())
			counter++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	stdout.Flush()
}
