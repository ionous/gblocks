package block

import (
	"github.com/fatih/camelcase"
	"strings"
)

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

// PascalCased -> lower spaces
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

// PascalCased -> CAPITAL_NAME
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

// under_score (or CAPS_NAME) -> PascalCased
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
