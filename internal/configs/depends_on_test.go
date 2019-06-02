package configs

import (
	"testing"

	"github.com/envy-tool/envy/internal/addrs"
)

func TestDecodeReference(t *testing.T) {
	tests := []struct {
		str       string
		want      addrs.Referenceable
		remainLen int
	}{
		{
			`foo.bar`,
			addrs.MakeHelper("foo", "bar"),
			0,
		},
		{
			`foo.bar.baz`,
			addrs.MakeHelper("foo", "bar"),
			1,
		},
		{
			`command.foo`,
			addrs.MakeCommand("foo"),
			0,
		},
		{
			`command.foo.bar`,
			addrs.MakeCommand("foo"),
			1,
		},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			ref, remain, diags := ParseReferenceStr(test.str)
			if len(diags) > 0 {
				for _, diag := range diags {
					t.Errorf("Unexpected diagnostic: %s", diag)
				}
				return
			}

			if got, want := ref.Addr, test.want; got != want {
				t.Errorf("wrong address\ngot:  %#v\nwant: %#v", got, want)
			}
			if got, want := len(remain), test.remainLen; got != want {
				t.Errorf("wrong remain length %d; want %d", got, want)
			}
		})
	}
}
