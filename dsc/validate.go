package dsc

import (
	"fmt"
	"strings"

	"github.com/HeartPattern/deb822/internal/controlfile"
)

// Validate checks required fields and basic stanza constraints, not full row,
// version, relationship or address syntax, checksums or signature validity.
func (f *File) Validate() error {
	if f == nil {
		return fmt.Errorf("nil file")
	}
	if err := controlfile.Required(
		controlfile.Field{Name: "Format", Value: f.Format}, controlfile.Field{Name: "Source", Value: f.Source},
		controlfile.Field{Name: "Version", Value: f.Version}, controlfile.Field{Name: "Maintainer", Value: f.Maintainer},
		controlfile.Field{Name: "Checksums-Sha1", Value: f.ChecksumsSha1}, controlfile.Field{Name: "Checksums-Sha256", Value: f.ChecksumsSha256},
		controlfile.Field{Name: "Files", Value: f.Files},
	); err != nil {
		return err
	}
	if err := controlfile.Multiline(
		controlfile.Field{Name: "Checksums-Sha1", Value: f.ChecksumsSha1}, controlfile.Field{Name: "Checksums-Sha256", Value: f.ChecksumsSha256},
		controlfile.Field{Name: "Files", Value: f.Files}, controlfile.Field{Name: "Package-List", Value: f.PackageList},
	); err != nil {
		return err
	}
	if strings.TrimSpace(f.StandardsVersion) == "" && !onlyUdebs(f.PackageList) {
		return fmt.Errorf("missing Standards-Version")
	}
	architectures := strings.Fields(f.Architecture)
	for _, arch := range architectures {
		if arch == "source" {
			return fmt.Errorf("invalid Architecture %q", f.Architecture)
		}
		if arch == "any" {
			for _, other := range architectures {
				if other != "any" && other != "all" {
					return fmt.Errorf("invalid Architecture %q", f.Architecture)
				}
			}
		}
	}
	vcsCount := 0
	for _, value := range []string{f.VcsArch, f.VcsBzr, f.VcsCvs, f.VcsDarcs, f.VcsGit, f.VcsHg, f.VcsMtn, f.VcsSvn} {
		if value != "" {
			vcsCount++
		}
	}
	if vcsCount > 1 {
		return fmt.Errorf("multiple Vcs repository types")
	}
	return nil
}

// Package-List is the available evidence that a source produces only udebs.
func onlyUdebs(value string) bool {
	_, rows, found := strings.Cut(value, "\n")
	if !found || strings.TrimSpace(rows) == "" {
		return false
	}
	for _, row := range strings.Split(rows, "\n") {
		columns := strings.Fields(row)
		if len(columns) < 4 || columns[1] != "udeb" {
			return false
		}
	}
	return true
}
