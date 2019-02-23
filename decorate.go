package gblocks

import (
// 	"github.com/ionous/errutil"
// "github.com/ionous/gblocks/decor"
// 	"github.com/ionous/gblocks/block"
// 	"strings"
)

// // partially implements decor.Context
// type parentContext struct {
// 	parent decor.Context
// }

// type subContext struct {
// 	parentContext
// 	atomIndex  int
// 	fieldIndex int
// 	prevIndex  int
// }

// func (parentContext) String() string    { return "" }
// func (parentContext) IsConnected() bool { return false }
// func (parentContext) HasPrev() bool     { return false }
// func (parentContext) HasNext() bool     { return false }

// // mutableBlock
// type mutableBlock struct {
// 	parentContext
// 	block *Block
// }

// func (*mutableBlock) ContextType() decor.ContextType {
// 	return decor.BlockContext
// }
// func (ctx *mutableBlock) Parent() decor.Context {
// 	return nil
// }

// // mutableInput
// type mutableInput struct {
// 	parentContext
// 	mi *InputMutation
// }

// func (*mutableInput) ContextType() decor.ContextType {
// 	return decor.MutationContext
// }

// // fix:? should mutations should check next/prev for other mutations.
// func (ctx *mutableInput) Parent() decor.Context {
// 	block := ctx.mi.Input().Block()
// 	return &mutableBlock{parentContext{}, block}
// }

// // atomContext
// type atomContext struct {
// 	parentContext
// 	mi   *InputMutation
// 	atom int
// }

// func (*atomContext) ContextType() decor.ContextType {
// 	return decor.AtomContext
// }

// // fix:? should mutations should check next/prev for other mutations.
// func (ctx *atomContext) Parent() decor.Context {
// 	block := ctx.mi.Input().Block()
// 	return &mutableBlock{parentContext{}, block}
// }

// func (ctx *atomContext) HasPrev() bool {
// 	return ctx.atom > 0
// }

// func (ctx *atomContext) HasNext() decor.Context {
// 	return ctx.atom < ctx.mi.NumAtoms
// }

// //

// mutableInput
// type mutableInput struct {
// 	parentContext
// 	mi *Input
// }

// func (*mutableInput) ContextType() decor.ContextType {
// 	return decor.ItemContext
// }

// func (ctx *atomContext) HasPrev() bool {
// 	return ctx.atom > 0
// }

// func (ctx *atomContext) HasNext() decor.Context {
// 	return ctx.atom < ctx.mi.NumAtoms
// }

// // // // func redecorateInput(reg decor.Registry, container block.Type, in *Input) {
// // // // 	// although things like "text input" and "number inputs" are "fields"
// // // // 	// interpolate arranges each into its own dummy input.
// // // // 	// ( otherwise this will get super complicated.)
// // // // 	var decoration *FieldLabel
// // // // 	for fi, fcnt := 0, fields.Length(); fi < fcnt; fi++ {
// // // // 		field := fields.Field(fi)
// // // // 		if fieldName := field.Name(); fieldName == ItemDecor {
// // // // 			decoration = &FieldLabel{field}
// // // // 		}
// // // // 	}
// // // // 	parts := strings.Split(in.Name.String(), "/")
// // // // 	bareName := parts[len(parts)-1]

// // // // 	if fn, ok := reg.Find(container, block.Item(bareName)); ok {
// // // // 	}
// // // // }

// // return the field index of the indicated item, -1 if the item is an input; or false, if not found.
// func findItem(b *Block, inputIndex, fieldIndex int, itemPath string) (ret int, okay bool) {
// 	if pathOnly == in.Name {
// 		ret, okay = -1, true
// 	} else if nextIndex := fieldIndex + 1; nextIndex < fieldCnt {
// 		nextField := fields.Field(nextIndex)
// 		if nextField.Name() == pathOnly {
// 			ret, okay = nextIndex, true
// 		}
// 	}
// 	return
// }

// func (b *Block) redecorate(reg decor.Decorators) (err error) {
// 	blockCtx := &mutableBlock{block: b}
// 	var parent Context = blockCtx

// 	// every input
// 	for inputIndex, inputCnt := 0, b.NumInputs(); inputIndex < inputCnt; inputIndex++ {
// 		in := b.Input(inputIndex)
// 		fields := in.Fields()
// 		// every field in that input
// 		for fieldIndex, fieldCnt := 0, fields.Length(); fieldIndex < fieldCnt; fieldIndex++ {
// 			field := fields.Field(inputIndex)
// 			// filtering for decorations
// 			if fieldName := field.Name(); strings.HasPrefix(ItemDecor) {
// 				label := &FieldLabel{field}
// 				pathOnly := fieldName[len(ItemDecor):]
// 				parts := makePathParts(pathOnly) // ex. $decor/MUTANT/0/TEXT

// 				switch len(parts) {
// 				case 3: // ex. "MUTANT/0/TEXT", an item in a mutation or atom.
// 					atom, item := strconv.Atoi(parts[1]), parts[2]

// 					if item == StatementNext {
// 						// points to the next atom in a mutation input; parent is the mutation
// 						// we can use item index as -1?
// 					} else if itemIndex, ok := findItem(b, inputIndex, fieldIndex, itemPath); ok {
// 						// points to an item in an atom or the mutation
// 						if atom == 0 {
// 							// parent is the mutation
// 						} else {
// 							// parent is the atom
// 						}
// 					} else {
// 						err = errutil.Append(err, errutil.New("item not found", fieldName))
// 					}

// 				case 1: // ex. "MUTANT", an input, or field.
// 					if itemIndex, ok := findItem(b, inputIndex, fieldIndex, itemPath); !ok {
// 						err = errutil.Append(err, errutil.New("item not found", fieldName))
// 					} else if itemIsField := itemIndex >= 0; itemIsField {
// 						// a field; parent is block
// 					} else if mi := in.Mutation(); mi == nil {
// 						// a mutation input; parent is block
// 					} else {
// 						//m := &InputMutation{Object: m}
// 						// a normal input; parent is block
// 						// siblings are somewhat difficult:
// 						// prev would be a -- itd be a item, a field or an input -
// 						// pretty much we'd need to re/create the stack here
// 						// you pretty much will have to anyway --
// 						// an option would be to create the tree
// 						// then visit the tree.
// 						// each node would have a pointer to the decoration label.
// 						// there'd be a current parent, youd build the node -- or notice that you are done
// 						// this is not particularly "fun" -- noticing would be some pain, evaluting what kind of node
// 						// you'd need like a statemachine of some sort. ex. state: parsing list of atoms, item is not an atom --
// 						// so youd need triggers, and an alg to handle those triggers.
// 						// if you had your inspect in advance -- then youd run switches
// 						// youd still have to build your instance tree.
// 					}

// 				default:
// 					err = errutil.Append(err, errutil.New("unexpected path", fieldName))
// 				}
// 			}
// 		}
// 	}
// }
