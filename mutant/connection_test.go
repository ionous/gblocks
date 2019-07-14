package mutant_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mock"
	"github.com/ionous/gblocks/mutant"
	"github.com/stretchr/testify/require"
)

// what we want out is muiblock ids -> input connections
// we need a workspace block with atoms and some sort of faked output connections
// we need a mui container with blocks of quarks to match
// we need to test the save, verify it, muck it up, restore -- verify the restore.
func TestSaveConnections(t *testing.T) {
	// create a workspace block with some expanded atoms
	var reg mock.Registry
	ws := reg.NewMockSpace()
	b := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.expandedInputs))
	// connect all of b's inputs to fake connections of dst
	dst := mock.CreateBlock("dest", nil)
	b.Workspace, dst.Workspace = ws, ws // workspace is used to test for valid connections
	connections := make(map[string]*mock.MockConnection)
	for i := 0; i < b.NumInputs(); i++ {
		in := b.Input(i)
		if inc := in.Connection(); inc != nil {
			zeroIndexed := strconv.Itoa(i)
			tgt := &mock.MockConnection{
				Name:   "c" + zeroIndexed,
				Source: dst,
			}
			inc.Connect(tgt)
			connections[tgt.Name] = tgt
		}
	}

	// create a matching mui block
	muispace := reg.NewMockSpace()
	blocks := mutant.NewMutatedBlocks()
	blockMutations := mock.NewMutations(common.inputs, common.quarks)
	require.NoError(t, blockMutations.Preregister("mockType", &reg))
	mb := blocks.CreateMutatedBlock(b, common.inputAtoms)
	c, e := blockMutations.CreateMui(muispace, mb)
	require.NoError(t, e)
	//
	initialConnections := []string{
		block.Scope("a", "mock", "I1", "0", "TERM:c1"),
		block.Scope("a", "mock", "I1", "1", "TERM:c2"),
		block.Scope("a", "mock", "I3", "0", "TERM:c6"),
		block.Scope("a", "mock", "I3", "2", "STATE:c9")}
	require.Equal(t, initialConnections, inputConnections(b),
		"initial connections from block")
	cs := mutant.SaveConnections(b, c)
	require.Equal(t, initialConnections, savedConnections(cs),
		"initial connections from save")
	// rotate the first connections, leave the last alone
	// note: we cant delete a connection without actually deleting the input
	// we test that at the end by restoring into an empty block
	remap := mutant.SavedConnections{
		block.Scope("a", "mock", "I1", "0", "TERM"):  connections["c6"],
		block.Scope("a", "mock", "I1", "1", "TERM"):  connections["c2"],
		block.Scope("a", "mock", "I3", "0", "TERM"):  connections["c1"],
		block.Scope("a", "mock", "I3", "2", "STATE"): connections["c9"],
	}
	remap.RestoreConnections(b)
	// test the connections by walking the block
	remapped := []string{
		block.Scope("a", "mock", "I1", "0", "TERM:c6"),
		block.Scope("a", "mock", "I1", "1", "TERM:c2"),
		block.Scope("a", "mock", "I3", "0", "TERM:c1"),
		block.Scope("a", "mock", "I3", "2", "STATE:c9"),
	}
	require.Equal(t, remapped, inputConnections(b),
		"remapped connections")
	// finally, restore into an empty block
	emptyBlock := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
	remap.RestoreConnections(emptyBlock)
	require.Empty(t, inputConnections(emptyBlock),
		"restored connections into an empty block")
}

func inputConnections(b block.Shape) (ret []string) {
	for i := 0; i < b.NumInputs(); i++ {
		in := b.Input(i)
		if inc := in.Connection(); inc != nil {
			s := inc.(*mock.MockConnection).String()
			ret = append(ret, s)
		}
	}
	return
}

// returns the target connections, and the sources they point back to
func savedConnections(cs mutant.SavedConnections) (ret []string) {
	// sort the sources of the saved connections ( the mui block input names )
	var srcs []string
	for k, _ := range cs {
		srcs = append(srcs, k)
	}
	sort.Strings(srcs)
	for _, src := range srcs {
		c := cs[src]
		tgt := c.(*mock.MockConnection).Name
		ret = append(ret, src+":"+tgt)
	}
	return
}
