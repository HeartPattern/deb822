package deb822

import (
	"fmt"
	"reflect"
	"strings"
)

// UnmarshalStanza maps a stanza into the exported string fields of a struct.
// The destination must be a non-nil pointer to a struct. Field names are matched
// case-insensitively; input names differing only by case are rejected.
//
// Tags use control:"Name[,simple|folded|multiline]"; simple is the default.
// Simple values are trimmed and must not contain newlines. Folded values have
// spaces, tabs and newlines collapsed. Multiline values are preserved verbatim.
// Untagged fields and control:"-" are ignored. One exported Stanza field tagged
// control:",remain" may receive unmapped fields, with lowercase names. Without
// a remain field, unmapped fields are discarded.
//
// Empty values are mapped; absent fields leave destination values unchanged.
// The remain map is replaced on each call and never aliases the input stanza.
// On error, the destination may be partially updated. No schema validation or
// field-specific interpretation is performed.
func UnmarshalStanza(stanza Stanza, destination any) error {
	target := reflect.ValueOf(destination)
	if target.Kind() != reflect.Pointer || target.IsNil() || target.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("control destination must be a non-nil pointer to a struct")
	}
	target = target.Elem()

	extra := make(Stanza)
	for name, value := range stanza {
		key := strings.ToLower(name)
		if _, exists := extra[key]; exists {
			return fmt.Errorf("duplicate field: %s", name)
		}
		extra[key] = value
	}

	var remain reflect.Value
	seen := make(map[string]bool)

	for i := 0; i < target.NumField(); i++ {
		field := target.Type().Field(i)
		tag, tagged := field.Tag.Lookup("control")
		if !tagged || tag == "-" {
			continue
		}

		name, format, _ := strings.Cut(tag, ",")
		if format == "remain" {
			if name != "" {
				return fmt.Errorf("%s: remain tag must not specify a field name", field.Name)
			}
			if remain.IsValid() {
				return fmt.Errorf("%s: multiple remain fields", field.Name)
			}
			remain = target.Field(i)
			if !remain.CanSet() || remain.Type() != reflect.TypeFor[Stanza]() {
				return fmt.Errorf("%s: remain field must be an exported deb822.Stanza", field.Name)
			}
			continue
		}

		if name == "" || name == "-" {
			return fmt.Errorf("%s: invalid control tag %q", field.Name, tag)
		}

		key := strings.ToLower(name)
		if seen[key] {
			return fmt.Errorf("%s: duplicate control tag %q", field.Name, name)
		}
		seen[key] = true

		output := target.Field(i)
		if !output.CanSet() || output.Kind() != reflect.String {
			return fmt.Errorf("%s: control field must be an exported string", field.Name)
		}

		value, exists := extra[key]
		switch format {
		case "", "simple":
			normalized, err := normalizeSimple(value)
			if err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
			value = normalized
		case "folded":
			value = normalizeFolded(value)
		case "multiline":
			// Preserve the decoded value; field-specific interpretation belongs elsewhere.
		default:
			return fmt.Errorf("%s: unsupported control format %q", field.Name, format)
		}

		if !exists {
			continue
		}

		output.SetString(value)
		delete(extra, key)
	}

	if remain.IsValid() {
		remain.Set(reflect.ValueOf(extra))
	}

	return nil
}
