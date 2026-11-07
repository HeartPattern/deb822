package debiancontrol

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/HeartPattern/deb822"
)

// Decode reads and validates one complete debian/control file. Comments are
// allowed. Empty fields are ignored; unknown fields are preserved in ExtraFields.
// Validation covers structural and basic field rules, not full Debian Policy
// compliance or the grammar of dependency and substitution expressions.
func Decode(reader io.Reader) (*File, error) {
	decoder := deb822.NewDecoder(reader)
	decoder.AllowComments = true
	stanza, err := decoder.Decode()
	if err == io.EOF {
		return nil, fmt.Errorf("debian/control: missing source stanza")
	}
	if err != nil {
		return nil, fmt.Errorf("source stanza: %w", err)
	}
	source, err := decodeSource(stanza)
	if err != nil {
		return nil, fmt.Errorf("source stanza: %w", err)
	}
	file := &File{Source: source}
	for index := 2; ; index++ {
		stanza, err := decoder.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("stanza %d: %w", index, err)
		}
		binary, err := decodeBinary(stanza)
		if err != nil {
			return nil, fmt.Errorf("stanza %d: %w", index, err)
		}
		file.Binaries = append(file.Binaries, binary)
	}
	if err := file.Validate(); err != nil {
		return nil, err
	}
	return file, nil
}

// DecodeBytes is Decode for an in-memory template.
func DecodeBytes(data []byte) (*File, error) { return Decode(bytes.NewReader(data)) }

func decodeSource(stanza deb822.Stanza) (Source, error) {
	var source Source
	if err := deb822.UnmarshalStanza(nonemptyFields(stanza), &source); err != nil {
		return Source{}, err
	}
	return source, nil
}

func decodeBinary(stanza deb822.Stanza) (Binary, error) {
	var binary Binary
	if err := deb822.UnmarshalStanza(nonemptyFields(stanza), &binary); err != nil {
		return Binary{}, err
	}
	return binary, nil
}

// Empty fields are ignored specifically in debian/control templates.
func nonemptyFields(stanza deb822.Stanza) deb822.Stanza {
	fields := make(deb822.Stanza)
	for name, value := range stanza {
		if strings.Trim(value, " \t\n") != "" {
			fields[name] = value
		}
	}
	return fields
}
