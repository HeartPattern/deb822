// Package dsc defines source package control stanzas (.dsc), described in
// Debian Policy section 5.4:
// https://www.debian.org/doc/debian-policy/ch-controlfields.html
//
// Decode accepts the unsigned payload of one stanza. OpenPGP armor must be
// handled by the caller; this package does not authenticate signatures.
// Multiline tables and unknown fields are preserved. Validation is partial:
// file rows, checksums and dependency expressions are not fully interpreted.
package dsc
