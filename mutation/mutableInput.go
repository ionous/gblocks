package mutation

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/blockly"
	"github.com/ionous/gblocks/inspect"
)

// associates mutable inputs with workspace and block pointers; alleviates passing a large numbers of parameters to helper functions
type mutableInput struct {
	*mutableBlock
	name  block.Item
	index int
	*inputData
}

func (mi *mutableInput) palette() (ret *Palette, okay bool) {
	if palette, ok := mi.workspace.palette[mi.mutableInput]; ok {
		ret = palette
	}
	return
}

// AddAtom - some number of contiguous inputs (already added to the parent block).
func (mi *mutableInput) addAtom(atomType block.Type) (ret int, err error) {
	if ptrType, ok := mi.workspace.atoms[atomType]; !ok {
		err = errutil.Append(err, errutil.New("unknown atom type", atomType))
	} else {
		atomIndex := len(mi.atoms)
		scope := mi.name.Index(atomIndex + 1)
		args := inspect.NewArgs(scope, mi.workspace.enums, mi.workspace.deps)

		inspect.VisitItems(ptrType.Elem(), func(item *inspect.Item, e error) bool {
			if e != nil {
				err = errutil.Append(err, e)
			} else if e := args.AddItem(item); e != nil {
				err = errutil.Append(err, e)
			}
			return true // keepGoing
		})
		if args.Len() > 0 {
			b := mi.block
			// generate new inputs from the atom
			oldLen := b.NumInputs()
			b.Interpolate(args.Message(), args.List())
			newLen := b.NumInputs()
			numInputs := newLen - oldLen
			// reorder inputs so that the atom's inputs follow the mutation's inputs.
			if numInputs > 0 {
				// record the desired order
				scratch := make([]*blockly.Input, 0, newLen)
				// there are three sections
				// 1. up-to-and-including the mutation input
				// 2. the atom's added inputs ( which were appened to the input list )
				// 3. the inputs originally following the mutation input
				end := mi.index + mi.totalInputs + 1
				for _, rng := range [][]int{{0, end}, {oldLen, newLen}, {end, oldLen}} {
					for i, last := rng[0], rng[1]; i < last; i++ {
						scratch = append(scratch, b.Input(i))
					}
				}
				// rewrite the input order
				for i, in := range scratch {
					b.SetInput(i, in)
				}
			}
			// record the atom
			mi.store(atomType, numInputs)
			ret = numInputs
		}
	}
	return
}
