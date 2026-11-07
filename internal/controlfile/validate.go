package controlfile

import (
	"fmt"
	"strings"
)

// Field is a field name and value used by Required.
type Field struct{ Name, Value string }

func Required(fields ...Field) error {
	for _, field := range fields {
		if strings.TrimSpace(field.Value) == "" {
			return fmt.Errorf("missing %s", field.Name)
		}
	}
	return nil
}

// Multiline checks that table-style fields start on a continuation line.
// Row syntax and referenced file contents are not validated.
func Multiline(fields ...Field) error {
	for _, field := range fields {
		if field.Value == "" {
			continue
		}
		first, rest, found := strings.Cut(field.Value, "\n")
		if !found || strings.TrimSpace(first) != "" || strings.TrimSpace(rest) == "" {
			return fmt.Errorf("%s: expected an empty first line followed by content", field.Name)
		}
	}
	return nil
}
