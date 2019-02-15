package gblocks

import (
	// 	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/decor"
	// 	"github.com/ionous/gblocks/named"
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

// // blockContext
// type blockContext struct {
// 	parentContext
// 	block *Block
// }

// func (*blockContext) ContextType() decor.ContextType {
// 	return decor.BlockContext
// }
// func (ctx *blockContext) Parent() decor.Context {
// 	return nil
// }

// // mutationContext
// type mutationContext struct {
// 	parentContext
// 	mi *InputMutation
// }

// func (*mutationContext) ContextType() decor.ContextType {
// 	return decor.MutationContext
// }

// // fix:? should mutations should check next/prev for other mutations.
// func (ctx *mutationContext) Parent() decor.Context {
// 	block := ctx.mi.Input().Block()
// 	return &blockContext{parentContext{}, block}
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
// 	return &blockContext{parentContext{}, block}
// }

// func (ctx *atomContext) HasPrev() bool {
// 	return ctx.atom > 0
// }

// func (ctx *atomContext) HasNext() decor.Context {
// 	return ctx.atom < ctx.mi.NumAtoms
// }

// //

// func (ctx *subContext) ContextType() (ret decor.ContextType) {

// }

// // // func redecorateInput(reg decor.Registry, container named.Type, in *Input) {
// // // 	// although things like "text input" and "number inputs" are "fields"
// // // 	// interpolate arranges each into its own dummy input.
// // // 	// ( otherwise this will get super complicated.)
// // // 	var decoration *FieldLabel
// // // 	for fi, fcnt := 0, fields.Length(); fi < fcnt; fi++ {
// // // 		field := fields.Field(fi)
// // // 		if fieldName := field.Name(); fieldName == ItemDecor {
// // // 			decoration = &FieldLabel{field}
// // // 		}
// // // 	}
// // // 	parts := strings.Split(in.Name.String(), "/")
// // // 	bareName := parts[len(parts)-1]

// // // 	if fn, ok := reg.Find(container, named.Item(bareName)); ok {
// // // 	}
// // // }

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

func (b *Block) redecorate(reg decor.Decorators) (err error) {
	return
}

// 	blockCtx := &blockContext{block: b}
// 	var parent Context = blockCtx

// 	for inputIndex, inputCnt := 0, b.NumInputs(); inputIndex < inputCnt; inputIndex++ {
// 		in := b.Input(inputIndex)
// 		fields := in.Fields()
// 		for fieldIndex, fieldCnt := 0, fields.Length(); fieldIndex < fieldCnt; fieldIndex++ {
// 			field := fields.Field(inputIndex)
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
// 					} else if isInput := itemIndex < 0; !isInput {
// 						// a field; parent is block
// 					} else if mi := in.Mutation(); mi == nil {
// 						// a mutation input; parent is block
// 					} else {
// 						// a normal input; parent is block
// 					}

// 				default:
// 					err = errutil.Append(err, errutil.New("unexpected path", fieldName))
// 				}
// 			}
// 		}
// 	}
// }
