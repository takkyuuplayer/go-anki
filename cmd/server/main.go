package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rakyll/statik/fs"
	anki "github.com/takkyuuplayer/go-anki"
	"github.com/takkyuuplayer/go-anki/mw"
	_ "github.com/takkyuuplayer/go-anki/web/statik"
	"github.com/takkyuuplayer/go-anki/wiktionary"
)

var dictionaries = map[string]anki.Dictionary{
	"mw":         mw.New(os.Getenv("MW_API_KEY"), "learners"),
	"wiktionary": wiktionary.New(),
}

func post(w http.ResponseWriter, r *http.Request) {
	wc := &anki.Client{
		Dictionary: dictionaries[r.PostFormValue("dictionary")],
		HttpClient: &http.Client{},
	}

	if wc.Dictionary == nil {
		log.Printf("Unknown Dictionary: %s", r.PostFormValue("dictionary"))
		http.Redirect(w, r, "/", 301)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=anki.tsv")

	in := strings.NewReader(r.PostFormValue("words"))

	out := csv.NewWriter(w)
	out.Comma = '\t'

	wc.Run(in, out, out)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		post(w, r)
		return
	}

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	res, err := statikFS.Open("/index.html")
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(w, res)
}

func main() {
	addr := flag.String("addr", ":8080", "addr to bind")

	flag.Parse()

	http.HandleFunc("/", handler)

	log.Printf("start listening on %s", *addr)

	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Print(err)
	}
}
