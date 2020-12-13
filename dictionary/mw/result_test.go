package mw

import "testing"

func TestFormat(t *testing.T) {
	tests := []struct {
		args string
		want string
	}{
		{args: "{bc} to move or travel to a place", want: "to move or travel to a place"},
		{args: "We {it}went{/it} many miles that day.", want: "We <i>went</i> many miles that day."},
		{args: "Something (such as a law or contract) that {phrase}goes through{/phrase} is officially accepted and approved.",
			want: "Something (such as a law or contract) that <i>goes through</i> is officially accepted and approved.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if got := Format(tt.args); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
