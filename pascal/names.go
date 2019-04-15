package pascal

import (
	"strings"

	"github.com/fatih/camelcase"
)

// camelCase is used for custom strings.
func ToCamel(s string) (ret string) {
	if len(s) > 0 {
		split := camelcase.Split(s)
		split[0] = strings.ToLower(split[0])
		ret = strings.Join(split, "")
	}
	return
}

// CAPS_CASE is used for input names.
func ToCaps(s string) (ret string) {
	if len(s) > 0 {
		split := camelcase.Split(s)
		for i, s := range split {
			split[i] = strings.ToUpper(s)
		}
		ret = strings.Join(split, "_")
	}
	return
}

// space separated names are used for display.
func ToSpaces(s string) (ret string) {
	if len(s) > 0 {
		split := camelcase.Split(s)
		for i, s := range split {
			split[i] = strings.ToLower(s)
		}
		ret = strings.Join(split, " ")
	}
	return
}

// underscored_names are used for type names.
func ToUnderscore(s string) (ret string) {
	if len(s) > 0 {
		split := camelcase.Split(s)
		for i, s := range split {
			split[i] = strings.ToLower(s)
		}
		ret = strings.Join(split, "_")
	}
	return
}
