package gblocks

import (
	"github.com/ionous/errutil"
	"strconv"
	"strings"
)

func (b *Block) updateShape(ws *Workspace) (err error) {
	if ws == nil {
		err = errutil.New("update shape into nil workspace")
	} else if ctx := ws.Context(b.Id); ctx == nil {
		err = errutil.New("update shape into invalid workspace")
	} else {
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

			// append inputs for each each element of the mutation
		// for each mutable input
		// input name format: INPUT_NAME/AtomIndex/INPUT_NAME
		offset := 0 // adjust the input position to account for input array growth
		for _, index := range indexOfMutations {
			in := b.Input(index + offset)
			els := ctx.FieldForInput(in.Name)
			if n := els.Len(); n > 0 {
				oldLen := b.NumInputs()

				// append atoms (sets of inputs) for each element of the mutation
				// they wind up at the end of all inputs, and need to be swapped.
				for i, prevLen := 0, oldLen; i < n; i++ {
					iface := els.Index(i)
					ptr := iface.Elem()
					el := ptr.Elem()
					t := el.Type()
					path := strings.Join([]string{in.Name.String(), strconv.Itoa(i), ""}, "/")
					if msg, args, _, e := TheRegistry.makeArgs(t, path); e != nil {
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

				// swap appended atoms into their correct spot
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
	}
	return
}
