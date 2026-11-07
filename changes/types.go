package changes

import "heartpattern.io/deb822"

// File describes an upload's .changes control stanza. Structured field values
// remain strings; checksums and OpenPGP signatures are not verified.
type File struct {
	Format          string        `control:"Format"`
	Date            string        `control:"Date"`
	Source          string        `control:"Source"`
	Binary          string        `control:"Binary,folded"`
	Architecture    string        `control:"Architecture"`
	Version         string        `control:"Version"`
	Distribution    string        `control:"Distribution"`
	Urgency         string        `control:"Urgency"`
	Maintainer      string        `control:"Maintainer"`
	ChangedBy       string        `control:"Changed-By"`
	GitTagTagger    string        `control:"Git-Tag-Tagger"`
	GitTagInfo      string        `control:"Git-Tag-Info"`
	Description     string        `control:"Description,multiline"`
	Closes          string        `control:"Closes"`
	Changes         string        `control:"Changes,multiline"`
	ChecksumsSha1   string        `control:"Checksums-Sha1,multiline"`
	ChecksumsSha256 string        `control:"Checksums-Sha256,multiline"`
	Files           string        `control:"Files,multiline"`
	ExtraFields     deb822.Stanza `control:",remain"`
}
