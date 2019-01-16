package gblocks

import (
	"github.com/fatih/camelcase"
	r "reflect"
	"strings"
)

// PascalCase -> under_score
func toTypeName(t r.Type) (ret TypeName) {
	str := pascalToUnderscore(t.Name())
	return TypeName(str)
}

func pascalToUnderscore(s string) (ret string) {
	if len(s) > 0 {
		split := camelcase.Split(s)
		for i, s := range split {
			split[i] = strings.ToLower(s)
		}
		ret = strings.Join(split, "_")
	}
	return
}

func pascalToSpace(s string) (ret string) {
	if len(s) > 0 {
		split := camelcase.Split(s)
		for i, s := range split {
			split[i] = strings.ToLower(s)
		}
		ret = strings.Join(split, " ")
	}
	return
}

// PascalCase -> CAPITAL_NAME
func pascalToCaps(s string) (ret string) {
	if len(s) > 0 {
		split := camelcase.Split(s)
		for i, s := range split {
			split[i] = strings.ToUpper(s)
		}
		ret = strings.Join(split, "_")
	}
	return
}

// under_score (or CAPS_NAME) -> PascalCase
func underscoreToPascal(s string) (ret string) {
	if len(s) > 0 {
		split := strings.Split(s, "_")
		for i, s := range split {
			split[i] = strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
		}
		ret = strings.Join(split, "")
	}
	return
}
