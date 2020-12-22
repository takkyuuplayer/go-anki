package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/takkyuuplayer/go-anki/dictionary/eijiro"

	"github.com/takkyuuplayer/go-anki/anki"
	"github.com/takkyuuplayer/go-anki/dictionary"

	"github.com/rakyll/statik/fs"
	"github.com/takkyuuplayer/go-anki/dictionary/mw"

	// https://github.com/rakyll/statik#usage
	_ "github.com/takkyuuplayer/go-anki/web/statik"
)

var dictionaries = map[string]dictionary.Dictionary{
	"mw":     mw.NewLearners(os.Getenv("MW_LEARNERS_KEY"), http.DefaultClient),
	"eijiro": eijiro.NewEijiro(http.DefaultClient),
}

func post(w http.ResponseWriter, r *http.Request) {
	dic := dictionaries[r.PostFormValue("dictionary")]
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Disposition", "attachment; filename=anki.tsv")

	in := strings.NewReader(r.PostFormValue("words"))

	out := csv.NewWriter(w)
	out.Comma = '\t'

	outErr := new(bytes.Buffer)
	outErrCsv := csv.NewWriter(outErr)
	outErrCsv.Comma = '\t'

	anki.Run(dic, in, out, outErrCsv)

	outErr.WriteTo(w)
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
