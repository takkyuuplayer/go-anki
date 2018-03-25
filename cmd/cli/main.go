package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
		HttpClient: &http.Client{},
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
