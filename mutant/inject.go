package mutant

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

// add inputs to a block immediately after an existing input
// ( unfortunately, they wind up at the end of the block by default )
type injector struct {
	target      *MutatedBlock
	min         *MutatedInput
	dst, oldLen int
}

func newInjector(mb *MutatedBlock, min *MutatedInput) (ret *injector, err error) {
	if _, inputIndex := mb.block.InputByName(min.InputName); inputIndex < 0 {
		err = errutil.New("unknown input", min.InputName)
	} else {
		ret = &injector{mb, min, inputIndex + 1, mb.block.NumInputs()}
	}
	return
}

func (j *injector) inject(args block.Args) int {
	b := j.target.block
	was := b.NumInputs()
	b.Interpolate(args.Message(), args.List())
	now := b.NumInputs()
	return now - was
}

// reorder inputs so that the atom's inputs follow the mutation's inputs.
// returns the number of added inputs
func (j *injector) finalizeInputs() (ret int) {
	b := j.target.block
	dst, newLen, oldLen := j.dst, b.NumInputs(), j.oldLen
	if newInputs := newLen - oldLen; newInputs > 0 {
		// if there were no elements beyond the original destination
		// then the added elements are already in the right place: at the end of the list.
		if oldEls := oldLen - dst; oldEls > 0 {
			// record the desired order
			inputOrder := make([]block.Input, 0, newInputs+oldEls)
			// there are three sections; the first of which we can leave as is.
			// 1. up-to-and-including the mutation input [:dst).
			// 2. all the appended inputs [oldLen:newLen).
			// 3. the inputs following the mutation input [dst:oldLen).
			for _, rng := range [][]int{{oldLen, newLen}, {dst, oldLen}} {
				start, end := rng[0], rng[1]
				for i := start; i < end; i++ {
					inputOrder = append(inputOrder, b.Input(i))
				}
			}
			// enforce the desired order
			for i, in := range inputOrder {
				b.SetInput(dst+i, in)
			}
		}
		ret = newInputs
	}
	return
}
