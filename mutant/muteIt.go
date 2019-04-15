package mutant

// walk over all quarks in a block's worth of input mutations
type muteIt struct {
	*InMutations
	input int
	Quark
	paletteOnly bool
}

func (it *muteIt) NextQuark() (ret Quark, okay bool) {
	if sub, ok := it.Quark.NextQuark(); ok {
		ret, okay = it.clone(it.input, sub), true
	} else if q, ok := it.advance(); ok {
		ret, okay = q, true
	}
	return
}

// create a new iterator
func (it *muteIt) clone(next int, sub Quark) *muteIt {
	return &muteIt{it.InMutations, next, sub, it.paletteOnly}
}

func (it *muteIt) advance() (ret *muteIt, okay bool) {
	for next := it.input + 1; next < len(it.Inputs); next++ {
		inputName := it.Inputs[next]
		if m, ok := it.GetMutation(inputName); !ok {
			panic("unstable mutator")
		} else if q, ok := m.Quarks(it.paletteOnly); ok {
			ret = it.clone(next, q)
			okay = true
			break
		}
	}
	return
}
