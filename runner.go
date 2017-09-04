package main

import (
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

func fatalf(fmtStr string, args interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args)
	os.Exit(-1)
}

func main() {
	// Create a transport that uses Tor Browser's SocksPort.  If
	// talking to a system tor, this may be an AF_UNIX socket, or
	// 127.0.0.1:9050 instead.
	tbProxyURL, err := url.Parse("socks5://proxy:9050")
	if err != nil {
		fatalf("Failed to parse proxy URL: %v\n", err)
	}

	// Get a proxy Dialer that will create the connection on our
	// behalf via the SOCKS5 proxy.  Specify the authentication
	// and re-create the dialer/transport/client if tor's
	// IsolateSOCKSAuth is needed.
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		fatalf("Failed to obtain proxy dialer: %v\n", err)
	}

	// Make a http.Transport that uses the proxy dialer, and a
	// http.Client that uses the transport.
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	client = &http.Client{Transport: tbTransport}

	// Example: Fetch something.  Real code will probably want to use
	// client.Do() so they can change the User-Agent.
	resp, err := client.Get(getWiktionaryUrl("put up with"))
	if err != nil {
		fatalf("Failed to issue GET request: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fatalf("Failed to read the body: %v\n", err)
	}

	fmt.Println(string(body))
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
