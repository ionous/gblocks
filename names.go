package gblocks

import (
	"github.com/fatih/camelcase"
	r "reflect"
	"strings"
)

// PascalCase -> under_score
func toTypeName(t r.Type) string {
	split := camelcase.Split(t.Name())
	for i, s := range split {
		split[i] = strings.ToLower(s)
	}
	return strings.Join(split, "_")
}

// CAPITAL_NAME -> PascalCase
func toFieldName(s string) string {
	split := strings.Split(s, "_")
	for i, s := range split {
		split[i] = string(s[0]) + strings.ToLower(s[1:])
	}
	return strings.Join(split, "")
}
