package gblocks

import (
	"github.com/ionous/errutil"
)

func (b *Block) updateShape(ws *Workspace) (err error) {
	// collapse all dynamic inputs
	var indexOfMutations []int // track which are indexOfMutations
	for i, collapse := 0, 0; i < b.NumInputs(); {
		in := b.Input(i)
		if collapse > 0 {
			b.RemoveInput(in.Name)
			collapse--
		} else {
			if m := in.Mutation(); m != nil {
				indexOfMutations = append(indexOfMutations, i)
				collapse = m.Reset() // collapse equals the total number of inputs used by this mutation
			}
			i++
		}
	}
	// rebuild the block
	ctx := ws.Context(b.Id)
	// for each mutable input
	offset := 0 // adjust the input position to account for input array growth
	for _, index := range indexOfMutations {
		in := b.Input(index + offset)
		els := ctx.Elem().FieldByName(in.Name.FieldName())
		if n := els.Len(); n > 0 {
			oldLen := b.NumInputs()

			// append inputs for each each element of the mutation
			for i, prevLen := 0, oldLen; i < n; i++ {
				iface := els.Index(i)
				ptr := iface.Elem()
				el := ptr.Elem()
				t := el.Type()
				if msg, args, _, e := TheRegistry.makeArgs(t); e != nil {
					err = e
					break
				} else if m := in.Mutation(); m == nil {
					err = errutil.New("missing mutation", in.Name)
					break
				} else {
					b.interpolate(msg, args)
					newLen := b.NumInputs()
					m.AddAtom(newLen - prevLen)
					prevLen = newLen
				}
			}

			// swap appended inputs into their correct spot
			newLen := b.NumInputs()
			if addedInputs := newLen - oldLen; addedInputs > 0 {
				scratch := make([]*Input, 0, newLen)
				// mutation element, all new elements go directly after this
				el := index + offset
				for _, rng := range [][]int{{0, el + 1}, {oldLen, newLen}, {el + 1, oldLen}} {
					for i, last := rng[0], rng[1]; i < last; i++ {
						in := b.Input(i)
						scratch = append(scratch, in)
					}
				}
				for i, in := range scratch {
					b.setInput(i, in)
				}
				offset += addedInputs
			}
		}
	}
	return
}
