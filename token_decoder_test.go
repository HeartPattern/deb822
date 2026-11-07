package deb822

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name, input string
		want        Token
	}{
		{"empty line", "\n", &Separator{}},
		{"whitespace line", " \t\n", &Separator{}},
		{"CRLF separator", "\r\n", &Separator{}},
		{"final whitespace", " \t", &Separator{}},
		{"field", "Package: hello\n", &Field{Name: "Package", Value: "hello"}},
		{"CRLF field", "Package: hello\r\n", &Field{Name: "Package", Value: "hello"}},
		{"final field", "Package: hello", &Field{Name: "Package", Value: "hello"}},
		{"colon in value", "URL: https://example.org\n", &Field{Name: "URL", Value: "https://example.org"}},
		{"continuation", "  indented  \n", &Continuation{Value: " indented  "}},
		{"tab continuation", "\ttext\r\n", &Continuation{Value: "text"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewTokenDecoder(strings.NewReader(tt.input))
			got, err := d.Decode()
			if err != nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Decode() = %#v, %v; want %#v", got, err, tt.want)
			}
			if got, err := d.Decode(); got != nil || err != io.EOF {
				t.Fatalf("after final line = %#v, %v", got, err)
			}
		})
	}
}

func TestDecodeRejectsInvalidFields(t *testing.T) {
	for _, input := range []string{"#Bad: value\n", "-Bad: value\n", ": value\n", "Bad Name: value\n", "Náme: value\n", "missing colon\n"} {
		if token, err := NewTokenDecoder(strings.NewReader(input)).Decode(); err == nil || token != nil {
			t.Errorf("%q: got %#v, %v", input, token, err)
		}
	}
}

func TestDecodesTerminates(t *testing.T) {
	for _, tt := range []struct {
		input        string
		tokens, errs int
	}{
		{"", 0, 0}, {"Package: hello\n\nPackage: world", 3, 0}, {"invalid\nPackage: ignored\n", 0, 1},
	} {
		tokens, errs := 0, 0
		NewTokenDecoder(strings.NewReader(tt.input)).Decodes()(func(token Token, err error) bool {
			if err != nil {
				errs++
			} else {
				tokens++
			}
			// Bound a broken iterator so a regression fails instead of hanging.
			return tokens+errs < 10
		})
		if tokens != tt.tokens || errs != tt.errs {
			t.Errorf("%q: tokens=%d errors=%d; want %d, %d", tt.input, tokens, errs, tt.tokens, tt.errs)
		}
	}
}

func TestDecodesEarlyStop(t *testing.T) {
	d := NewTokenDecoder(strings.NewReader("First: one\nSecond: two\n"))
	for range d.Decodes() {
		break
	}
	got, err := d.Decode()
	if err != nil || !reflect.DeepEqual(got, &Field{Name: "Second", Value: "two"}) {
		t.Fatalf("after early stop = %#v, %v", got, err)
	}
}

type failingReader struct{ err error }

func (r failingReader) Read(p []byte) (int, error) { return copy(p, "Name: partial"), r.err }
func TestDecodePropagatesReadError(t *testing.T) {
	want := errors.New("read failed")
	got, err := NewTokenDecoder(failingReader{want}).Decode()
	if got != nil || !errors.Is(err, want) {
		t.Fatalf("Decode() = %#v, %v", got, err)
	}
}

func TestDecodeComments(t *testing.T) {
	tests := []struct {
		name, input string
		want        []Token
	}{
		{"comments only", "# comment\n# last", nil},
		{"empty comment", "#", nil},
		{"before field", "# comment\r\n# another\nPackage: hello", []Token{&Field{Name: "Package", Value: "hello"}}},
		{"between continuations", "Description: summary\n first\n# comment\n second\n# last", []Token{&Field{Name: "Description", Value: "summary"}, &Continuation{Value: "first"}, &Continuation{Value: "second"}}},
		{"separator preserved", "# first\n\n# last\n", []Token{&Separator{}}},
		{"hash in values", "Field: value # text\n # text\n\t# tab\n", []Token{&Field{Name: "Field", Value: "value # text"}, &Continuation{Value: "# text"}, &Continuation{Value: "# tab"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewTokenDecoder(strings.NewReader(tt.input))
			d.AllowComments = true
			for i, want := range tt.want {
				got, err := d.Decode()
				if err != nil || !reflect.DeepEqual(got, want) {
					t.Fatalf("token %d = %#v, %v; want %#v", i, got, err, want)
				}
			}
			if got, err := d.Decode(); got != nil || err != io.EOF {
				t.Fatalf("after comments = %#v, %v; want EOF", got, err)
			}
		})
	}
}

func TestCommentsDisabledByDefault(t *testing.T) {
	for _, input := range []string{"# comment\n", "#Field: value\n"} {
		if got, err := NewTokenDecoder(strings.NewReader(input)).Decode(); got != nil || err == nil {
			t.Errorf("%q: got %#v, %v; want error", input, got, err)
		}
	}
}
