package deb822

import (
	"fmt"
	"strings"
)

func normalizeSimple(value string) (string, error) {
	if strings.ContainsRune(value, '\n') {
		return "", fmt.Errorf("must be a single-line field")
	}
	return strings.Trim(value, " \t"), nil
}

func normalizeFolded(value string) string {
	return strings.Join(strings.FieldsFunc(value, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	}), " ")
}
