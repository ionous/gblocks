package mutant

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

// add inputs to a block immediately after an existing input
// ( unfortunately, they wind up at the end of the block )
type injector struct {
	b           block.Shape
	in          block.Input
	dst, oldLen int
}

func newInjector(b block.Shape, input string) (ret *injector, err error) {
	if in, inputIndex := b.InputByName(input); inputIndex < 0 {
		err = errutil.New("unknown input", input)
	} else {
		ret = &injector{b, in, inputIndex + 1, b.NumInputs()}
	}
	return
}

func (j *injector) target() block.Input {
	return j.in
}

func (j *injector) inject(args block.Args) (err error) {
	j.b.Interpolate(args.Message(), args.List())
	return nil
}

// reorder inputs so that the atom's inputs follow the mutation's inputs.
func (j *injector) finalizeInputs() {
	b := j.b
	dst, newLen, oldLen := j.dst, b.NumInputs(), j.oldLen
	if numInputs := newLen - oldLen; numInputs > 0 {
		// if there were no elements beyond the original destination
		// then the added elements are already in the right place: at the end of the list.
		if oldEls := oldLen - dst; oldEls > 0 {
			// record the desired order
			inputOrder := make([]block.Input, 0, numInputs+oldEls)
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
	}
}
