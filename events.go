package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

type Events struct {
	*js.Object
	fireQueue *js.Object `js:"FIRE_QUEUE_"` // array
	Move      string     `js:"COMMENT_MOVE"`
}

func GetEvents() (ret *Events) {
	if blockly := js.Global.Get("Blockly"); blockly.Bool() {
		obj := blockly.Get("Events")
		ret = &Events{Object: obj}
	}
	return
}

func (ns *Events) IsEnabled() bool {
	return ns.Call("isEnabled").Bool()
}

func (ns *Events) TestFire(evt *js.Object) {
	if ns.IsEnabled() {
		ns.push(evt)
		ns.fireNow()
	}
}

func (ns *Events) push(evt *js.Object) {
	ns.fireQueue.Call("push", evt)
}

func (ns *Events) fireNow() {
	ns.Call("fireNow_")
}
