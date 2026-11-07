// Package changes defines upload control stanzas (.changes), described in
// Debian Policy section 5.5:
// https://www.debian.org/doc/debian-policy/ch-controlfields.html
//
// Decode accepts unsigned payloads, not OpenPGP armor. Description is a
// multiline table of package summaries, distinct from DEBIAN/control's prose
// Description. Changes, file tables, dot escapes and unknown fields are
// preserved. Validation is partial and does not verify checksums or signatures.
package changes
