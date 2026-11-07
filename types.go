package deb822

import "strings"

type Token interface {
	isToken()
}

type Field struct {
	Name  string
	Value string
}

func (*Field) isToken() {}

type Continuation struct {
	Value string
}

func (*Continuation) isToken() {}

type Separator struct {
}

func (*Separator) isToken() {}

type Stanza map[string]string

func (s Stanza) Get(key string) string {
	if value, ok := s[strings.ToLower(key)]; ok {
		return value
	}
	for name, value := range s {
		if strings.EqualFold(name, key) {
			return value
		}
	}
	return ""
}
