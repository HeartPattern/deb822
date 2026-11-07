package deb822

import "testing"

func TestStanzaGet(t *testing.T) {
	for _, name := range []string{"package", "Package", "PACKAGE"} {
		s := Stanza{name: "hello"}
		for _, query := range []string{"package", "Package", "PACKAGE"} {
			if got := s.Get(query); got != "hello" {
				t.Errorf("key=%q query=%q: got %q", name, query, got)
			}
		}
		if got := s.Get("Missing"); got != "" {
			t.Errorf("missing key: got %q", got)
		}
	}
	var s Stanza
	if got := s.Get("Package"); got != "" {
		t.Errorf("nil stanza: got %q", got)
	}
}
