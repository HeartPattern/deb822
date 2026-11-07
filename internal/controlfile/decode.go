// Package controlfile provides shared decoding for generated control files.
package controlfile

import (
	"fmt"
	"io"
	"strings"

	"heartpattern.io/deb822"
)

// Decode reads exactly one unsigned stanza. Unlike debian/control templates,
// generated control files do not allow comments or empty fields.
func Decode(reader io.Reader, destination any) error {
	decoder := deb822.NewDecoder(reader)
	stanza, err := decoder.Decode()
	if err == io.EOF {
		return fmt.Errorf("missing stanza")
	}
	if err != nil {
		return err
	}
	for name, value := range stanza {
		if strings.Trim(value, " \t\n") == "" {
			return fmt.Errorf("%s: empty field", name)
		}
	}
	if _, err := decoder.Decode(); err != io.EOF {
		if err != nil {
			return err
		}
		return fmt.Errorf("expected exactly one stanza")
	}
	return deb822.UnmarshalStanza(stanza, destination)
}
