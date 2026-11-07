package binarycontrol

import "heartpattern.io/deb822"

// File describes a generated binary package's DEBIAN/control file.
// Values such as Source (which may include a version), dependency expressions,
// and InstalledSize are kept as strings rather than interpreted.
type File struct {
	Name          string        `control:"Package"`
	Source        string        `control:"Source"`
	Version       string        `control:"Version"`
	Section       string        `control:"Section"`
	Priority      string        `control:"Priority"`
	Architecture  string        `control:"Architecture"`
	Essential     string        `control:"Essential"`
	Depends       string        `control:"Depends"`
	PreDepends    string        `control:"Pre-Depends"`
	Recommends    string        `control:"Recommends"`
	Suggests      string        `control:"Suggests"`
	Breaks        string        `control:"Breaks"`
	Conflicts     string        `control:"Conflicts"`
	Provides      string        `control:"Provides"`
	Replaces      string        `control:"Replaces"`
	Enhances      string        `control:"Enhances"`
	InstalledSize string        `control:"Installed-Size"`
	Maintainer    string        `control:"Maintainer"`
	Description   string        `control:"Description,multiline"`
	Homepage      string        `control:"Homepage"`
	BuiltUsing    string        `control:"Built-Using"`
	MultiArch     string        `control:"Multi-Arch"`
	ExtraFields   deb822.Stanza `control:",remain"`
}
