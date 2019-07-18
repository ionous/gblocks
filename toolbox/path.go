package toolbox

import (
	"strconv"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
)

func (g *blockGen) newMutationGenerator(itemName string) *blockGen {
	scope := block.Scope("a", itemName)
	return &blockGen{g.block, g.domGenerator, g.shadowing, scope, g.atomNum}
}

// path for atom.
func (g *blockGen) newAtomGenerator() *blockGen {
	// place the inputs of atoms inside the same block
	// increase the shadowing depth for sub-blocks
	// MOD-stravis: this was g.atomNum+1 -- but if we do that then the paths arent zero indexed
	return &blockGen{g.block, g.domGenerator, g.shadowing.Children(), g.scope, g.atomNum}
}

func (g *blockGen) newStatement(itemName string, block *dom.Block) *dom.Statement {
	pathName := g.newInput(itemName)
	return &dom.Statement{Name: pathName, Input: dom.BlockInput{block}}
}

func (g *blockGen) newValue(itemName string, block *dom.Block) *dom.Value {
	pathName := g.newInput(itemName)
	return &dom.Value{Name: pathName, Input: dom.BlockInput{block}}
}

func (g *blockGen) newField(itemName string, val string) *dom.Field {
	pathName := g.newInput(itemName)
	return &dom.Field{pathName, val}
}

func (g *blockGen) mutating() bool {
	return len(g.scope) > 0
}

func (g *blockGen) newInput(item string) (ret string) {
	if !g.mutating() {
		ret = item
	} else {
		zeroIndexed := strconv.Itoa(g.atomNum)
		g.atomNum++
		ret = block.Scope(g.scope, zeroIndexed, item)
	}
	return
}
