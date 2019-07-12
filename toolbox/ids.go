package toolbox

import "strconv"

type Ids interface {
	// note: while blockly's examples do not include "id" in the dom,
	// Blockly.Xml.domToBlockHeadless_() in xml.js does correctly parse it.
	// id is used in gblocks for "generic mutations" to distinguish b.t mutation blocks in a defined way
	// so, we autogenerate one if it is not assigned
	NewId() string
}

type IdGenerator struct {
	nextId int // unique id generator for toolbox.Builder
}

func (g *IdGenerator) NewId() string {
	oneIndexed := strconv.Itoa(g.nextId + 1)
	newId := "bl" + oneIndexed
	g.nextId++
	return newId
}
