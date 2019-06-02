package configs

import (
	"testing"

	"github.com/hashicorp/hcl2/gohcl"
)

func TestLoadConfigFile(t *testing.T) {
	t.Run("command", func(t *testing.T) {
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
	t.Run("helper", func(t *testing.T) {
		f, diags := LoadConfigFile("testdata/helper.nv.hcl")
		if diags.HasErrors() {
			for _, diag := range diags {
				t.Errorf("unexpected diagnostic: %s", diag)
			}
			return
		}

		if got, want := len(f.Helpers), 1; got != want {
			t.Errorf("wrong number of helpers %d; want %d", got, want)
			if got == 0 {
				return
			}
		}

		helper := f.Helpers[0]
		if got, want := helper.Type, "foo"; got != want {
			t.Errorf("wrong type %q; want %q", got, want)
		}
		if got, want := helper.Name, "bar"; got != want {
			t.Errorf("wrong name %q; want %q", got, want)
		}

		type BodyContent struct {
			Foo string `hcl:"foo"`
		}
		var content BodyContent
		diags = gohcl.DecodeBody(helper.Body, nil, &content)
		if got, want := content.Foo, "bar"; got != want {
			t.Errorf("wrong \"foo\" value %q; want %q", got, want)
		}
	})
	t.Run("shared", func(t *testing.T) {
		f, diags := LoadConfigFile("testdata/shared_object.nv.hcl")
		if diags.HasErrors() {
			for _, diag := range diags {
				t.Errorf("unexpected diagnostic: %s", diag)
			}
			return
		}

		if got, want := len(f.SharedObjects), 1; got != want {
			t.Errorf("wrong number of shared objects %d; want %d", got, want)
			if got == 0 {
				return
			}
		}

		so := f.SharedObjects[0]
		if got, want := so.Name, "example"; got != want {
			t.Errorf("wrong name %q; want %q", got, want)
		}
	})
}
