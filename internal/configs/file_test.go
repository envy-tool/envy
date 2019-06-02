package configs

import (
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	t.Run("command", func (t *testing.T) {
		f, diags := LoadConfigFile("testdata/command.nv.hcl")
		if diags.HasErrors() {
			for _, diag := range diags {
				t.Errorf("unexpected diagnostic: %s", diag)
			}
			return
		}

		if got, want := len(f.Commands), 1; got != want {
			t.Errorf("wrong number of commands %d; want %d", got, want)
			if got == 0 {
				return
			}
		}

		cmd := f.Commands[0]
		if got, want := cmd.Name, "terraform"; got != want {
			t.Errorf("wrong name %q; want %q", got, want)
		}
	})
}
