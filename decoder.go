package deb822

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"strings"
)

type Decoder struct {
	tokenDecoder *TokenDecoder

	// AllowComments skips lines beginning with '#'. It is disabled by default.
	AllowComments bool
}

func NewDecoderBytes(b []byte) *Decoder {
	return NewDecoder(bytes.NewReader(b))
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		tokenDecoder: NewTokenDecoder(reader),
	}
}

// Decode reads the next nonempty stanza, storing field names in lowercase.
// Continuations are joined with newlines, preserving their remaining whitespace.
// Field-specific folding and dot escapes are left to the caller.
// The final stanza is returned before io.EOF, even without a trailing separator.
// After a syntax or read error, callers should stop using the decoder.
func (d *Decoder) Decode() (Stanza, error) {
	d.tokenDecoder.AllowComments = d.AllowComments
	stanza := make(Stanza)
	var name string
	var value strings.Builder
	finish := func() {
		if name != "" {
			stanza[name] = value.String()
		}
		value.Reset()
	}
	for {
		token, err := d.tokenDecoder.Decode()
		if err == io.EOF {
			finish()
			if len(stanza) == 0 {
				return nil, io.EOF
			}
			return stanza, nil
		}
		if err != nil {
			return nil, err
		}
		switch token := token.(type) {
		case *Separator:
			finish()
			if len(stanza) > 0 {
				return stanza, nil
			}
		case *Field:
			finish()
			name = strings.ToLower(token.Name)
			if _, exists := stanza[name]; exists {
				return nil, fmt.Errorf("duplicate field: %s", token.Name)
			}
			value.WriteString(token.Value)
		case *Continuation:
			if name == "" {
				return nil, fmt.Errorf("continuation without a preceding field")
			}
			value.WriteByte('\n')
			value.WriteString(token.Value)
		}
	}
}

// Decodes yields stanzas until EOF or the first error.
func (d *Decoder) Decodes() iter.Seq2[Stanza, error] {
	return func(yield func(Stanza, error) bool) {
		for {
			stanza, err := d.Decode()
			if err == io.EOF {
				return
			}
			if !yield(stanza, err) || err != nil {
				return
			}
		}
	}
}
