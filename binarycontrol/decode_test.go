package binarycontrol_test

import (
	"strings"
	"testing"

	"heartpattern.io/deb822/binarycontrol"
)

const valid = `Package: example
Source: example-src (1.0-1)
Version: 1.0-1
Architecture: amd64
Maintainer: Example <maint@example.org>
Description: example package
 paragraph
 .
  preformatted
Depends: libc6 (>= 2.36), other | alternative
Pre-Depends: pre
Recommends: recommended
Suggests: suggested
Breaks: broken
Conflicts: conflicting
Provides: virtual
Replaces: replaced
Enhances: enhanced
Installed-Size: 123
Essential: no
Multi-Arch: same
Homepage: https://example.org
Built-Using: library (= 1.0)
Section: utils
Priority: optional
X-Custom: value
 continuation
`

func TestDecode(t *testing.T) {
	for _, input := range []string{valid, strings.ReplaceAll(valid, "\n", "\r\n")} {
		file, err := binarycontrol.DecodeBytes([]byte(input))
		if err != nil {
			t.Fatal(err)
		}
		if file.Name != "example" || file.Source != "example-src (1.0-1)" || file.Version != "1.0-1" || file.Architecture != "amd64" {
			t.Fatalf("identity: %#v", file)
		}
		if file.Description != "example package\nparagraph\n.\n preformatted" {
			t.Fatalf("Description=%q", file.Description)
		}
		if file.Depends != "libc6 (>= 2.36), other | alternative" || file.PreDepends != "pre" || file.Recommends != "recommended" || file.Suggests != "suggested" || file.Breaks != "broken" || file.Conflicts != "conflicting" || file.Provides != "virtual" || file.Replaces != "replaced" || file.Enhances != "enhanced" {
			t.Fatalf("relationships: %#v", file)
		}
		if file.InstalledSize != "123" || file.Essential != "no" || file.MultiArch != "same" || file.Homepage != "https://example.org" || file.BuiltUsing != "library (= 1.0)" || file.Section != "utils" || file.Priority != "optional" {
			t.Fatalf("metadata: %#v", file)
		}
		if file.ExtraFields.Get("X-Custom") != "value\ncontinuation" || len(file.ExtraFields) != 1 {
			t.Fatalf("extra: %#v", file.ExtraFields)
		}
	}
}

func TestInvalid(t *testing.T) {
	tests := []struct{ name, input, want string }{
		{"empty", "", "missing stanza"},
		{"second stanza", valid + "\n" + valid, "exactly one stanza"},
		{"comment", "# comment\n" + valid, "field separator"},
		{"empty extra", valid + "X-Empty:\n", "empty field"},
		{"missing version", strings.Replace(valid, "Version: 1.0-1\n", "", 1), "Version"},
		{"missing package", strings.Replace(valid, "Package: example\n", "", 1), "Package"},
		{"template architecture", strings.Replace(valid, "Architecture: amd64", "Architecture: any", 1), "Architecture"},
		{"wildcard architecture", strings.Replace(valid, "Architecture: amd64", "Architecture: linux-any", 1), "Architecture"},
		{"multiple architectures", strings.Replace(valid, "Architecture: amd64", "Architecture: amd64 arm64", 1), "Architecture"},
		{"folded dependency", strings.Replace(valid, "Depends: libc6 (>= 2.36), other | alternative", "Depends: libc6,\n other", 1), "single-line"},
		{"invalid size", strings.Replace(valid, "Installed-Size: 123", "Installed-Size: -1", 1), "Installed-Size"},
		{"invalid essential", strings.Replace(valid, "Essential: no", "Essential: true", 1), "Essential"},
		{"invalid multiarch", strings.Replace(valid, "Multi-Arch: same", "Multi-Arch: invalid", 1), "Multi-Arch"},
		{"empty synopsis", strings.Replace(valid, "Description: example package", "Description:", 1), "synopsis"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := binarycontrol.Decode(strings.NewReader(tt.input))
			if file != nil || err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %#v, %v; want %q", file, err, tt.want)
			}
		})
	}
	var file *binarycontrol.File
	if err := file.Validate(); err == nil {
		t.Fatal("nil file accepted")
	}
}
