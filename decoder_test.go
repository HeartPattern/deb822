package deb822

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestDecoderStanzas(t *testing.T) {
	d := NewDecoderBytes([]byte("\n \t\nPackage: first\nVersion: 1\n\n\nPackage: second\r\nDescription: summary\r\n first line\r\n .\r\n  indented\r\nEmpty:\r\n"))
	for i, want := range []Stanza{
		{"package": "first", "version": "1"},
		{"package": "second", "description": "summary\nfirst line\n.\n indented", "empty": ""},
	} {
		got, err := d.Decode()
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("stanza %d = %#v, %v; want %#v", i, got, err, want)
		}
	}
	for range 2 {
		if got, err := d.Decode(); got != nil || err != io.EOF {
			t.Fatalf("after last stanza = %#v, %v", got, err)
		}
	}
}

func TestDecoderEOF(t *testing.T) {
	for _, input := range []string{"", "\n\n", " \t"} {
		if got, err := NewDecoder(strings.NewReader(input)).Decode(); got != nil || err != io.EOF {
			t.Errorf("%q = %#v, %v", input, got, err)
		}
	}
	for _, ending := range []string{"", "\n", "\n\n"} {
		d := NewDecoder(strings.NewReader("Package: hello" + ending))
		if got, err := d.Decode(); err != nil || !reflect.DeepEqual(got, Stanza{"package": "hello"}) {
			t.Fatalf("ending %q = %#v, %v", ending, got, err)
		}
		if got, err := d.Decode(); got != nil || err != io.EOF {
			t.Fatalf("ending %q after stanza = %#v, %v", ending, got, err)
		}
	}
}

func TestDecoderInvalidStanza(t *testing.T) {
	for _, input := range []string{
		"Package: one\nPackage: two", "Package: one\npACKAGE: two",
		"Empty:\nOther: x\nEMPTY: y", " orphan", "\n orphan", "Package: one\nbroken",
		"# comment\nPackage: one",
	} {
		if got, err := NewDecoder(strings.NewReader(input)).Decode(); got != nil || err == nil || err == io.EOF {
			t.Errorf("%q = %#v, %v; want syntax error", input, got, err)
		}
	}
	d := NewDecoder(strings.NewReader("Package: one\n\n orphan"))
	if _, err := d.Decode(); err != nil {
		t.Fatal(err)
	}
	if got, err := d.Decode(); got != nil || err == nil {
		t.Fatalf("continuation after separator = %#v, %v", got, err)
	}
}

func TestDecoderComments(t *testing.T) {
	d := NewDecoder(strings.NewReader("# header\nDescription: summary # literal\n first\n# between\n # indented\n# end"))
	d.AllowComments = true
	want := Stanza{"description": "summary # literal\nfirst\n# indented"}
	if got, err := d.Decode(); err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("Decode() = %#v, %v", got, err)
	}
	if got, err := d.Decode(); got != nil || err != io.EOF {
		t.Fatalf("after stanza = %#v, %v", got, err)
	}
	d = NewDecoder(strings.NewReader("# only comment"))
	d.AllowComments = true
	if got, err := d.Decode(); got != nil || err != io.EOF {
		t.Fatalf("comments only = %#v, %v", got, err)
	}
}

func TestDecoderReadError(t *testing.T) {
	want := errors.New("read failed")
	if got, err := NewDecoder(failingReader{want}).Decode(); got != nil || !errors.Is(err, want) {
		t.Fatalf("Decode() = %#v, %v", got, err)
	}
}

func TestDecoderIterator(t *testing.T) {
	for _, tt := range []struct {
		input         string
		stanzas, errs int
	}{
		{"", 0, 0}, {"Package: one\n\nPackage: two", 2, 0}, {"Package: one\n\nbroken\n\nPackage: three", 1, 1},
	} {
		stanzas, errs := 0, 0
		NewDecoder(strings.NewReader(tt.input)).Decodes()(func(s Stanza, err error) bool {
			if err != nil {
				errs++
			} else {
				stanzas++
			}
			return stanzas+errs < 10
		})
		if stanzas != tt.stanzas || errs != tt.errs {
			t.Errorf("%q: got %d stanzas, %d errors", tt.input, stanzas, errs)
		}
	}
	d := NewDecoder(strings.NewReader("Package: one\n\nPackage: two"))
	for range d.Decodes() {
		break
	}
	if got, err := d.Decode(); err != nil || got.Get("Package") != "two" {
		t.Fatalf("after early stop = %#v, %v", got, err)
	}
}
