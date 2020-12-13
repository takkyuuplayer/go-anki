package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/takkyuuplayer/go-anki/anki"
	"github.com/takkyuuplayer/go-anki/dictionary"

	"github.com/rakyll/statik/fs"
	"github.com/takkyuuplayer/go-anki/dictionary/mw"

	// https://github.com/rakyll/statik#usage
	_ "github.com/takkyuuplayer/go-anki/web/statik"
)

const Concurrency = 10

var dictionaries = map[string]dictionary.Dictionary{
	"mw": mw.NewLearners(os.Getenv("MW_LEARNERS_KEY"), &http.Client{}),
}

func post(w http.ResponseWriter, r *http.Request) {
	dictionaryApi := dictionaries[r.PostFormValue("dictionary")]

	w.Header().Set("Content-Disposition", "attachment; filename=anki.tsv")
	in := strings.NewReader(r.PostFormValue("words"))

	out := csv.NewWriter(w)
	out.Comma = '\t'

	anki.Run(dictionaryApi, in, out, out)
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
