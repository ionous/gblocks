package gblocks

import (
	// 	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/decor"
	// 	"github.com/ionous/gblocks/named"
	// 	"strings"
)

// type parentContext struct{ parent Context }

// func (p parentContext) Parent() Context { return p.parent }
// func (parentContext) IsConnected() bool { return false }
// func (parentContext) HasPrev() bool     { return false }
// func (parentContext) HasNext() bool     { return false }

// // inset implements decor.Contexts
// type blockContext struct {
// 	parentContext
// 	block *Block
// }

// // inset implements decor.Contexts
// type inputContext struct {
// 	parentContext
// 	block *Block
// 	input *Input
// }

// type fieldContext struct {
// 	parentContext
// 	block *Block
// 	input *Input
// 	index int // field row
// }

// type mutationContext struct {
// 	parentContext
// 	mutation *InputMutation
// }

// type atomContext struct {
// 	parentContext
// 	mutation *InputMutation
// 	index    int // atoms
// }

// func (parentContext) IsConnected() bool { return false }
// func (parentContext) HasPrev() bool     { return false }
// func (parentContext) HasNext() bool     { return false }

// type outputContext struct {
// 	parentContext
// 	connects *Connection
// }

// func (ctx *outputContext) IsConnected() bool {
// 	return ctx.connects.IsConnected()
// }

// func redecorateInput(reg decor.Registry, container named.Type, in *Input) {
// 	// although things like "text input" and "number inputs" are "fields"
// 	// interpolate arranges each into its own dummy input.
// 	// ( otherwise this will get super complicated.)
// 	var decoration *FieldLabel
// 	for fi, fcnt := 0, fields.Length(); fi < fcnt; fi++ {
// 		field := fields.Field(fi)
// 		if fieldName := field.Name(); fieldName == FieldDecor {
// 			decoration = &FieldLabel{field}
// 		}
// 	}
// 	parts := strings.Split(in.Name.String(), "/")
// 	bareName := parts[len(parts)-1]

// 	if fn, ok := reg.Find(container, named.Input(bareName)); ok {
// 	}
// }

func (b *Block) redecorate(reg decor.Decorators) {
	// 	blockCtx := &blockContext{block: b}
	// 	//
	// 	for inputIndex, cnt := 0, b.NumInputs(); inputIndex < cnt; inputIndex++ {
	// 		in := b.Input(inputIndex)
	// 		if decorate, ok := reg.Find(b.Type, in.Name); ok {
	// 			var field *FieldLabel
	// 			fields := in.Fields()
	// 			// should probably just be one field, but inputIndex guess you never know
	// 			for inputIndex, cnt := 0, fields.Length(); inputIndex < cnt; inputIndex++ {
	// 				if f := fields.Field(inputIndex); f.Name() == FieldDecor {
	// 					field = &FieldLabel{f}
	// 					break
	// 				}
	// 			}

	// 			inputCtx := &inputContext{parentContext{blockCtx}, in}
	// 			newText := decorate(inputCtx)
	// 			if field != nil {
	// 				field.SetText(newText)
	// 			} else if len(text) > 0 {
	// 				cssClass := "decor"
	// 				in.AppendNamedField(FieldDecor, NewFieldLabel(newText, cssClass))
	// 			}

	// 			if m := in.Mutation(); m != nil {
	// 				mutationCtx := mutationContext{parentContext{inputCtx}, m}

	// 				for i := 0; i < m.LeadingInputs; i++ {
	// 					inputIndex++
	// 					in := b.Input(inputIndex)
	// 					leadingCtx := &inputContext{parentContext{mutationCtx}, in}
	// 				}

	// 				atoms := m.NumAtoms()
	// 				for a, atoms := 0; a < atoms; a++ {
	// 					atom := m.Atom(a)
	// 					atomCtx := &atomContext{m, a}
	// 					atomInputs := atom.NumInputs()
	// 					for atomInput := 0; atomInput < atomInputs; atomInput++ {
	// 						in := b.Input(i + atomInput + 1)
	// 						// each atom has its own type
	// 						if decorate, ok := reg.Find(atom.Type, in.Name); ok {
	// 							// iiiiiiii + atomInputs
	// 							// ctx:= &inputContext{parentContext{atomCtx},
	// 						}

	// 					}

	// 				}
	// 			}

	// 			for i := 0; i < m.TrailingInputs; i++ {
	// 				inputIndex++
	// 				in := b.Input(inputIndex + i)
	// 				trailingCtx := &inputContext{parentContext{mutationCtx}, in}
	// 			}
	// 		}
}

// 	// we'll need some sort of trailing input and field for this.
// 	// if nc := b.NextConnection(); nc != nil {
// 	// 	if decorate, ok := reg.Find(b.Type, NextInput); ok {
// 	//    ctx:= &outputContext{ nc}
// 	// 		decorate(ctx)
// 	// 	}
// 	// }
// }
