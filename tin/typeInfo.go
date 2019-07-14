package tin

import (
	r "reflect"

	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/mutant"
)

// wrap a golang type with helpers for adapting to a blockly block.
type TypeInfo struct {
	Name    string // under_scores
	Model   Model
	ptrType r.Type
}

// string representation of this type: "ex. block_name (TopBlock)"
func (t *TypeInfo) String() string {
	return t.Name + " (" + t.Model.String() + ")"
}

// expand the contents of this type into fields and inputs usable by blockly
func (t *TypeInfo) BuildItems(scope string, db mutant.Atomizer, mutables Mutations, out *mutant.BlockMutations) (block.Args, error) {
	c := context{db, mutables}
	return c.buildItems(scope, t.ptrType, out)
}

// return the subset of types which can be assigned this type.
func (t *TypeInfo) LimitsOfNext(tins []*TypeInfo) (ret block.Limits) {
	if len(tins) > 0 {
		ret = limitsOfNext(t.ptrType, &typeFilter{tins, MidBlock})
	} else {
		ret = block.MakeOffLimits()
	}
	return
}

// return the subset of types which have an ouput that can be assigned this type.
func (t *TypeInfo) LimitsOfOutput(tins []*TypeInfo) (ret block.Limits, err error) {
	if len(tins) > 0 {
		ret, err = limitsOfOutput(t.ptrType, &typeFilter{tins, TermBlock})
	} else {
		ret = block.MakeOffLimits()
	}
	return
}

// visit all block.Option tags
func (t *TypeInfo) VisitOptions(cb func(k, v string)) {
	structType := t.ptrType.Elem()
	for i, cnt := 0, structType.NumField(); i < cnt; i++ {
		if field := structType.Field(i); len(field.PkgPath) == 0 {
			if field.Name != block.NextStatement {
				if Classify(field.Type) == Option {
					visitTags(string(field.Tag), cb)
				}
			}
		}
	}
}
