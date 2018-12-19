package gblocks

import (
	"github.com/fatih/camelcase"
	r "reflect"
	"strings"
)

// PascalCase -> under_score
func toTypeName(t r.Type) string {
	return pascalToUnderscore(t.Name())
}

func pascalToUnderscore(s string) string {
	split := camelcase.Split(s)
	for i, s := range split {
		split[i] = strings.ToLower(s)
	}
	return strings.Join(split, "_")
}

// PascalCase -> CAPITAL_NAME
func pascalToCaps(s string) string {
	split := camelcase.Split(s)
	for i, s := range split {
		split[i] = strings.ToUpper(s)
	}
	return strings.Join(split, "_")
}

// under_score (or CAPS_NAME) -> PascalCase
func underscoreToPascal(s string) string {
	split := strings.Split(s, "_")
	for i, s := range split {
		split[i] = strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
	}
	return strings.Join(split, "")
}
