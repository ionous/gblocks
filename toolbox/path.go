package toolbox

import (
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/dom"
)

func (g *blockGen) newMutationGenerator(itemName string, atoms []*dom.Atom) *blockGen {
	mgen := &mutatorGen{atoms: atoms, atomIndex: -1}
	return &blockGen{g.block, g.domGenerator, g.shadowing, mgen}
}

func (g *blockGen) isMutating() bool {
	return g.mutating != nil
}

// next atom in a list of atoms
func (g *blockGen) newAtomGenerator() *blockGen {
	return &blockGen{g.block, g.domGenerator, g.shadowing.Children(), g.mutating.nextAtom()}
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

func (g *blockGen) newInput(item string) (ret string) {
	if !g.isMutating() {
		ret = item
	} else if a, ok := g.mutating.currentAtom(); !ok {
		ret = item // still handling the mutation struct
	} else {
		ret = block.Scope("a", a.Name, item)
	}
	return
}
