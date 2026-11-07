package dsc_test

import (
	"strings"
	"testing"

	"heartpattern.io/deb822/dsc"
)

const valid = `Format: 3.0 (quilt)
Source: example
Binary: example,
 example-doc
Architecture: any all
Version: 1.0-1
Maintainer: Example <maint@example.org>
Uploaders: Other <other@example.org>
Standards-Version: 4.7.4
Homepage: https://example.org
Vcs-Git: https://example.org/repo -b debian
Vcs-Browser: https://example.org/repo
Testsuite: autopkgtest
Dgit: abcdef
 extra-data
Git-Tag-Tagger: Example <maint@example.org>
Git-Tag-Info: tag=abcdef fp=123456
Build-Depends: debhelper-compat (= 13)
Build-Depends-Arch: compiler
Build-Depends-Indep: docs
Build-Conflicts: old-tools
Build-Conflicts-Arch: old-compiler
Build-Conflicts-Indep: old-docs
Package-List:
 example deb utils optional arch=any
 example-doc deb doc optional arch=all
Checksums-Sha1:
 0123456789012345678901234567890123456789 123 example.tar.xz
Checksums-Sha256:
 0123456789012345678901234567890123456789012345678901234567890123 123 example.tar.xz
Files:
 01234567890123456789012345678901 123 example.tar.xz
X-Custom: retained
`

func TestDecode(t *testing.T) {
	file, err := dsc.DecodeBytes([]byte(valid))
	if err != nil {
		t.Fatal(err)
	}
	if file.Format != "3.0 (quilt)" || file.Source != "example" || file.Binary != "example, example-doc" || file.Architecture != "any all" || file.Version != "1.0-1" {
		t.Fatalf("identity: %#v", file)
	}
	if file.Uploaders != "Other <other@example.org>" || file.Homepage != "https://example.org" || file.Testsuite != "autopkgtest" || file.Dgit != "abcdef extra-data" || file.GitTagInfo != "tag=abcdef fp=123456" || file.GitTagTagger != "Example <maint@example.org>" {
		t.Fatalf("metadata: %#v", file)
	}
	if file.VcsGit != "https://example.org/repo -b debian" || file.VcsBrowser != "https://example.org/repo" {
		t.Fatalf("VCS: %#v", file)
	}
	if file.BuildDepends != "debhelper-compat (= 13)" || file.BuildDependsArch != "compiler" || file.BuildDependsIndep != "docs" || file.BuildConflicts != "old-tools" || file.BuildConflictsArch != "old-compiler" || file.BuildConflictsIndep != "old-docs" {
		t.Fatalf("relationships: %#v", file)
	}
	if file.PackageList != "\nexample deb utils optional arch=any\nexample-doc deb doc optional arch=all" {
		t.Fatalf("PackageList=%q", file.PackageList)
	}
	if !strings.HasPrefix(file.ChecksumsSha1, "\n0123") || !strings.HasPrefix(file.ChecksumsSha256, "\n0123") || !strings.HasPrefix(file.Files, "\n0123") {
		t.Fatal("multiline tables changed")
	}
	if file.ExtraFields.Get("X-Custom") != "retained" || len(file.ExtraFields) != 1 {
		t.Fatalf("extra: %#v", file.ExtraFields)
	}
}

func TestInvalid(t *testing.T) {
	for _, tt := range []struct{ name, input, want string }{
		{"missing format", strings.Replace(valid, "Format: 3.0 (quilt)\n", "", 1), "Format"},
		{"missing standards", strings.Replace(valid, "Standards-Version: 4.7.4\n", "", 1), "Standards-Version"},
		{"missing sha1", strings.Replace(valid, "Checksums-Sha1:\n 0123456789012345678901234567890123456789 123 example.tar.xz\n", "", 1), "Checksums-Sha1"},
		{"table first line", strings.Replace(valid, "Files:\n ", "Files: ", 1), "empty first line"},
		{"bad any list", strings.Replace(valid, "Architecture: any all", "Architecture: any amd64", 1), "Architecture"},
		{"source architecture", strings.Replace(valid, "Architecture: any all", "Architecture: source", 1), "Architecture"},
		{"multiple vcs", valid + "Vcs-Hg: repo\n", "multiple Vcs"},
		{"folded uploaders", strings.Replace(valid, "Uploaders: Other <other@example.org>", "Uploaders:\n Other <other@example.org>", 1), "single-line"},
		{"signed input", "-----BEGIN PGP SIGNED MESSAGE-----\nHash: SHA256\n\n" + valid, "field separator"},
		{"second stanza", valid + "\n" + valid, "exactly one stanza"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			file, err := dsc.Decode(strings.NewReader(tt.input))
			if file != nil || err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %#v, %v; want %q", file, err, tt.want)
			}
		})
	}
	var file *dsc.File
	if err := file.Validate(); err == nil {
		t.Fatal("nil file accepted")
	}
}

func TestUdebException(t *testing.T) {
	input := strings.Replace(valid, "Standards-Version: 4.7.4\n", "", 1)
	input = strings.ReplaceAll(input, " deb ", " udeb ")
	if _, err := dsc.DecodeBytes([]byte(input)); err != nil {
		t.Fatal(err)
	}
	input = strings.Replace(input, " udeb ", " deb ", 1)
	if _, err := dsc.DecodeBytes([]byte(input)); err == nil {
		t.Fatal("mixed packages require Standards-Version")
	}
}
