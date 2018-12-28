package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/stretchr/testify/require"
	"testing"
)

// see also: blockly/tests/jsunit/block_test.js setupStackBlocks.
type StackBlock struct {
	PreviousStatement,
	NextStatement interface{}
}

// see also: blockly/tests/jsunit/block_test.js setUpRowBlocks.
type RowBlock struct {
	Input interface {
		Output() interface{}
	}
}

type FieldBlock struct {
	Number float32
}

// Output - implement a generic output
func (b *RowBlock) Output() interface{} {
	return b
}

func testMirror(t *testing.T, fn func(*Workspace)) {
	var reg Registry
	require.NoError(t, reg.RegisterBlock((*StackBlock)(nil), nil), "register stack")
	require.NoError(t, reg.RegisterBlock((*RowBlock)(nil), nil), "register row")
	require.NoError(t, reg.RegisterBlock((*FieldBlock)(nil), nil), "register fields")
	ws := NewBlankWorkspace(&reg)
	// replace timed event queue with direct event queue
	events := &Events{Object: js.Global.Get("Blockly").Get("Events")}
	events.Set("fire", js.MakeFunc(func(_ *js.Object, args []*js.Object) interface{} {
		events.TestFire(args[0])
		return nil
	}))
	//
	fn(ws)
	ws.Dispose()
}

// TestJsonStackBlock - verify the json generated by gblocks reflection matches expections.
func TestJsonStackBlock(t *testing.T) {
	var reg Registry
	opts := make(Options)
	require.NoError(t, reg.RegisterBlock((*StackBlock)(nil), opts), "register stack")
	expected := Options{
		"type":              TypeName("stack_block"),
		"previousStatement": nil,
		"nextStatement":     nil,
	}
	require.Equal(t, expected, opts)
}

// TestJsonRowBlock - verify the json generated by gblocks reflection matches expections.
func TestJsonRowBlock(t *testing.T) {
	var reg Registry
	opts := make(Options)
	require.NoError(t, reg.RegisterBlock((*RowBlock)(nil), opts), "register row")
	expected := Options{
		"type":     TypeName("row_block"),
		"message0": "%1",
		"args0": []Options{{
			"type": "input_value",
			"name": "INPUT",
		}},
		"output": nil,
	}
	require.Equal(t, expected, opts)
}

// TestJsonFieldBlock - verify the json generated by gblocks reflection matches expections.
func TestJsonFieldBlock(t *testing.T) {
	var reg Registry
	opts := make(Options)
	require.NoError(t, reg.RegisterBlock((*FieldBlock)(nil), opts), "register stack")
	expected := Options{
		"type":     TypeName("field_block"),
		"message0": "%1",
		"args0": []Options{{
			"name": "NUMBER",
			"type": "field_number",
		}},
	}
	require.Equal(t, expected, opts)
}

// TestMirrorStackCreate -
// verify appropriate gblock data created/destroyed for an appropriate block
func TestMirrorStackCreate(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		blockA, e := ws.NewBlock((*StackBlock)(nil))
		require.NoError(t, e)
		require.NotNil(t, blockA, "first block")
		//
		data := ws.GetDataById(blockA.Id)
		require.NotNil(t, data, "expected non-nil data")
		require.IsType(t, (*StackBlock)(nil), data, "expected block data")
		//
		blockA.Dispose()
		again := ws.GetDataById(blockA.Id)
		require.Nil(t, again, "data should be gone after delete")
	})
}

// TestMirrorStackConnect -
// verify that stacking block mirrors next / prev connections in go
func TestMirrorStackConnect(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var block [3](*Block)
		var data [3](*StackBlock)

		t.Log("building block data")
		for i := 0; i < len(block); i++ {
			b, e := ws.NewBlock((*StackBlock)(nil))
			require.NoError(t, e)
			require.NotNilf(t, b, "new block %d", i)
			block[i] = b
			d := ws.GetDataById(b.Id).(*StackBlock)
			require.NotNilf(t, d, "get data %d", i)
			require.Nilf(t, d.PreviousStatement, "prev empty %d", i)
			require.Nilf(t, d.NextStatement, "next empty %d", i)
			data[i] = d
		}

		//
		t.Log("connecting a->b->c")
		block[0].NextConnection.Connect(block[1].PreviousConnection)
		block[1].NextConnection.Connect(block[2].PreviousConnection)

		prev := [3]interface{}{nil, data[0], data[1]}
		next := [3]interface{}{data[1], data[2], nil}

		t.Log("testing connections")
		for i := 0; i < len(data); i++ {
			d, p, n := data[i], prev[i], next[i]
			require.Equalf(t, p, d.PreviousStatement, "data prev %d", i)
			require.Equalf(t, n, d.NextStatement, "data next %d", i)
		}

		t.Log("verifying data")
		for i := 0; i < len(block); i++ {
			d := ws.GetDataById(block[i].Id).(*StackBlock)
			require.Truef(t, d == data[i], "get again %d")
		}
	})
}

// TestMirrorStackDisconnectHeal -
// after stacking three blocks, delete the middle one, patching the first and third together.
// verify that the mirrored gblock data follows along.
func TestMirrorStackDisconnectHeal(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var block [3](*Block)
		var data [3](*StackBlock)
		for i := 0; i < len(block); i++ {
			b, e := ws.NewBlock((*StackBlock)(nil))
			require.NoError(t, e)
			data[i] = ws.GetDataById(b.Id).(*StackBlock)
			block[i] = b
		}

		// a->b->c
		block[0].NextConnection.Connect(block[1].PreviousConnection)
		block[1].NextConnection.Connect(block[2].PreviousConnection)

		// heal the rift
		block[1].Unplug(true)
		block[1].Dispose()
		{
			prev := [3]interface{}{nil, nil, data[0]}
			next := [3]interface{}{data[2], nil, nil}
			exists := [3]bool{true, false, true}

			for i := 0; i < 3; i++ {
				d, p, n := data[i], prev[i], next[i]
				require.Equalf(t, p, d.PreviousStatement, "data prev %d", i)
				require.Equalf(t, n, d.NextStatement, "data next %d", i)
				require.Equalf(t, exists[i], ws.GetDataById(block[i].Id) != nil, "data not deleted %d", i)
			}
		}
	})
}

// TestMirrorStackDisconnectBreak -
// after stacking three blocks, delete the middle one and its child blocks
// verify that the mirrored gblock data follows along.
func TestMirrorStackDisconnectBreak(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var block [3](*Block)
		for i := 0; i < len(block); i++ {
			if b, e := ws.NewBlock((*StackBlock)(nil)); e != nil {
				require.NoError(t, e)
			} else {
				block[i] = b
			}
		}

		// a->b->c
		block[0].NextConnection.Connect(block[1].PreviousConnection)
		block[1].NextConnection.Connect(block[2].PreviousConnection)

		require.NotNil(t, ws.GetDataById(block[1].Id), "data exists 1")
		require.NotNil(t, ws.GetDataById(block[2].Id), "data exists 2")

		// dont heal the rift
		block[1].Dispose()
		{
			d0 := ws.GetDataById(block[0].Id).(*StackBlock)
			require.NotNil(t, d0, "data deleted 0")
			require.Equal(t, nil, d0.PreviousStatement, "data prev 0")
			require.Equal(t, nil, d0.NextStatement, "data next 0")

			// the second one is disposed of too
			// because it isnt re-attached.
			require.Nil(t, ws.GetDataById(block[1].Id), "data deleted 1")
			require.Nil(t, ws.GetDataById(block[2].Id), "data deleted 2")
		}
	})
}

// TestMirrorRowCreate -- create three blocks which connect in a single row.
func TestMirrorRowCreate(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var block [3](*Block)
		for i := 0; i < len(block); i++ {
			if b, e := ws.NewBlock((*RowBlock)(nil)); e != nil {
				require.NoError(t, e)
			} else {
				block[i] = b
			}
		}

		block[0].Input(0).Connection.Connect(block[1].OutputConnection)
		p1 := block[1].GetParent()
		require.NotNil(t, p1, "connected 1")
		require.Equal(t, block[0].Id, p1.Id, "connected 1<-0")

		block[1].Input(0).Connection.Connect(block[2].OutputConnection)
		p2 := block[2].GetParent()
		require.NotNil(t, p2, "connected 2")
		require.Equal(t, block[1].Id, p2.Id, "connected 2<-1")
	})
}

// TestMirrorRowConnect -- create three blocks which connect in a single row:
// verify that the gblock data follows along.
func TestMirrorRowConnect(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var block [3](*Block)
		var data [3](*RowBlock)
		for i := 0; i < len(block); i++ {
			if b, e := ws.NewBlock((*RowBlock)(nil)); e != nil {
				require.NoError(t, e)
			} else {
				block[i] = b
				data[i] = ws.GetDataById(block[i].Id).(*RowBlock)
			}
		}
		block[0].Input(0).Connection.Connect(block[1].OutputConnection)
		block[1].Input(0).Connection.Connect(block[2].OutputConnection)

		require.Equal(t, data[0].Input, data[1], "mirrored 0->1")
		require.Equal(t, data[1].Input, data[2], "mirrored 1->2")
		require.Equal(t, data[2].Input, nil, "no change 2->")
	})
}

// TestMirrorRowDisconnectHeal -
// after joining three blocks in a row, delete the middle one, patching the first and third together.
// verify that the mirrored gblock data follows along.
func TestMirrorRowDisconnectHeal(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var block [3](*Block)
		var data [3](*RowBlock)
		for i := 0; i < len(block); i++ {
			if b, e := ws.NewBlock((*RowBlock)(nil)); e != nil {
				require.NoError(t, e)
			} else {
				block[i] = b
				data[i] = ws.GetDataById(block[i].Id).(*RowBlock)
			}

		}
		block[0].Input(0).Connection.Connect(block[1].OutputConnection)
		block[1].Input(0).Connection.Connect(block[2].OutputConnection)

		// heal the rift
		block[1].Unplug(true)
		block[1].Dispose()
		{
			next := [3]interface{}{data[2], nil, nil}
			exists := [3]interface{}{true, false, true}
			for i := 0; i < 3; i++ {
				d, n := data[i], next[i]
				require.Equalf(t, n, d.Input, "data next %d", i)
				require.Equalf(t, exists[i], ws.GetDataById(block[i].Id) != nil, "data not deleted %d", i)
			}
		}
	})
}

// TestMirrorStackDisconnectBreak -
// after joining three blocks in a row, delete the middle one and its children (right-side) blocks.
// verify that the mirrored gblock data follows along.
func TestMirrorRowDisconnectBreak(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var block [3](*Block)
		for i := 0; i < len(block); i++ {
			if b, e := ws.NewBlock((*RowBlock)(nil)); e != nil {
				require.NoError(t, e)
			} else {
				block[i] = b
			}
		}
		block[0].Input(0).Connection.Connect(block[1].OutputConnection)
		block[1].Input(0).Connection.Connect(block[2].OutputConnection)

		require.NotNil(t, ws.GetDataById(block[1].Id), "data exists 1")
		require.NotNil(t, ws.GetDataById(block[2].Id), "data exists 2")

		// dont heal the rift
		block[1].Dispose()
		{
			d0 := ws.GetDataById(block[0].Id).(*RowBlock)
			require.NotNil(t, d0, "data not deleted")
			require.Equal(t, nil, d0.Input, "data input")

			// the second one is disposed of too
			// because it isnt re-attached.
			require.Nil(t, ws.GetDataById(block[1].Id), "data deleted 1")
			require.Nil(t, ws.GetDataById(block[2].Id), "data deleted 2")
		}
	})
}

// verify that setting fields in a block mirrors values in go.
func TestMirrorFields(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		block, e := ws.NewBlock((*FieldBlock)(nil))
		require.NoError(t, e)
		f := block.GetField("NUMBER")
		require.NotNil(t, f, "get field")
		f.SetText("42")
		data := ws.GetDataById(block.Id)
		require.NotNil(t, data, "get data")
		require.IsType(t, (*FieldBlock)(nil), data, "expected field block")
		res := data.(*FieldBlock)
		require.Equal(t, float32(42.0), res.Number, "expected life, the universe and everything")
	})
}
