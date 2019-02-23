package gblocks

// import (
// 	// "github.com/kr/pretty"
// 	"github.com/stretchr/testify/require"
// 	// r "reflect"
// 	//	"github.com/ionous/gblocks/decor"
// 	"testing"
// )

// type DecorTest struct {
// 	Mutant DecorMutation
// }

// type DecorMutation struct {
// 	Text          string         `decor:"commasAnd"`
// 	NextStatement *DecorMutation `decor:"fullStop"`
// }

// // func commasAnd(d decor.Context) (ret string) {
// // 	ctx := d.Container()
// // 	// potentially, we could do this by running multiple decorations:
// // 	// first the comma has next, then the is last.
// // 	if ctx.HasPrev() && ctx.IsLast() {
// // 		ret = ", and"
// // 	} else if ctx.HasNext() {
// // 		ret = ","
// // 	}
// // 	return
// // }

// // func fullStop(d decor.Context) (ret string) {
// // 	if !ctx.IsConnected() {
// // 		ret = "."
// // 	}
// // 	return
// // }

// // really just testing MUTANT/0 part of toolbox construction
// func TestDecorToolbox(t *testing.T) {
// 	three := NewTool(&DecorTest{DecorMutation{
// 		Text: "one", NextStatement: &DecorMutation{
// 			Text: "two", NextStatement: &DecorMutation{
// 				Text: "three"}}}})

// 	expected :=
// 		`<block type="decor_test">` +
// 			/**/ `<mutation>` +
// 			/* */ `<atoms name="MUTANT">` +
// 			/*  */ `<atom type="decor_mutation"/>` +
// 			/*  */ `<atom type="decor_mutation"/>` +
// 			/* */ `</atoms>` +
// 			/**/ `</mutation>` +
// 			/**/ `<field name="MUTANT/0/TEXT">one</field>` +
// 			/**/ `<field name="MUTANT/1/TEXT">two</field>` +
// 			/**/ `<field name="MUTANT/2/TEXT">three</field>` +
// 			`</block>`
// 	require.Equal(t, expected, three.OuterHTML(), "toolbox")
// }

// func TestDecorInputs(t *testing.T) {
// 	testDecor(t, func(ws *Workspace, reg *Registry, b *Block) {
// 		x := reduceInputs(b)
// 		expected := []string{"MUTANT",
// 			"$decor/MUTANT/0/TEXT", "MUTANT/0/TEXT", "$decor/MUTANT/0/NEXT_STATEMENT",
// 			"$decor/MUTANT/1/TEXT", "MUTANT/1/TEXT", "$decor/MUTANT/1/NEXT_STATEMENT",
// 			"$decor/MUTANT/2/TEXT", "MUTANT/2/TEXT", "$decor/MUTANT/2/NEXT_STATEMENT",
// 		}
// 		require.Equal(t, expected, x, "inputs")
// 	})
// }

// // func TestDecorList(t *testing.T) {
// // 	testDecor(t, func(ws *Workspace, reg *Registry, b *Block) {
// // 	})

// // // var d decor.Registry
// // // d.Register("commasAnd", commasAnd)
// // // d.Register("fullStop", fullStop)
// // // // one.
// // // // one, and two.
// // // // one, two, and three.
// // // one := &DecorTest{&DecorMutation{"one", nil}}

// // // two := &DecorTest{&DecorMutation{"one", &DecorMutation{"two", nil}}}
// // // three := &DecorTest{&DecorMutation{"one", &DecorMutation{"two", &DecorMutation{"three", nil}}}}
// // }

// func testDecor(t *testing.T, fn func(*Workspace, *Registry, *Block)) {
// 	var reg Registry
// 	require.NoError(t,
// 		reg.RegisterMutation((*DecorMutation)(nil),
// 			Mutation{"decor", (*DecorMutation)(nil)},
// 		), "register mutation")
// 	require.NoError(t, reg.RegisterBlocks(nil,
// 		(*DecorMutation)(nil),
// 		(*DecorTest)(nil),
// 	), "register blocks")
// 	ws := NewBlankWorkspace(false,&orderedGenerator{name: "main"})w
// 	three := NewTool(&DecorTest{DecorMutation{
// 		Text: "one", NextStatement: &DecorMutation{
// 			Text: "two", NextStatement: &DecorMutation{
// 				Text: "three"}}}})
// 	if b := blockly.Xml.DomToBlock(three, ws); b == nil {
// 		t.Fatal("no block")
// 	} else {
// 		fn(ws, &reg, b)
// 	}
// }
