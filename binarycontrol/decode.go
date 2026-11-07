package binarycontrol

import (
	"bytes"
	"fmt"
	"io"

	"heartpattern.io/deb822/internal/controlfile"
)

// Decode reads one generated DEBIAN/control stanza and checks basic constraints.
// Comments, empty fields and additional stanzas are rejected.
func Decode(reader io.Reader) (*File, error) {
	var file File
	if err := controlfile.Decode(reader, &file); err != nil {
		return nil, fmt.Errorf("DEBIAN/control: %w", err)
	}
	if err := file.Validate(); err != nil {
		return nil, fmt.Errorf("DEBIAN/control: %w", err)
	}
	return &file, nil
}

func DecodeBytes(data []byte) (*File, error) { return Decode(bytes.NewReader(data)) }
