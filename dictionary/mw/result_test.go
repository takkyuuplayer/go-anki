package mw_test

import (
	"testing"

	"github.com/takkyuuplayer/go-anki/dictionary/mw"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		args string
		want string
	}{
		{args: "{bc} to move or travel to a place", want: "<b>:</b>  to move or travel to a place"},
		{args: "We {it}went{/it} many miles that day.", want: "We <i>went</i> many miles that day."},
		{args: "Something (such as a law or contract) that {phrase}goes through{/phrase} is officially accepted and approved.",
			want: "Something (such as a law or contract) that <i>goes through</i> is officially accepted and approved.",
		},
		{args: "{b}string{/b}", want: "<b>string</b>"},
		{args: "{inf}string{/inf}", want: "<sub>string</sub>"},
		{args: "{ldquo} {rdquo}", want: "&ldquo; &rdquo;"},
		{args: "{sc}string{/sc}", want: `<span style="font-variant: small-caps;">string</span>`},
		{args: "{sup}string{/sup}", want: "<sup>string</sup>"},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if got := mw.Format(tt.args); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
