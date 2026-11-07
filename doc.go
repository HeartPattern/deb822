// Package deb822 decodes Debian control data into tokens and stanzas.
// UnmarshalStanza maps stanzas to structs using control tags.
// It handles common syntax and field formats; file-specific schemas, validation
// and interpretation belong in higher-level packages.
package deb822
