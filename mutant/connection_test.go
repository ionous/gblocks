package mutant_test

// what we want out is muiblock ids -> input connections
// we need a workspace block with atoms and some sort of faked output connections
// we need a mui container with blocks of quarks to match
// we need to test the save, verify it, muck it up, restore -- verify the restore.
// func TestSaveConnections(t *testing.T) {
// 	// create a workspace block with some expanded atoms
// 	var reg mock.Registry
// 	ws := reg.NewMockSpace()
// 	// input 1 has two atoms
// 	desc := mock.MakeDesc("mockType", common.expandedInputs)
// 	wsBlock := mock.CreateBlock("mock", desc)
// 	// connect all of wsBlock's inputs to fake connections of dst
// 	dst := mock.CreateBlock("dst", nil)
// 	wsBlock.Workspace, dst.Workspace = ws, ws // workspace is used to test for valid connections
// 	connections := make(map[string]*mock.MockConnection)
// 	bound := make(map[string]string)
// 	for i := 0; i < wsBlock.NumInputs(); i++ {
// 		in := wsBlock.Input(i)
// 		if inc := in.Connection(); inc != nil {
// 			zeroIndexed := strconv.Itoa(i)
// 			tgt := &mock.MockConnection{
// 				Name:   "c" + zeroIndexed,
// 				Source: dst,
// 			}
// 			bound[in.InputName()] = tgt.Name
// 			inc.Connect(tgt)
// 			connections[tgt.Name] = tgt
// 		}
// 	}
// 	require.Equal(t, bound, map[string]string{
// 		block.Scope("a", "M1", "0", "TERM"):  "c1",
// 		block.Scope("a", "M1", "1", "TERM"):  "c2",
// 		block.Scope("a", "M3", "0", "TERM"):  "c6",
// 		block.Scope("a", "M3", "2", "STATE"): "c9",
// 	})
// 	// create a matching mui block
// 	muispace := reg.NewMockSpace()
// 	blocks := mutant.NewMutatedBlocks()
// 	db := &mock.MockDatabase{common.atomProducts}
// 	arch := mock.NewMutations(common.inputs, common.quarks)
// 	require.NoError(t, arch.Preregister("mockType", &reg))
// 	//
// 	mb := blocks.CreateMutatedBlock(wsBlock, arch, db)
// 	for k, v := range common.inputAtoms {
// 		mutant.AtomizedInputs{v }
// 		mb.GetMutatedInput(k).Atoms= SetAtomsForInput(k, v)
// 	}
// 	muiContainer, e := mb.CreateMui(muispace)
// 	require.NoError(t, e)
// 	//
// 	initialConnections := []string{
// 		block.Scope("a", "M1", "0", "TERM:c1"),
// 		block.Scope("a", "M1", "1", "TERM:c2"),
// 		block.Scope("a", "M3", "0", "TERM:c6"),
// 		block.Scope("a", "M3", "2", "STATE:c9")}
// 	require.Equal(t, initialConnections, inputConnections(wsBlock),
// 		"initial connections from block")
// 	cs := mb.SaveConnections(muiContainer)
// 	generatedConnections := map[string]string{
// 		// the mock block ids are (blockType#countOfBlocks)
// 		// the block type for quarks are: (mui, inputName, quarkName)
// 		// the atoms come from common.inputAtoms
// 		block.Scope("mui", "M1", "A#0"): "c1",
// 		block.Scope("mui", "M1", "A#1"): "c2", // <-- its making c1
// 		block.Scope("mui", "M2", "B#0"): "???",
// 		block.Scope("mui", "M3", "A#0"): "c6",
// 		block.Scope("mui", "M3", "B#0"): "c6", // <-- why?
// 		block.Scope("mui", "M3", "C#0"): "c9", // <--- Zits making c6
// 	}
// 	require.Equal(t, generatedConnections, savedConnections(cs),
// 		"initial connections from save")
// 	// rotate the first connections, leave the last alone
// 	// note: we cant delete a connection without actually deleting the input
// 	// we test that at the end by restoring into an empty block
// 	makeConnection := func(c string) []block.Connection {
// 		x := connections[c]
// 		return []block.Connection{x}
// 	}
// 	remap := mutant.SavedConnections{
// 		block.Scope("a", "M1", "0", "TERM"):  makeConnection("c6"),
// 		block.Scope("a", "M1", "1", "TERM"):  makeConnection("c2"),
// 		block.Scope("a", "M3", "0", "TERM"):  makeConnection("c1"),
// 		block.Scope("a", "M3", "2", "STATE"): makeConnection("c9"),
// 	}
// 	mb.UpdateConnections(remap)
// 	require.NoError(t, mb.CreateFromMui(muiContainer))

// 	// test the connections by walking the block
// 	remapped := []string{
// 		block.Scope("a", "M1", "0", "TERM:c6"),
// 		block.Scope("a", "M1", "1", "TERM:c2"),
// 		block.Scope("a", "M3", "0", "TERM:c1"),
// 		block.Scope("a", "M3", "2", "STATE:c9"),
// 	}
// 	require.Equal(t, remapped, inputConnections(wsBlock),
// 		"remapped connections")

// 	// finally, restore into an empty block
// 	// emptyBlock := mock.CreateBlock("mock", mock.MakeDesc("mockType", common.inputs))
// 	// remap.RestoreConnections(emptyBlock)
// 	// require.Empty(t, inputConnections(emptyBlock),
// 	// 	"restored connections into an empty block")
// }

// // list the names of the block's input connection
// func inputConnections(b block.Shape) (ret []string) {
// 	for i := 0; i < b.NumInputs(); i++ {
// 		in := b.Input(i)
// 		if inc := in.Connection(); inc != nil {
// 			s := inc.(*mock.MockConnection).String()
// 			ret = append(ret, s)
// 		}
// 	}
// 	return
// }

// // returns the target connections, and the sources they point back to
// func savedConnections(cs mutant.SavedConnections) (ret map[string]string) {
// 	ret = make(map[string]string)
// 	for k, ca := range cs {
// 		tgt := "???"
// 		if len(ca) > 0 {
// 			cel := ca[0].(*mock.MockConnection)
// 			tgt = cel.Name
// 		}
// 		ret[k] = tgt
// 	}
// 	return
// }
