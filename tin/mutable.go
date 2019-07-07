package tin

import (
	r "reflect"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
)

// records the original golang representation of a mutation description.
// implements InMutation
type Mutable struct {
	name    string // name of the mutation type
	ptrType r.Type // original description of the mutation
	quarks  []r.Type
}

func (m *Mutable) Name() string {
	return m.name
}

// quarks which can appear at the top of an mui input stack.
func (m *Mutable) Limits() (ret block.Limits) {
	// if the mutation has "extra content" -- members other than just the next statement
	// limit the first block to the just that one block.
	// otherwise, return the quarks that can attach to the ptr's next statement.
	if first, ok := m.firstBlock(); ok {
		ret = block.MakeLimits([]string{first.Name()})
	} else {
		ret = m.limitsOfNext(m.ptrType)
	}
	return
}

// what mui blocks in this mutable can attach to the next statement of the passed type.
func (m *Mutable) limitsOfNext(ptrType r.Type) (ret block.Limits) {
	if types, ok := m.types(); ok {
		ret = limitsOfNext(ptrType, types)
	} else {
		ret = block.MakeOffLimits() // when no types are compatible
	}
	return
}

//
func (m *Mutable) FirstBlock() (ret mutant.Quark, okay bool) {
	return m.firstBlock()
}

// list of elements for a mutation
func (m *Mutable) Quarks(paletteOnly bool) (ret mutant.Quark, okay bool) {
	if !paletteOnly {
		ret, okay = m.firstBlock()
	}
	if !okay {
		ret, okay = m.firstQuark()
	}
	return
}

//
func (m *Mutable) types() (ret typeIterator, okay bool) {
	if x, ok := m.firstBlock(); ok {
		ret, okay = x, true
	} else {
		ret, okay = m.firstQuark()
	}
	return
}

// return an iterator over the "extra" members of a mutation.
func (m *Mutable) firstBlock() (ret *fixedIt, okay bool) {
	if HasContent(m.ptrType.Elem()) {
		ret, okay = &fixedIt{mutable: m}, true
	}
	return
}

// does the passed struct have members other than "NextStatement"?
func HasContent(elem r.Type) bool {
	var extraFields int
	if _, ok := elem.FieldByName(block.NextStatement); ok {
		extraFields = 1
	}
	return elem.NumField() > extraFields
}

//
func (m *Mutable) firstQuark() (ret *quarkIt, okay bool) {
	if len(m.quarks) > 0 {
		ret, okay = &quarkIt{mutable: m}, true
	}
	return
}
