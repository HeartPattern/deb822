package binarycontrol

import (
	"fmt"
	"strings"

	"github.com/HeartPattern/deb822/internal/controlfile"
)

// Validate checks required fields and basic architecture, synopsis and enum
// constraints. It does not validate full version, dependency or address syntax.
func (f *File) Validate() error {
	if f == nil {
		return fmt.Errorf("nil file")
	}
	if err := controlfile.Required(
		controlfile.Field{Name: "Package", Value: f.Name}, controlfile.Field{Name: "Version", Value: f.Version},
		controlfile.Field{Name: "Architecture", Value: f.Architecture}, controlfile.Field{Name: "Maintainer", Value: f.Maintainer},
		controlfile.Field{Name: "Description", Value: f.Description},
	); err != nil {
		return err
	}
	architectures := strings.Fields(f.Architecture)
	if len(architectures) != 1 || architectures[0] == "source" || architectures[0] == "any" || strings.HasPrefix(architectures[0], "any-") || strings.HasSuffix(architectures[0], "-any") {
		return fmt.Errorf("invalid Architecture %q", f.Architecture)
	}
	synopsis, _, _ := strings.Cut(f.Description, "\n")
	if strings.TrimSpace(synopsis) == "" {
		return fmt.Errorf("missing Description synopsis")
	}
	if f.Essential != "" && f.Essential != "yes" && f.Essential != "no" {
		return fmt.Errorf("invalid Essential %q", f.Essential)
	}
	switch f.MultiArch {
	case "", "no", "same", "foreign", "allowed":
	default:
		return fmt.Errorf("invalid Multi-Arch %q", f.MultiArch)
	}
	if f.InstalledSize != "" {
		for _, r := range f.InstalledSize {
			if r < '0' || r > '9' {
				return fmt.Errorf("invalid Installed-Size %q", f.InstalledSize)
			}
		}
	}
	return nil
}
