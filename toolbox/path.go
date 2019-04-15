package toolbox

import (
	"strconv"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
)

func (g *blockGen) newMutationGenerator(itemName string) *blockGen {
	scope := block.Scope("a", g.block.Id, itemName)
	return &blockGen{g.block, g.domGenerator, g.shadowing, scope, g.atomNum}
}

// path for atom.
func (g *blockGen) newAtomGenerator() *blockGen {
	// place the inputs of atoms inside the same block
	// increase the shadowing depth for sub-blocks
	return &blockGen{g.block, g.domGenerator, g.shadowing.Children(), g.scope, g.atomNum + 1}
}

func (g *blockGen) newStatement(itemName string, block *dom.Block) *dom.Statement {
	pathName := g.pathName(itemName)
	return &dom.Statement{Name: pathName, Input: dom.BlockInput{block}}
}

func (g *blockGen) newValue(itemName string, block *dom.Block) *dom.Value {
	pathName := g.pathName(itemName)
	return &dom.Value{Name: pathName, Input: dom.BlockInput{block}}
}

func (g *blockGen) newField(itemName string, val string) *dom.Field {
	pathName := g.pathName(itemName)
	return &dom.Field{pathName, val}
}

func (g *blockGen) mutating() bool {
	return len(g.scope) > 0
}

func (g *blockGen) pathName(item string) (ret string) {
	if !g.mutating() {
		ret = item
	} else {
		ret = block.Scope(g.scope, strconv.Itoa(g.atomNum), item)
	}
	return
}