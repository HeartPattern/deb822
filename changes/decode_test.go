package changes_test

import (
	"strings"
	"testing"

	"heartpattern.io/deb822/changes"
)

const valid = `Format: 1.8
Date: Sun, 08 Nov 2026 12:00:00 +0900
Source: example (1.0-1)
Binary: example
 example-doc
Architecture: source amd64 all
Version: 1.0-1
Distribution: unstable
Urgency: medium
Maintainer: Example <maint@example.org>
Changed-By: Other <other@example.org>
Git-Tag-Tagger: Example <maint@example.org>
Git-Tag-Info: tag=abcdef fp=123456
Description:
 example - example package
 example-doc - documentation
Closes: 123 456
Changes:
 example (1.0-1) unstable; urgency=medium
 .
   * Initial release.
Checksums-Sha1:
 0123456789012345678901234567890123456789 123 example.deb
Checksums-Sha256:
 0123456789012345678901234567890123456789012345678901234567890123 123 example.deb
Files:
 01234567890123456789012345678901 123 utils optional example.deb
X-Custom: retained
`

func TestDecode(t *testing.T) {
	file, err := changes.DecodeBytes([]byte(valid))
	if err != nil {
		t.Fatal(err)
	}
	if file.Format != "1.8" || file.Source != "example (1.0-1)" || file.Binary != "example example-doc" || file.Architecture != "source amd64 all" || file.Version != "1.0-1" {
		t.Fatalf("identity: %#v", file)
	}
	if file.Date != "Sun, 08 Nov 2026 12:00:00 +0900" || file.Distribution != "unstable" || file.Urgency != "medium" || file.Maintainer != "Example <maint@example.org>" || file.ChangedBy != "Other <other@example.org>" || file.Closes != "123 456" {
		t.Fatalf("metadata: %#v", file)
	}
	if file.GitTagInfo != "tag=abcdef fp=123456" || file.GitTagTagger != "Example <maint@example.org>" {
		t.Fatal("tag metadata lost")
	}
	if file.Description != "\nexample - example package\nexample-doc - documentation" || file.Changes != "\nexample (1.0-1) unstable; urgency=medium\n.\n  * Initial release." {
		t.Fatalf("multiline content: %#v", file)
	}
	if !strings.HasPrefix(file.ChecksumsSha1, "\n0123") || !strings.HasPrefix(file.ChecksumsSha256, "\n0123") || !strings.HasPrefix(file.Files, "\n0123") {
		t.Fatal("multiline tables changed")
	}
	if file.ExtraFields.Get("X-Custom") != "retained" || len(file.ExtraFields) != 1 {
		t.Fatalf("extra: %#v", file.ExtraFields)
	}
}

func TestSourceOnly(t *testing.T) {
	input := strings.Replace(valid, "Architecture: source amd64 all", "Architecture: source", 1)
	input = strings.Replace(input, "Binary: example\n example-doc\n", "", 1)
	input = strings.Replace(input, "Description:\n example - example package\n example-doc - documentation\n", "", 1)
	if _, err := changes.DecodeBytes([]byte(input)); err != nil {
		t.Fatal(err)
	}
	// Binary uploads require both fields, even if the source is also uploaded.
	for _, arch := range []string{"amd64", "source amd64"} {
		if _, err := changes.DecodeBytes([]byte(strings.Replace(input, "Architecture: source", "Architecture: "+arch, 1))); err == nil {
			t.Fatalf("%s upload accepted without Binary", arch)
		}
	}
}

func TestInvalid(t *testing.T) {
	for _, tt := range []struct{ name, input, want string }{
		{"missing date", strings.Replace(valid, "Date: Sun, 08 Nov 2026 12:00:00 +0900\n", "", 1), "Date"},
		{"missing binary", strings.Replace(valid, "Binary: example\n example-doc\n", "", 1), "Binary"},
		{"missing description", strings.Replace(valid, "Description:\n example - example package\n example-doc - documentation\n", "", 1), "Description"},
		{"synopsis instead of table", strings.Replace(valid, "Description:\n example - example package", "Description: example - example package", 1), "empty first line"},
		{"inline changes", strings.Replace(valid, "Changes:\n example", "Changes: example", 1), "Changes"},
		{"wildcard", strings.Replace(valid, "Architecture: source amd64 all", "Architecture: any", 1), "Architecture"},
		{"second stanza", valid + "\n" + valid, "exactly one stanza"},
		{"empty extra", valid + "X-Empty:\n", "empty field"},
		{"signed input", "-----BEGIN PGP SIGNED MESSAGE-----\nHash: SHA256\n\n" + valid, "field separator"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			file, err := changes.Decode(strings.NewReader(tt.input))
			if file != nil || err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %#v, %v; want %q", file, err, tt.want)
			}
		})
	}
	var file *changes.File
	if err := file.Validate(); err == nil {
		t.Fatal("nil file accepted")
	}
}
