package debiancontrol

import (
	"fmt"
	"strings"
)

func validatePackageName(name string) error {
	if len(name) < 2 {
		return fmt.Errorf("package name %q must contain at least two characters", name)
	}
	for i, r := range name {
		alphanumeric := r >= 'a' && r <= 'z' || r >= '0' && r <= '9'
		if !alphanumeric && (i == 0 || r != '+' && r != '-' && r != '.') {
			return fmt.Errorf("invalid package name %q", name)
		}
	}
	return nil
}

func validateArchitecture(value string) error {
	architectures := strings.Fields(value)
	if len(architectures) == 0 {
		return fmt.Errorf("value is required")
	}
	for _, arch := range architectures {
		if arch == "source" || ((arch == "all" || arch == "any") && len(architectures) != 1) {
			return fmt.Errorf("invalid architecture list %q", value)
		}
	}
	return nil
}

func validateStandardsVersion(version string) error {
	parts := strings.Split(version, ".")
	if len(parts) != 3 && len(parts) != 4 {
		return fmt.Errorf("version %q must have three or four numeric components", version)
	}
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("version %q contains an empty component", version)
		}
		for _, r := range part {
			if r < '0' || r > '9' {
				return fmt.Errorf("version %q contains a nonnumeric component", version)
			}
		}
	}
	return nil
}

func validateEssential(value string) error {
	if value != "" && value != "yes" && value != "no" {
		return fmt.Errorf("expected yes or no, got %q", value)
	}
	return nil
}

func validateMultiArch(value string) error {
	switch value {
	case "", "no", "same", "foreign", "allowed":
		return nil
	default:
		return fmt.Errorf("invalid value %q", value)
	}
}
