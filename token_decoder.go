package deb822

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"strings"
)

type TokenDecoder struct {
	reader *bufio.Reader

	// AllowComments skips lines beginning with '#'. It is disabled by default.
	AllowComments bool
}

func NewTokenDecoder(reader io.Reader) *TokenDecoder {
	return &TokenDecoder{
		reader: bufio.NewReader(reader),
	}
}

// Decode returns one physical line token, without its LF or CRLF terminator.
// Continuation indentation beyond the first space or tab is preserved.
// When AllowComments is enabled, comment lines are skipped without emitting tokens.
func (t *TokenDecoder) Decode() (Token, error) {
	var line string
	for {
		var err error
		line, err = t.reader.ReadString('\n')
		if len(line) == 0 && err != nil {
			return nil, err
		}
		if err != nil && err != io.EOF {
			return nil, err
		}
		if t.AllowComments && strings.HasPrefix(line, "#") {
			continue
		}
		if before, ok := strings.CutSuffix(line, "\n"); ok {
			line = strings.TrimSuffix(before, "\r")
		}
		break
	}

	if len(line) == 0 || !strings.ContainsFunc(line, func(r rune) bool { return r != ' ' && r != '\t' }) {
		return &Separator{}, nil
	}

	if line[0] == ' ' || line[0] == '\t' {
		return &Continuation{
			Value: line[1:],
		}, nil
	}

	key, value, ok := strings.Cut(line, ":")
	if !ok {
		return nil, fmt.Errorf("field separator not found: %q", line)
	}

	if len(key) == 0 || key[0] == '-' || key[0] == '#' || strings.ContainsFunc(key, func(r rune) bool { return r < '!' || r > '~' }) {
		return nil, fmt.Errorf("invalid key: %s", key)
	}

	return &Field{
		Name:  key,
		Value: strings.TrimLeft(value, " \t"),
	}, nil
}

func (t *TokenDecoder) Decodes() iter.Seq2[Token, error] {
	return func(yield func(Token, error) bool) {
		for {
			token, err := t.Decode()
			if err == io.EOF {
				return
			}
			if !yield(token, err) || err != nil {
				break
			}
		}
	}
}
