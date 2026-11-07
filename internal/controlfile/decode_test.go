package controlfile

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type errorReader struct{ err error }

func (r errorReader) Read([]byte) (int, error) { return 0, r.err }

func TestDecodeReadErrors(t *testing.T) {
	want := errors.New("read failed")
	for _, prefix := range []string{"", "Field: value\n\n"} {
		var result struct {
			Field string `control:"Field"`
		}
		err := Decode(io.MultiReader(strings.NewReader(prefix), errorReader{want}), &result)
		if !errors.Is(err, want) {
			t.Fatalf("prefix %q: %v", prefix, err)
		}
	}
}

func TestDecodeTrailingSeparators(t *testing.T) {
	var result struct {
		Field string `control:"Field"`
	}
	if err := Decode(strings.NewReader("\nField: value\n\n\n"), &result); err != nil {
		t.Fatal(err)
	}
	if result.Field != "value" {
		t.Fatalf("got %#v", result)
	}
}
