package gblocks

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/tin"
	"github.com/ionous/gblocks/toolbox"
)

// provide a sink ( similar to gblocks.TypeCollector ) to collect blocks visited during toolbox creation
type TypeCollector struct {
	types     []*tin.TypeInfo
	lastError error
	names     UniqueNames
}

type UniqueNames interface{ GenerateUniqueName() string }

func NewTypeCollector(names UniqueNames) *TypeCollector {
	return &TypeCollector{names: names}
}

// golang pimpl
type toolboxHelper struct {
	collector *TypeCollector
}

func (cb *toolboxHelper) GenerateUniqueName() string {
	return cb.collector.names.GenerateUniqueName()
}

func (cb *toolboxHelper) OnBlock(t *tin.TypeInfo) {
	cb.collector.addOnce(t)
}

func (cb *toolboxHelper) OnError(err error) {
	cb.collector.lastError = errutil.Append(cb.collector.lastError, err)
}

// return a toolbox builder which registers blocks to this maker
func (tc *TypeCollector) NewBlocks() *toolbox.Builder {
	// FIX: nested shadowing doesnt provide a nice user experience
	// need some sort of listener convert a whole placed shadow chain with an unshadowed equivalent
	// whenever a user ties to edit or place something into one of the nested element
	shadowing := toolbox.NoShadow //toolbox.SubShadow
	return tc.NewShadows(shadowing)
}

// return a toolbox builder which registers blocks to this maker
// same as: TypeCollector.NewBlocks(toolbox.SubShadow)
func (tc *TypeCollector) NewShadows(s toolbox.Shadowing) *toolbox.Builder {
	return toolbox.NewBlocks(s, &toolboxHelper{tc}, tc.names)
}

// you're either a term, which can contain input statements;
// or you're statement, which can contain input terms.
// fix: used for testing only; public b/c system tests are in a different package
func (tc *TypeCollector) AddTerm(ptr interface{}) *TypeCollector {
	tc.addType(ptr, tin.TermBlock)
	return tc
}

// you're either a statement, which can contain input terms;
// or you're a term, which can contain input statements.
// fix: used for testing only; public b/c system tests are in a different package
func (tc *TypeCollector) AddStatement(ptr interface{}) *TypeCollector {
	tc.addType(ptr, tin.MidBlock)
	return tc
}

// a top statement is a statement with no previous connection.
// fix: used for testing only; public b/c system tests are in a different package
func (tc *TypeCollector) AddTopStatement(ptr interface{}) *TypeCollector {
	tc.addType(ptr, tin.TopBlock)
	return tc
}

func (tc *TypeCollector) GetTypes() (ret []*tin.TypeInfo, err error) {
	if tc.lastError != nil {
		err = tc.lastError
	} else {
		ret = tc.types
	}
	return
}

func (tc *TypeCollector) addType(ptr interface{}, model tin.Model) (okay bool) {
	var err error
	if t, e := model.PtrInfo(ptr); e != nil {
		err = e
	} else if _, found := tin.FindByName(tc.types, t.Name); found {
		err = errutil.New("type already registered", t.Name)
	} else {
		tc.types = append(tc.types, t)
	}
	if err != nil {
		tc.lastError = errutil.Append(tc.lastError, err)
	} else {
		okay = true
	}
	return
}

// helper for toolbuilder; where its okay to add more than once
func (tc *TypeCollector) addOnce(t *tin.TypeInfo) (okay bool) {
	var err error
	if nm, found := tin.FindByName(tc.types, t.Name); found && nm.Model != t.Model {
		err = errutil.New("type mismatch", t.Name, "was", nm.Model, "now", t.Model)
	} else if !found {
		tc.types = append(tc.types, t)
	}
	if err != nil {
		tc.lastError = errutil.Append(tc.lastError, err)
	} else {
		okay = true
	}
	return
}
