package debiancontrol

import "testing"

func validFile() File {
	return File{
		Source:   Source{Name: "example", Maintainer: "Maintainer <maint@example.org>", StandardsVersion: "4.7.4"},
		Binaries: []Binary{{Name: "example", Architecture: "any", Description: "example package"}},
	}
}

func TestFileValidationRules(t *testing.T) {
	tests := []struct {
		name   string
		change func(*File)
		valid  bool
	}{
		{"baseline", func(f *File) {}, true},
		{"package punctuation", func(f *File) { f.Source.Name = "libfoo1.0+dev" }, true},
		{"short source", func(f *File) { f.Source.Name = "a" }, false},
		{"uppercase source", func(f *File) { f.Source.Name = "Example" }, false},
		{"leading punctuation", func(f *File) { f.Source.Name = "+example" }, false},
		{"unicode package", func(f *File) { f.Binaries[0].Name = "exämple" }, false},
		{"blank maintainer", func(f *File) { f.Source.Maintainer = " \t" }, false},
		{"four part standards", func(f *File) { f.Source.StandardsVersion = "4.7.4.1" }, true},
		{"two part standards", func(f *File) { f.Source.StandardsVersion = "4.7" }, false},
		{"empty version component", func(f *File) { f.Source.StandardsVersion = "4..4" }, false},
		{"invalid version character", func(f *File) { f.Source.StandardsVersion = "4.a.4" }, false},
		{"no binaries", func(f *File) { f.Binaries = nil }, false},
		{"duplicate binaries", func(f *File) { f.Binaries = append(f.Binaries, f.Binaries[0]) }, false},
		{"missing architecture", func(f *File) { f.Binaries[0].Architecture = " \t" }, false},
		{"architecture list", func(f *File) { f.Binaries[0].Architecture = "linux-any arm64" }, true},
		{"all architecture", func(f *File) { f.Binaries[0].Architecture = "all" }, true},
		{"mixed all", func(f *File) { f.Binaries[0].Architecture = "all amd64" }, false},
		{"mixed any", func(f *File) { f.Binaries[0].Architecture = "any amd64" }, false},
		{"source architecture", func(f *File) { f.Binaries[0].Architecture = "source" }, false},
		{"missing synopsis", func(f *File) { f.Binaries[0].Description = "\nlong description" }, false},
		{"essential yes", func(f *File) { f.Binaries[0].Essential = "yes" }, true},
		{"essential no", func(f *File) { f.Binaries[0].Essential = "no" }, true},
		{"invalid essential", func(f *File) { f.Binaries[0].Essential = "true" }, false},
		{"multiarch no", func(f *File) { f.Binaries[0].MultiArch = "no" }, true},
		{"multiarch same", func(f *File) { f.Binaries[0].MultiArch = "same" }, true},
		{"multiarch foreign", func(f *File) { f.Binaries[0].MultiArch = "foreign" }, true},
		{"multiarch allowed", func(f *File) { f.Binaries[0].MultiArch = "allowed" }, true},
		{"invalid multiarch", func(f *File) { f.Binaries[0].MultiArch = "yes" }, false},
		{"vcs and browser", func(f *File) { f.Source.VcsGit = "repo"; f.Source.VcsBrowser = "browser" }, true},
		{"multiple vcs", func(f *File) { f.Source.VcsGit = "repo"; f.Source.VcsHg = "repo" }, false},
		{"missing standards", func(f *File) { f.Source.StandardsVersion = "" }, false},
		{"udeb exception", func(f *File) { f.Source.StandardsVersion = ""; f.Binaries[0].PackageType = "udeb" }, true},
		{"udeb invalid standards", func(f *File) { f.Source.StandardsVersion = "invalid"; f.Binaries[0].PackageType = "udeb" }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := validFile()
			tt.change(&file)
			if err := file.Validate(); (err == nil) != tt.valid {
				t.Fatalf("Validate() = %v; want valid=%v", err, tt.valid)
			}
		})
	}
	var file *File
	if err := file.Validate(); err == nil {
		t.Fatal("nil file must fail validation")
	}
}

func TestStanzaValidation(t *testing.T) {
	file := validFile()
	if err := file.Source.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := file.Binaries[0].Validate(); err != nil {
		t.Fatal(err)
	}
	// Whether Standards-Version is mandatory depends on all binary stanzas.
	file.Source.StandardsVersion = ""
	if err := file.Source.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := file.Validate(); err == nil {
		t.Fatal("normal binary requires Standards-Version")
	}
	file.Source.VcsGit, file.Source.VcsHg = "git", "hg"
	if err := file.Source.Validate(); err == nil {
		t.Fatal("source must reject multiple VCS types")
	}
	file.Binaries[0].Essential = "true"
	if err := file.Binaries[0].Validate(); err == nil {
		t.Fatal("binary must reject invalid Essential")
	}
}
