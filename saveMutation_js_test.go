package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShapeSave(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		b := ws.NewBlock((*ShapeTest)(nil))
		//
		d := ws.GetDataById(b.Id).(*ShapeTest)
		d.Mutant.els = append(d.Mutant.els, &MutationEl{}, &MutationAlt{}, &MutationAlt{})
		//
		el := b.mutationToDom(ws)
		text := `<mutation><input name="MUTANT" types="mutation_el,mutation_alt" elements="0,1,1"></mutation>`
		require.Equal(t, text, el.OuterHTML())
	})
}

func TestShapeRestore(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		b := ws.NewBlock((*ShapeTest)(nil))
		//
		d := ws.GetDataById(b.Id).(*ShapeTest)
		d.Mutant.els = append(d.Mutant.els, &MutationEl{}, &MutationAlt{}, &MutationAlt{})
		//
		el := b.mutationToDom(ws)
		text := `<mutation><input name="MUTANT" types="mutation_el,mutation_alt" elements="0,1,1"></mutation>`
		require.Equal(t, text, el.OuterHTML())
	})
}
