package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/stretchr/testify/require"
	"testing"
)

type StackBlock struct {
	PreviousStatement,
	NextStatement interface{}
}

type RowBlock struct {
	Input interface {
		Output() interface{}
	}
}

// implement a generic output
func (b *RowBlock) Output() interface{} {
	return b
}

func testMirror(t *testing.T, fn func(*Workspace)) {
	var reg Registry
	require.NoError(t, reg.RegisterBlock((*StackBlock)(nil), nil), "register stack")
	require.NoError(t, reg.RegisterBlock((*RowBlock)(nil), nil), "register row")
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

func TestJsonStackBlock(t *testing.T) {
	var reg Registry
	opts := make(Options)
	require.NoError(t, reg.RegisterBlock((*StackBlock)(nil), opts), "register stack")
	expected := Options{
		"type":              "stack_block",
		"message0":          "",
		"previousStatement": nil,
		"nextStatement":     nil,
	}
	require.Equal(t, expected, opts)
}

func TestJsonRowBlock(t *testing.T) {
	var reg Registry
	opts := make(Options)
	require.NoError(t, reg.RegisterBlock((*RowBlock)(nil), opts), "register row")
	expected := Options{
		"type":     "row_block",
		"message0": "%1",
		"args0": []Options{{
			"type": "input_value",
			"name": "INPUT",
		}},
		"output": nil,
	}
	require.Equal(t, expected, opts)
}

// verify that the workspace gets default data created, destroyed for an appropriate block
// see also: blockly/tests/jsunit/block_test.js setupStackBlocks.
func TestMirrorStackCreate(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		blockA := ws.NewBlock((*StackBlock)(nil))
		require.NotNil(t, blockA, "first block")
		//
		data := ws.GetDataById(blockA.Id)
		require.NotNil(t, data, "expected non-nil data")
		require.IsType(t, (*StackBlock)(nil), data, "expected stack data")
		//
		blockA.Dispose(false)
		again := ws.GetDataById(blockA.Id)
		require.Nil(t, again, "data should be gone after delete")
	})
}

// verify that stacking stack mirrors next / prev connections in go
func TestMirrorStackConnect(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var stack [3](*Block)
		var data [3](*StackBlock)
		for i := 0; i < len(stack); i++ {
			b := ws.NewBlock((*StackBlock)(nil))
			require.NotNilf(t, b, "new block %", i)
			stack[i] = b
			d := ws.GetDataById(b.Id).(*StackBlock)
			require.NotNilf(t, d, "get data %d", i)
			require.Nilf(t, d.PreviousStatement, "prev empty %d", i)
			require.Nilf(t, d.NextStatement, "next empty %d", i)
			data[i] = d
		}

		// a->b->c
		stack[0].NextConnection.Connect(stack[1].PreviousConnection)
		stack[1].NextConnection.Connect(stack[2].PreviousConnection)

		prev := [3]interface{}{nil, data[0], data[1]}
		next := [3]interface{}{data[1], data[2], nil}

		for i := 0; i < len(data); i++ {
			d, p, n := data[i], prev[i], next[i]
			require.Equalf(t, p, d.PreviousStatement, "data prev %d", i)
			require.Equalf(t, n, d.NextStatement, "data next %d", i)
		}

		for i := 0; i < len(stack); i++ {
			d := ws.GetDataById(stack[i].Id).(*StackBlock)
			require.Truef(t, d == data[i], "get again %d")
		}
	})
}

func TestMirrorStackDisconnectHeal(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var stack [3](*Block)
		var data [3](*StackBlock)
		for i := 0; i < len(stack); i++ {
			b := ws.NewBlock((*StackBlock)(nil))
			data[i] = ws.GetDataById(b.Id).(*StackBlock)
			stack[i] = b
		}

		// a->b->c
		stack[0].NextConnection.Connect(stack[1].PreviousConnection)
		stack[1].NextConnection.Connect(stack[2].PreviousConnection)

		// heal the rift
		stack[1].Dispose(true)
		{
			prev := [3]interface{}{nil, nil, data[0]}
			next := [3]interface{}{data[2], nil, nil}

			for i := 0; i < 3; i++ {
				d, p, n := data[i], prev[i], next[i]
				require.Equalf(t, p, d.PreviousStatement, "data prev %d", i)
				require.Equalf(t, n, d.NextStatement, "data next %d", i)
			}
		}
	})
}

func TestMirrorStackDisconnectBreak(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var stack [3](*Block)
		for i := 0; i < len(stack); i++ {
			stack[i] = ws.NewBlock((*StackBlock)(nil))
		}

		// a->b->c
		stack[0].NextConnection.Connect(stack[1].PreviousConnection)
		stack[1].NextConnection.Connect(stack[2].PreviousConnection)

		// dont heal the rift
		stack[1].Dispose(false)
		{
			d0 := ws.GetDataById(stack[0].Id).(*StackBlock)
			require.NotNil(t, d0, "data deleted 0")
			require.Equal(t, nil, d0.PreviousStatement, "data prev 0")
			require.Equal(t, nil, d0.NextStatement, "data next 0")

			d1 := ws.GetDataById(stack[1].Id)
			require.Nil(t, d1, "data deleted 1")

			// the second one is disposed of too
			// because it isnt re-attached.
			d2 := ws.GetDataById(stack[2].Id)
			require.Nil(t, d2, "data deleted 2")
		}
	})
}

func TestMirrorRowCreate(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var row [3](*Block)
		for i := 0; i < len(row); i++ {
			row[i] = ws.NewBlock((*RowBlock)(nil))
		}

		row[0].GetInput(0).Connection.Connect(row[1].OutputConnection)
		p1 := row[1].GetParent()
		require.NotNil(t, p1, "connected 1")
		require.Equal(t, row[0].Id, p1.Id, "connected 1<-0")

		row[1].GetInput(0).Connection.Connect(row[2].OutputConnection)
		p2 := row[2].GetParent()
		require.NotNil(t, p2, "connected 2")
		require.Equal(t, row[1].Id, p2.Id, "connected 2<-1")
	})
}

// verify that linking inputs mirrors connections in go
func TestMirrorRowConnect(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
		var row [3](*Block)
		var data [3](*RowBlock)
		for i := 0; i < len(row); i++ {
			row[i] = ws.NewBlock((*RowBlock)(nil))
			data[i] = ws.GetDataById(row[i].Id).(*RowBlock)
		}
		row[0].GetInput(0).Connection.Connect(row[1].OutputConnection)
		row[1].GetInput(0).Connection.Connect(row[2].OutputConnection)

		require.Equal(t, data[0].Input, data[1], "mirrored 0->1")
		require.Equal(t, data[1].Input, data[2], "mirrored 1->2")
		require.Equal(t, data[2].Input, nil, "no change 2->")
	})
}

func TestMirrorRowDiscnnect(t *testing.T) {
	testMirror(t, func(ws *Workspace) {
	})
}

// verify that setting fields in a block mirrors values in go.
func TestMirrorFields(t *testing.T) {

	testMirror(t, func(ws *Workspace) {
	})
}
