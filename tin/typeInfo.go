package tin

import (
	r "reflect"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
)

// wrap a golang type with helpers for adapting to a blockly block.
type TypeInfo struct {
	Name    string
	Model   Model
	ptrType r.Type
}

// string representation of this type: "ex. block_name (TopBlock)"
func (t *TypeInfo) String() string {
	return t.Name + " (" + t.Model.String() + ")"
}

// expand the contents of this type into fields and inputs usable by blockly
func (t *TypeInfo) BuildItems(scope string, db mutant.Atomizer, mutables Mutables, out *mutant.InMutations) (block.Args, error) {
	c := context{db, mutables}
	return c.buildItems(scope, t.ptrType, out)
}

// return the subset of types which can be assigned this type.
func (t *TypeInfo) LimitsOfNext(tins []*TypeInfo) block.Limits {
	return limitsOfNext(t.ptrType, &typeFilter{tins, MidBlock})
}

// return the subset of types which have an ouput that can be assigned this type.
func (t *TypeInfo) LimitsOfOutput(tins []*TypeInfo) (block.Limits, error) {
	return limitsOfOutput(t.ptrType, &typeFilter{tins, TermBlock})
}
