package blockly

import "github.com/gopherjs/gopherjs/js"

// via Blockly.Utils
type Utils struct {
	*js.Object
}

// returns string
func (u *Utils) GenUid() (ret string) {
	if obj := u.Call("genUid", u.Object); obj != nil && obj.Bool() {
		ret = obj.String()
	}
	return
}
