package tin

import (
	r "reflect"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
)

// iteration over the extra fields of a mutation.
type fixedIt struct {
	mutable *Mutable
}

// return the unscoped name
func (it *fixedIt) Name() string {
	return block.Scope(it.mutable.name, "")
}

// return the unscoped name
func (it *fixedIt) PtrType() r.Type {
	return it.mutable.ptrType
}

// return the mui block type, scoped to its mutation.
// ex. "mui$mutation$atomw"
func (it *fixedIt) BlockType() string {
	return block.Scope("mui", it.mutable.name, "")
}

// return the displayed label of the indexed mui block
func (it *fixedIt) Label() string {
	// strip mutation off of name
	name, suffix := it.mutable.name, "_mutation"
	if strings.HasSuffix(name, suffix) {
		name = name[:len(name)-len(suffix)]
	}
	// replace underscores with spaces
	return strings.Join(strings.Split(name, "_"), " ")
}

func (it *fixedIt) LimitsOfNext() (ret block.Limits) {
	return it.mutable.limitsOfNext(it.PtrType())
}

// scope provides a unique name for this atom's inputs
func (it *fixedIt) Atomize(scope string, db mutant.Atomizer) (ret block.Args, err error) {
	c := context{Atomizer: db}
	if args, e := c.buildItems(scope, it.mutable.ptrType, nil); e != nil {
		err = errutil.Fmt("%s while atomizing the fixed fields of %q", e, it.mutable.name)
	} else {
		ret = args
	}
	return
}

// implements quarkIterator
func (it *fixedIt) NextQuark() (mutant.Quark, bool) {
	return it.mutable.firstQuark()
}

// implements typeIterator
func (it *fixedIt) NextType() (typeIterator, bool) {
	return it.mutable.firstQuark()
}
