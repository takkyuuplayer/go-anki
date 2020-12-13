package anki

import (
	"html/template"
	"io/ioutil"

	"github.com/rakyll/statik/fs"
	// https://github.com/rakyll/statik#usage
	_ "github.com/takkyuuplayer/go-anki/anki/statik"
)

//go:generate statik -src assets -f

var tmpl = *template.New("anki")

func init() {
	mustParseAssets("entry", "/entry.html.tmpl")
}

func mustParseAssets(name, path string) *template.Template {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	r, err := statikFS.Open(path)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return template.Must(tmpl.New(name).Parse(string(contents)))
}
