package gblocks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type MutationElStart struct {
	NextStatement interface{}
}

type MutationElControl struct {
	PreviousStatement, NextStatement interface{}
}
type MutationAltControl struct {
	PreviousStatement, NextStatement interface{}
}

func TestShapeDecompose(t *testing.T) {
	testShape(t, func(ws *Workspace) {
		b := ws.NewBlock((*ShapeTest)(nil))
		d := ws.GetDataById(b.Id).(*ShapeTest)
		d.Mutant.els = append(d.Mutant.els, &MutationEl{}, &MutationAlt{}, &MutationAlt{})
		//
		data := ws.BlockData(b)
		m, ok := data.Mutation(b.GetInput("MUTANT"))
		require.True(t, ok)

		mutationUiBlocks := reduceBlocks(decompose(ws, m))
		require.Equal(t, []string{
			"mutation_el_start", "mutation_el_control", "mutation_alt_control", "mutation_alt_control",
		}, str)
	})
}

func reduceBlocks(block *Block) (ret []string) {
	for i := 0; block != nil && i < 100; i++ {
		ret = append(ret, block.Type)
		block = block.GetNextBlock()
	}
	return
}
