package anki

type Result struct {
	Word       string
	Definition string
	IsSuccess  bool
}

type DictionaryClient interface {
    SearchDefinition(chan<- *Result, string)
}
