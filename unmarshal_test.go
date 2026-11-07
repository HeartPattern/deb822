package deb822_test

import (
	"reflect"
	"strings"
	"testing"

	"heartpattern.io/deb822"
)

func TestUnmarshalStanza(t *testing.T) {
	// Arbitrary Go and control names verify the mapper has no schema-specific logic.
	var target struct {
		Extra     deb822.Stanza `control:",remain"`
		Name      string        `control:"X-Name"`
		Simple    string        `control:"X-Simple,simple"`
		Folded    string        `control:"X-Folded,folded"`
		Multiline string        `control:"X-Lines,multiline"`
		Ignored   int           `control:"-"`
		Untagged  int
	}
	stanza := deb822.Stanza{
		"x-name": " name \t", "x-simple": "\t simple ",
		"x-folded": " first\n\t second  third ",
		"x-lines":  " synopsis  \n.\n indented \n", "x-extra": " raw\n. ", "x-empty": " \t\n",
	}
	original := make(deb822.Stanza)
	for name, value := range stanza {
		original[name] = value
	}
	err := deb822.UnmarshalStanza(stanza, &target)
	if err != nil {
		t.Fatal(err)
	}
	if target.Name != "name" || target.Simple != "simple" || target.Folded != "first second third" || target.Multiline != " synopsis  \n.\n indented \n" {
		t.Fatalf("mapped values: %#v", target)
	}
	if !reflect.DeepEqual(target.Extra, deb822.Stanza{"x-extra": " raw\n. ", "x-empty": " \t\n"}) {
		t.Fatalf("extra fields: %#v", target.Extra)
	}
	if !reflect.DeepEqual(stanza, original) {
		t.Fatal("input stanza was modified")
	}
	target.Extra["x-extra"] = "changed"
	if stanza["x-extra"] != " raw\n. " {
		t.Fatal("extra fields alias input stanza")
	}
}

func TestUnmarshalStanzaErrors(t *testing.T) {
	tests := []struct {
		name   string
		target any
		want   string
	}{
		{"nil", nil, "non-nil pointer"},
		{"nil pointer", (*struct{})(nil), "non-nil pointer"},
		{"value", struct{}{}, "non-nil pointer"},
		{"pointer to scalar", new(string), "non-nil pointer"},
		{"unknown format", &struct {
			Value string `control:"Value,unknown"`
		}{}, "unsupported control format"},
		{"multiple formats", &struct {
			Value string `control:"Value,simple,folded"`
		}{}, "unsupported control format"},
		{"missing name", &struct {
			Value string `control:",folded"`
		}{}, "invalid control tag"},
		{"empty tag", &struct {
			Value string `control:""`
		}{}, "invalid control tag"},
		{"wrong type", &struct {
			Value int `control:"Value"`
		}{}, "exported string"},
		{"unexported", &struct {
			value string `control:"Value"`
		}{}, "exported string"},
		{"duplicate", &struct {
			First  string `control:"Value"`
			Second string `control:"vALUE"`
		}{}, "duplicate control tag"},
		{"named remain", &struct {
			Extra deb822.Stanza `control:"Extra,remain"`
		}{}, "must not specify a field name"},
		{"multiple remain", &struct {
			First  deb822.Stanza `control:",remain"`
			Second deb822.Stanza `control:",remain"`
		}{}, "multiple remain fields"},
		{"wrong remain type", &struct {
			Extra map[string]string `control:",remain"`
		}{}, "exported deb822.Stanza"},
		{"unexported remain", &struct {
			extra deb822.Stanza `control:",remain"`
		}{}, "exported deb822.Stanza"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := deb822.UnmarshalStanza(nil, tt.target)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %v; want %q", err, tt.want)
			}
		})
	}
	var target struct {
		Value string `control:"Value"`
	}
	err := deb822.UnmarshalStanza(deb822.Stanza{"value": "first\nsecond"}, &target)
	if err == nil || !strings.Contains(err.Error(), "Value: must be a single-line field") {
		t.Fatalf("simple multiline value: %v", err)
	}
}

func TestUnmarshalStanzaWithoutRemain(t *testing.T) {
	var target struct {
		Value   string        `control:"Value"`
		Ignored deb822.Stanza `control:"-"`
	}
	if err := deb822.UnmarshalStanza(deb822.Stanza{"value": "known", "other": "unknown"}, &target); err != nil {
		t.Fatal(err)
	}
	if target.Value != "known" || target.Ignored != nil {
		t.Fatalf("mapped values: %#v", target)
	}
}

func TestUnmarshalStanzaEmptyRemain(t *testing.T) {
	target := struct {
		Value string        `control:"Value"`
		Extra deb822.Stanza `control:",remain"`
	}{Extra: deb822.Stanza{"stale": "old value"}}
	if err := deb822.UnmarshalStanza(deb822.Stanza{"value": "known"}, &target); err != nil {
		t.Fatal(err)
	}
	if target.Value != "known" || len(target.Extra) != 0 {
		t.Fatalf("mapped values: %#v", target)
	}
}

func TestUnmarshalStanzaEmptyAndAbsentFields(t *testing.T) {
	target := struct {
		Empty      string        `control:"Empty"`
		Whitespace string        `control:"Whitespace"`
		Folded     string        `control:"Folded,folded"`
		Multiline  string        `control:"Multiline,multiline"`
		Absent     string        `control:"Absent"`
		Extra      deb822.Stanza `control:",remain"`
	}{Empty: "old", Whitespace: "old", Folded: "old", Multiline: "old", Absent: "keep"}
	stanza := deb822.Stanza{
		"Empty": "", "Whitespace": " \t", "Folded": " \t\n",
		"Multiline": " \t\n", "Unknown": "", "Unknown-Spaces": " \t",
	}
	if err := deb822.UnmarshalStanza(stanza, &target); err != nil {
		t.Fatal(err)
	}
	if target.Empty != "" || target.Whitespace != "" || target.Folded != "" || target.Multiline != " \t\n" || target.Absent != "keep" {
		t.Fatalf("empty and absent values: %#v", target)
	}
	if !reflect.DeepEqual(target.Extra, deb822.Stanza{"unknown": "", "unknown-spaces": " \t"}) {
		t.Fatalf("remaining empty fields: %#v", target.Extra)
	}
}

func TestUnmarshalStanzaCaseInsensitive(t *testing.T) {
	var target struct {
		Value string `control:"VALUE"`
	}
	for _, name := range []string{"value", "Value", "vAlUe"} {
		if err := deb822.UnmarshalStanza(deb822.Stanza{name: "matched"}, &target); err != nil {
			t.Fatal(err)
		}
		if target.Value != "matched" {
			t.Fatalf("%q was not matched", name)
		}
		target.Value = ""
	}
	if err := deb822.UnmarshalStanza(deb822.Stanza{"Value": "", "value": "second"}, &target); err == nil || !strings.Contains(err.Error(), "duplicate field") {
		t.Fatalf("duplicate case variants: %v", err)
	}
}
