package gblocks

import (
	"github.com/ionous/errutil"
)

func (b *Block) updateShape(ws *Workspace) (err error) {
	// collapse all dynamic inputs
	var mutations []int // track which are mutations
	for i, collapse := 0, 0; i < b.NumInputs(); {
		in := b.Input(i)
		if collapse > 0 {
			b.RemoveInput(in.Name)
			collapse--
		} else {
			if in.mutations != 0 {
				mutations = append(mutations, i)
				collapse = in.mutations
				in.mutations = -1
			}
			i++ // note: we dont update the index for removed elements
		}
	}

	// rebuild the block
	data := ws.BlockData(b)
	// for each mutable input
	offset := 0 // adjust the input position to account for input array growth
	for _, index := range mutations {
		in := b.Input(index + offset)
		if els, ok := data.Elements(in); !ok {
			err = errutil.Append(err, errutil.New("bad input", in.Name))
		} else if n := els.Len(); n > 0 {
			oldLen := b.NumInputs()
			// expand each element of the mutation into its own set of inputs.
			for i := 0; i < n; i++ {
				iface := els.Index(i)
				ptr := iface.Elem()
				el := ptr.Elem()
				t := el.Type()
				if msg, args, _, e := ws.reg.makeArgs(t); e != nil {
					err = e
					break
				} else {
					b.interpolate(msg, args)
				}
			}
			// swap appended inputs into their correct spot
			newLen := b.NumInputs()
			if addedElements := newLen - oldLen; addedElements > 0 {
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
				offset += addedElements
				in.mutations = addedElements
			}
		}
	}
	return
}
