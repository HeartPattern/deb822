package debiancontrol

import "github.com/HeartPattern/deb822"

// File is a debian/control template, not a generated DEBIAN/control file.
// Binary fields are stored as declared; source defaults are not copied into them.
type File struct {
	Source   Source
	Binaries []Binary
}

// Source describes the first stanza. Relationship expressions, addresses,
// repository locations and substitution variables are retained as strings.
type Source struct {
	Name                string `control:"Source"`
	Maintainer          string `control:"Maintainer"`
	Uploaders           string `control:"Uploaders,folded"`
	Section             string `control:"Section"`
	Priority            string `control:"Priority"`
	BuildDepends        string `control:"Build-Depends,folded"`
	BuildDependsArch    string `control:"Build-Depends-Arch,folded"`
	BuildDependsIndep   string `control:"Build-Depends-Indep,folded"`
	BuildConflicts      string `control:"Build-Conflicts,folded"`
	BuildConflictsArch  string `control:"Build-Conflicts-Arch,folded"`
	BuildConflictsIndep string `control:"Build-Conflicts-Indep,folded"`
	StandardsVersion    string `control:"Standards-Version"`
	Homepage            string `control:"Homepage"`
	VcsBrowser          string `control:"Vcs-Browser"`
	VcsArch             string `control:"Vcs-Arch"`
	VcsBzr              string `control:"Vcs-Bzr"`
	VcsCvs              string `control:"Vcs-Cvs"`
	VcsDarcs            string `control:"Vcs-Darcs"`
	VcsGit              string `control:"Vcs-Git"`
	VcsHg               string `control:"Vcs-Hg"`
	VcsMtn              string `control:"Vcs-Mtn"`
	VcsSvn              string `control:"Vcs-Svn"`
	Testsuite           string `control:"Testsuite"`
	RulesRequiresRoot   string `control:"Rules-Requires-Root"`
	// ExtraFields retains unknown, nonempty fields with lowercase names.
	ExtraFields deb822.Stanza `control:",remain"`
}

// Binary describes a binary package to build, not its generated control file.
// Optional string fields are empty when absent. Relationship expressions and
// BuildProfiles are retained as strings, without evaluating substitutions.
type Binary struct {
	Name         string `control:"Package"`
	Architecture string `control:"Architecture"`
	Section      string `control:"Section"`
	Priority     string `control:"Priority"`
	Essential    string `control:"Essential"`
	Depends      string `control:"Depends,folded"`
	PreDepends   string `control:"Pre-Depends,folded"`
	Recommends   string `control:"Recommends,folded"`
	Suggests     string `control:"Suggests,folded"`
	Breaks       string `control:"Breaks,folded"`
	Conflicts    string `control:"Conflicts,folded"`
	Provides     string `control:"Provides,folded"`
	Replaces     string `control:"Replaces,folded"`
	Enhances     string `control:"Enhances,folded"`
	// Description contains the synopsis and extended description separated by
	// newlines. Dot escapes and remaining indentation are preserved as decoded.
	Description   string        `control:"Description,multiline"`
	Homepage      string        `control:"Homepage"`
	BuiltUsing    string        `control:"Built-Using,folded"`
	PackageType   string        `control:"Package-Type"`
	BuildProfiles string        `control:"Build-Profiles,folded"`
	MultiArch     string        `control:"Multi-Arch"`
	ExtraFields   deb822.Stanza `control:",remain"`
}
