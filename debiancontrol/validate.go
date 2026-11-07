package debiancontrol

import (
	"fmt"
	"strings"
)

// Validate checks stanza constraints and file-wide rules. It does not validate
// relationship expressions, maintainer addresses, repository URLs, build
// profiles or substitution syntax.
func (f *File) Validate() error {
	if f == nil {
		return fmt.Errorf("debian/control: nil file")
	}
	if err := f.Source.Validate(); err != nil {
		return fmt.Errorf("source stanza: %w", err)
	}
	if len(f.Binaries) == 0 {
		return fmt.Errorf("debian/control: missing binary stanza")
	}
	allUdeb := true
	names := make(map[string]bool)
	for i, binary := range f.Binaries {
		if err := binary.Validate(); err != nil {
			return fmt.Errorf("stanza %d: %w", i+2, err)
		}
		if names[binary.Name] {
			return fmt.Errorf("stanza %d: duplicate Package %q", i+2, binary.Name)
		}
		names[binary.Name] = true
		allUdeb = allUdeb && binary.PackageType == "udeb"
	}
	if f.Source.StandardsVersion == "" && !allUdeb {
		return fmt.Errorf("source stanza: missing Standards-Version")
	}
	return nil
}

// Validate checks the source stanza's basic field constraints. An absent
// Standards-Version is allowed here; File.Validate checks the udeb exception.
func (s Source) Validate() error {
	if err := validatePackageName(s.Name); err != nil {
		return fmt.Errorf("Source: %w", err)
	}
	if strings.TrimSpace(s.Maintainer) == "" {
		return fmt.Errorf("missing Maintainer")
	}
	if s.StandardsVersion != "" {
		if err := validateStandardsVersion(s.StandardsVersion); err != nil {
			return fmt.Errorf("Standards-Version: %w", err)
		}
	}
	vcsCount := 0
	for _, value := range []string{s.VcsArch, s.VcsBzr, s.VcsCvs, s.VcsDarcs, s.VcsGit, s.VcsHg, s.VcsMtn, s.VcsSvn} {
		if value != "" {
			vcsCount++
		}
	}
	if vcsCount > 1 {
		return fmt.Errorf("multiple Vcs repository types")
	}
	return nil
}

// Validate checks the binary stanza's basic field constraints. Package-name
// uniqueness across stanzas is checked by File.Validate.
func (b Binary) Validate() error {
	if err := validatePackageName(b.Name); err != nil {
		return fmt.Errorf("Package: %w", err)
	}
	if err := validateArchitecture(b.Architecture); err != nil {
		return fmt.Errorf("Architecture: %w", err)
	}
	synopsis, _, _ := strings.Cut(b.Description, "\n")
	if strings.TrimSpace(synopsis) == "" {
		return fmt.Errorf("missing Description synopsis")
	}
	if err := validateEssential(b.Essential); err != nil {
		return fmt.Errorf("Essential: %w", err)
	}
	if err := validateMultiArch(b.MultiArch); err != nil {
		return fmt.Errorf("Multi-Arch: %w", err)
	}
	return nil
}
