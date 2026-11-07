package debiancontrol

import (
	"errors"
	"io"
	"strings"
	"testing"
)

const sourceFixture = "Source: example\nMaintainer: Example Maintainer <maint@example.org>\nStandards-Version: 4.7.4\n"
const binaryFixture = "Package: example\nArchitecture: any\nDescription: example package\n"

func TestDecode(t *testing.T) {
	input := `# source comment
sOuRcE: example
Maintainer: Example Maintainer <maint@example.org>
Uploaders: First <first@example.org>,
 Second <second@example.org>
Standards-Version: 4.7.4
Section: utils
Priority: optional
Build-Depends:
 debhelper-compat (= 13),
# comment between continuations
 pkg-config [linux-any] <!nocheck>
Build-Depends-Arch: compiler
Build-Depends-Indep: docs
Build-Conflicts: old-tools
Build-Conflicts-Arch: old-compiler
Build-Conflicts-Indep: old-docs
Homepage: https://example.org
Vcs-Git: https://example.org/repo.git -b debian [pkg]
Vcs-Browser: https://example.org/repo
Testsuite: autopkgtest
Rules-Requires-Root: no
X-Custom: keep
 continuation
X-Empty:

Package: example
Architecture: linux-any arm64
Essential: no
Depends: ${shlibs:Depends},
 ${misc:Depends}
Pre-Depends: pre
Recommends: recommended
Suggests: suggested
Breaks: broken
Conflicts: conflicting
Provides: virtual
Replaces: replaced
Enhances: enhanced
Built-Using: library (= 1.0)
Build-Profiles: <!stage1>
 <cross>
Multi-Arch: foreign
Description: example package
 first paragraph
 .
  preformatted text
X-Binary: unchanged

Package: example-doc
Architecture: all
Section: doc
Priority: optional
Homepage: https://example.org/docs
Package-Type: deb
Description: documentation`
	file, err := DecodeBytes([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	s := file.Source
	if s.Name != "example" || s.Maintainer != "Example Maintainer <maint@example.org>" || s.StandardsVersion != "4.7.4" {
		t.Fatalf("source identity: %#v", s)
	}
	if s.Uploaders != "First <first@example.org>, Second <second@example.org>" {
		t.Fatalf("Uploaders = %q", s.Uploaders)
	}
	if s.BuildDepends != "debhelper-compat (= 13), pkg-config [linux-any] <!nocheck>" {
		t.Fatalf("BuildDepends = %q", s.BuildDepends)
	}
	if s.BuildDependsArch != "compiler" || s.BuildDependsIndep != "docs" || s.BuildConflicts != "old-tools" || s.BuildConflictsArch != "old-compiler" || s.BuildConflictsIndep != "old-docs" {
		t.Fatalf("build relationships: %#v", s)
	}
	if s.VcsGit != "https://example.org/repo.git -b debian [pkg]" || s.VcsBrowser != "https://example.org/repo" || s.Homepage != "https://example.org" || s.Testsuite != "autopkgtest" || s.RulesRequiresRoot != "no" {
		t.Fatalf("source metadata: %#v", s)
	}
	if s.ExtraFields.Get("X-Custom") != "keep\ncontinuation" || len(s.ExtraFields) != 1 {
		t.Fatalf("extra fields: %#v", s.ExtraFields)
	}
	if len(file.Binaries) != 2 {
		t.Fatalf("binaries: %#v", file.Binaries)
	}
	b := file.Binaries[0]
	if b.Depends != "${shlibs:Depends}, ${misc:Depends}" || b.Description != "example package\nfirst paragraph\n.\n preformatted text" {
		t.Fatalf("binary content: %#v", b)
	}
	if b.PreDepends != "pre" || b.Recommends != "recommended" || b.Suggests != "suggested" || b.Breaks != "broken" || b.Conflicts != "conflicting" || b.Provides != "virtual" || b.Replaces != "replaced" || b.Enhances != "enhanced" || b.BuiltUsing != "library (= 1.0)" || b.BuildProfiles != "<!stage1> <cross>" {
		t.Fatalf("binary relationships: %#v", b)
	}
	if b.Section != "" || b.Priority != "" || s.Section != "utils" || s.Priority != "optional" {
		t.Fatal("source defaults must remain separate")
	}
	if b.ExtraFields.Get("X-Binary") != "unchanged" || b.MultiArch != "foreign" || b.Essential != "no" {
		t.Fatalf("binary metadata: %#v", b)
	}
	second := file.Binaries[1]
	if second.Name != "example-doc" || second.Section != "doc" || second.Priority != "optional" || second.Homepage != "https://example.org/docs" || second.PackageType != "deb" {
		t.Fatalf("second binary: %#v", second)
	}
}

func TestDecodeInvalid(t *testing.T) {
	valid := sourceFixture + "\n" + binaryFixture
	tests := []struct{ name, input, contains string }{
		{"empty", "", "source stanza"},
		{"source only", sourceFixture, "binary stanza"},
		{"binary first", binaryFixture, "source stanza: Source:"},
		{"source later", sourceFixture + "\n" + sourceFixture, "stanza 2: Package:"},
		{"missing maintainer", strings.Replace(valid, "Maintainer: Example Maintainer <maint@example.org>\n", "", 1), "Maintainer"},
		{"missing standards", strings.Replace(valid, "Standards-Version: 4.7.4\n", "", 1), "Standards-Version"},
		{"bad standards", strings.Replace(valid, "4.7.4", "4.x.4", 1), "Standards-Version"},
		{"invalid name", strings.Replace(valid, "Source: example", "Source: Upper", 1), "Source"},
		{"missing package", strings.Replace(valid, "Package: example\n", "", 1), "Package"},
		{"missing architecture", strings.Replace(valid, "Architecture: any\n", "", 1), "Architecture"},
		{"mixed any", strings.Replace(valid, "Architecture: any", "Architecture: any amd64", 1), "Architecture"},
		{"source architecture", strings.Replace(valid, "Architecture: any", "Architecture: source", 1), "Architecture"},
		{"empty description", strings.Replace(valid, "Description: example package", "Description: \t", 1), "Description"},
		{"no synopsis", strings.Replace(valid, "Description: example package", "Description:\n extended only", 1), "synopsis"},
		{"folded simple", strings.Replace(valid, "Source: example", "Source: example\n continued", 1), "single-line"},
		{"essential", valid + "Essential: true\n", "Essential"},
		{"multiarch", valid + "Multi-Arch: invalid\n", "Multi-Arch"},
		{"duplicate binary", valid + "\n" + binaryFixture, "duplicate Package"},
		{"duplicate field", valid + "pACKAGE: example\n", "duplicate field"},
		{"malformed", valid + "bad line\n", "stanza 2"},
		{"multiple vcs", sourceFixture + "Vcs-Git: git-url\nVcs-Hg: hg-url\n\n" + binaryFixture, "multiple Vcs"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := Decode(strings.NewReader(tt.input))
			if file != nil || err == nil || !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("got %#v, %v; want %q", file, err, tt.contains)
			}
		})
	}
}

func TestUdebStandardsException(t *testing.T) {
	input := strings.Replace(sourceFixture, "Standards-Version: 4.7.4\n", "", 1) + "\n" + binaryFixture + "Package-Type: udeb\n"
	if _, err := DecodeBytes([]byte(input)); err != nil {
		t.Fatal(err)
	}
	input += "\nPackage: example-other\nArchitecture: all\nDescription: normal package\n"
	if _, err := DecodeBytes([]byte(input)); err == nil {
		t.Fatal("mixed deb/udeb requires Standards-Version")
	}
}

func TestDecodeCRLFAndEmptyOptional(t *testing.T) {
	input := sourceFixture + "Uploaders: \t\nX-Empty:\n\n" + binaryFixture + "Essential: \t\n"
	file, err := DecodeBytes([]byte(strings.ReplaceAll(input, "\n", "\r\n")))
	if err != nil {
		t.Fatal(err)
	}
	if file.Source.Uploaders != "" || len(file.Source.ExtraFields) != 0 || file.Binaries[0].Essential != "" {
		t.Fatalf("empty fields: %#v", file)
	}
}

type errorReader struct{ err error }

func (r errorReader) Read([]byte) (int, error) { return 0, r.err }
func TestReadError(t *testing.T) {
	want := errors.New("read failure")
	for _, prefix := range []string{"", sourceFixture + "\n"} {
		_, err := Decode(io.MultiReader(strings.NewReader(prefix), errorReader{want}))
		if !errors.Is(err, want) {
			t.Fatalf("got %v; want wrapped read error", err)
		}
	}
}

func TestNormalizationBoundaries(t *testing.T) {
	input := sourceFixture + "Homepage: \t https://example.org \t\nBuild-Depends: \t first,\n\t  second  |\n third \t\n\n" + "Package: example\nArchitecture: any\nDescription: \t . \t\n .\n  .\n text  \n"
	file, err := DecodeBytes([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if got := file.Source.Homepage; got != "https://example.org" {
		t.Errorf("Homepage = %q", got)
	}
	if got := file.Source.BuildDepends; got != "first, second | third" {
		t.Errorf("BuildDepends = %q", got)
	}
	if got := file.Binaries[0].Description; got != ". \t\n.\n .\ntext  " {
		t.Errorf("Description = %q", got)
	}
}

func TestOtherStanzaFieldsPreserved(t *testing.T) {
	input := sourceFixture + "Package: extra-package\n\n" + binaryFixture + "Source: extra-source\n"
	file, err := DecodeBytes([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if got := file.Source.ExtraFields.Get("Package"); got != "extra-package" {
		t.Errorf("source extra Package = %q", got)
	}
	if got := file.Binaries[0].ExtraFields.Get("Source"); got != "extra-source" {
		t.Errorf("binary extra Source = %q", got)
	}
}
