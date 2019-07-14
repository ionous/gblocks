package block

import "strings"

// helper to separate pieces of ids
// see Blockly.utils.genUid.soup_ for characters blockly uses
func Scope(str ...string) string {
	return strings.Join(str, " ")
}
