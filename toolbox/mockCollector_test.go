package toolbox_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/tin"
	"github.com/ionous/gblocks/toolbox"
)

// provide a sink ( similar to gblocks.Maker ) to collect blocks visited during toolbox creation
// note, this needs tin and toolbox does not -- so could go in gblocks maybe
type testCollector struct {
	Types     []*tin.TypeInfo
	LastError error
}

// return a toolbox builder which registers blocks to this maker
// same as: Maker.NewBlocks(toolbox.SubShadow)
func (m *testCollector) NewBlocks(shadowing toolbox.Shadowing) *toolbox.Builder {
	return toolbox.NewBlocks(shadowing, m, &atomTestNames{})
}

func (m *testCollector) OnBlock(t *tin.TypeInfo) {
	if nm, found := tin.FindByName(m.Types, t.Name); found && nm.Model != t.Model {
		m.OnError(errutil.New("type already registered", t))
	} else if !found {
		m.Types = append(m.Types, t)
	}
}

func (m *testCollector) OnError(err error) {
	m.LastError = errutil.Append(m.LastError, err)
}

func (m *testCollector) Collected() (ret []string, err error) {
	if m.LastError != nil {
		err = m.LastError
	} else {
		for _, t := range m.Types {
			ret = append(ret, t.String())
		}
	}
	return
}
