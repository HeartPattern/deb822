package dsc

import "heartpattern.io/deb822"

// File describes a source package's .dsc control stanza. File lists, checksums,
// package lists and relationship expressions remain strings. Checksums and
// OpenPGP signatures are not verified by this package.
type File struct {
	Format              string        `control:"Format"`
	Source              string        `control:"Source"`
	Binary              string        `control:"Binary,folded"`
	Architecture        string        `control:"Architecture"`
	Version             string        `control:"Version"`
	Maintainer          string        `control:"Maintainer"`
	Uploaders           string        `control:"Uploaders"`
	Homepage            string        `control:"Homepage"`
	VcsBrowser          string        `control:"Vcs-Browser"`
	VcsArch             string        `control:"Vcs-Arch"`
	VcsBzr              string        `control:"Vcs-Bzr"`
	VcsCvs              string        `control:"Vcs-Cvs"`
	VcsDarcs            string        `control:"Vcs-Darcs"`
	VcsGit              string        `control:"Vcs-Git"`
	VcsHg               string        `control:"Vcs-Hg"`
	VcsMtn              string        `control:"Vcs-Mtn"`
	VcsSvn              string        `control:"Vcs-Svn"`
	Testsuite           string        `control:"Testsuite"`
	Dgit                string        `control:"Dgit,folded"`
	GitTagTagger        string        `control:"Git-Tag-Tagger"`
	GitTagInfo          string        `control:"Git-Tag-Info"`
	StandardsVersion    string        `control:"Standards-Version"`
	BuildDepends        string        `control:"Build-Depends"`
	BuildDependsArch    string        `control:"Build-Depends-Arch"`
	BuildDependsIndep   string        `control:"Build-Depends-Indep"`
	BuildConflicts      string        `control:"Build-Conflicts"`
	BuildConflictsArch  string        `control:"Build-Conflicts-Arch"`
	BuildConflictsIndep string        `control:"Build-Conflicts-Indep"`
	PackageList         string        `control:"Package-List,multiline"`
	ChecksumsSha1       string        `control:"Checksums-Sha1,multiline"`
	ChecksumsSha256     string        `control:"Checksums-Sha256,multiline"`
	Files               string        `control:"Files,multiline"`
	ExtraFields         deb822.Stanza `control:",remain"`
}
