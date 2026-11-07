package changes

import (
	"fmt"
	"strings"

	"github.com/HeartPattern/deb822/internal/controlfile"
)

// Validate checks required fields, conditional Binary/Description presence,
// architecture constraints and multiline layout. It does not verify dates,
// versions, file rows, checksums, signatures or tag2upload provenance.
func (f *File) Validate() error {
	if f == nil {
		return fmt.Errorf("nil file")
	}
	if err := controlfile.Required(
		controlfile.Field{Name: "Format", Value: f.Format}, controlfile.Field{Name: "Date", Value: f.Date},
		controlfile.Field{Name: "Source", Value: f.Source}, controlfile.Field{Name: "Architecture", Value: f.Architecture},
		controlfile.Field{Name: "Version", Value: f.Version}, controlfile.Field{Name: "Distribution", Value: f.Distribution},
		controlfile.Field{Name: "Maintainer", Value: f.Maintainer}, controlfile.Field{Name: "Changes", Value: f.Changes},
		controlfile.Field{Name: "Checksums-Sha1", Value: f.ChecksumsSha1}, controlfile.Field{Name: "Checksums-Sha256", Value: f.ChecksumsSha256},
		controlfile.Field{Name: "Files", Value: f.Files},
	); err != nil {
		return err
	}
	hasBinaries := false
	for _, arch := range strings.Fields(f.Architecture) {
		if arch == "any" || strings.HasPrefix(arch, "any-") || strings.HasSuffix(arch, "-any") {
			return fmt.Errorf("invalid Architecture %q", f.Architecture)
		}
		if arch != "source" {
			hasBinaries = true
		}
	}
	if hasBinaries {
		if err := controlfile.Required(controlfile.Field{Name: "Binary", Value: f.Binary}, controlfile.Field{Name: "Description", Value: f.Description}); err != nil {
			return err
		}
	}
	return controlfile.Multiline(
		controlfile.Field{Name: "Changes", Value: f.Changes}, controlfile.Field{Name: "Description", Value: f.Description},
		controlfile.Field{Name: "Checksums-Sha1", Value: f.ChecksumsSha1}, controlfile.Field{Name: "Checksums-Sha256", Value: f.ChecksumsSha256},
		controlfile.Field{Name: "Files", Value: f.Files},
	)
}
