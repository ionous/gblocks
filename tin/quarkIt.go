package tin

import (
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
	"github.com/ionous/gblocks/pascal"
)

// walks a quark, returning each quark in turn.
type quarkIt struct {
	mutable *Mutable
	i       int
}

// return the unscoped name
func (it *quarkIt) Name() string {
	return pascal.ToUnderscore(it.PtrType().Elem().Name())
}

// implements typeIterator
func (it *quarkIt) PtrType() r.Type {
	return it.mutable.quarks[it.i]
}

// return the mui block type, scoped to its mutation.
// ex. "mui$block_mutation$test_atom"
func (it *quarkIt) BlockType() string {
	return block.Scope("mui", it.mutable.name, it.Name())
}

// return the displayed label of the indexed mui block
func (it *quarkIt) Label() string {
	return pascal.ToSpaces(it.PtrType().Elem().Name())
}

//
func (it *quarkIt) LimitsOfNext() block.Limits {
	return it.mutable.limitsOfNext(it.PtrType())
}

// scope provides a unique name for this atom's inputs
func (it *quarkIt) Atomize(scope string, db mutant.Atomizer) (ret block.Args, err error) {
	c := context{Atomizer: db}
	if args, e := c.buildItems(scope, it.PtrType(), nil); e != nil {
		err = errutil.New(e, "while atomizing", it.Name())
	} else {
		ret = args
	}
	return
}

// implements quarkIterator
func (it *quarkIt) NextQuark() (mutant.Quark, bool) {
	return it.nextQuark()
}

// implements typeIterator
func (it *quarkIt) NextType() (typeIterator, bool) {
	return it.nextQuark()
}

//
func (it *quarkIt) nextQuark() (ret *quarkIt, okay bool) {
	if n := it.i + 1; n < len(it.mutable.quarks) {
		ret, okay = &quarkIt{it.mutable, n}, true
	}
	return
}
