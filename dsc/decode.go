package dsc

import (
	"bytes"
	"fmt"
	"io"

	"heartpattern.io/deb822/internal/controlfile"
)

// Decode reads one unsigned .dsc stanza and checks basic constraints.
// For a clear-signed file, the caller must first extract its plaintext payload.
// This function does not verify signatures or checksums.
func Decode(reader io.Reader) (*File, error) {
	var file File
	if err := controlfile.Decode(reader, &file); err != nil {
		return nil, fmt.Errorf(".dsc: %w", err)
	}
	if err := file.Validate(); err != nil {
		return nil, fmt.Errorf(".dsc: %w", err)
	}
	return &file, nil
}

func DecodeBytes(data []byte) (*File, error) { return Decode(bytes.NewReader(data)) }
