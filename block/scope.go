package block

import "strings"

func Scope(str ...string) string {
	return strings.Join(str, "$")
}
