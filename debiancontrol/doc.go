// Package debiancontrol decodes source package templates (debian/control),
// as described in https://www.debian.org/doc/debian-policy/ch-controlfields.html.
// A File contains a Source stanza followed by Binary stanzas. These differ from
// generated binary package control files (DEBIAN/control).
//
// Fields are mapped using control struct tags. Simple fields are trimmed,
// folded fields are unfolded, and multiline values retain their newlines,
// remaining indentation and dot escapes. Description rendering is left to callers.
// Unknown fields are captured with control:",remain"; empty fields are ignored. Field names in
// ExtraFields are lowercase; original comments and formatting are not retained.
// Source defaults are not applied to binary stanzas, and optional fields retain
// their zero values. Dependency expressions and substitution variables are not
// evaluated. Validation is structural and partial, not a full Policy checker.
package debiancontrol
