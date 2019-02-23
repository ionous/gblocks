package mutation

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/blockly"
	"github.com/ionous/gblocks/dom"
	"github.com/ionous/gblocks/inspect"
	r "reflect"
)

type WorkspaceFactory struct {
	*blockly.Workspace
	mixin    blockly.Mixin             // functions used for generating the mutation
	mutators map[block.Type]block.Type // block types that contain mutations; multiple blocks can share the same mutation?
	palette  map[block.Type]*Palette   // mutations to palette data

	atoms        map[block.Type]r.Type // sub-blocks used by mutations; FIX: we dont really need the items.
	enums        inspect.EnumPairs
	deps         inspect.DependencyPool // all types that can connect with atom inputs
	blockDataMap blockDataMap           // custom per-input, per-block data about the mutation(s)
}

// FIX: and enus, etc.
func NewWorkspaceFactory(ws *blockly.Workspace) *WorkspaceFactory {
	mw := &WorkspaceFactory{Workspace: ws, blockDataMap: make(blockDataMap)}
	mw.AddChangeListener(mw.mutate)
	return mw
}

func (mw *WorkspaceFactory) RegisterAtom(rtype r.Type) (err error) {
	if _, e := inspect.MakeItems(rtype); e != nil {
		err = e
	} else {
		name := block.TypeFromStruct(rtype)
		if mw.atoms == nil {
			mw.atoms = map[block.Type]r.Type{name: rtype}
		} else {
			mw.atoms[name] = rtype
		}
	}
	return
}

func (mw *WorkspaceFactory) RegisterMutation(rtype r.Type) (err error) {
	// name := block.TypeFromStruct(rtype.Elem())
	// var items inspect.Items
	// inspect.VisitItems(rtype, func(it *inspect.Item, e error) (keepGoing bool) {
	// 	if e != nil {
	// 		err = errutil.Append(err, e)
	// 	} else if it.Class == inspect.Group {
	// 		// maybe --
	// 		// where do we keep the *group* types
	// 		// where do we build / keep the quarks
	// 	}
	// 	return true
	// })
	return
}

// note: all of the events we get are only for this workspace
// there are no global events, and none of the events carry information about other workspaces
// basically: using events its impossible to learn about the mutation workspace(s).
func (mw *WorkspaceFactory) mutate(evt interface{}) {
	switch evt := evt.(type) {
	case *blockly.BlockCreate:
		// ids is an array of js strings
		for i := 0; i < evt.Ids.Length(); i++ {
			id := block.Id(evt.Ids.Index(i).String())
			if b := mw.GetBlockById(id); b != nil {
				if mutator, ok := mw.mutators[b.Type]; ok {
					blockly.Extensions.Apply(mutator.String(), b, true)
					mw.blockDataMap[id] = &blockData{make(inputMap), make(connectionMap)}
				}
			}
		}
	case *blockly.BlockDelete:
		// ids is an array of js strings
		for i := 0; i < evt.Ids.Length(); i++ {
			id := block.Id(evt.Ids.Index(i).String())
			delete(mw.blockDataMap, id)
		}
	}
}

func (mw *WorkspaceFactory) getMixin() blockly.Mixin {
	if len(mw.mixin) == 0 {
		mw.mixin = blockly.Mixin{
			"mutationToDom": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
				if mb, e := mw.mutableBlock(obj); e != nil {
					panic(e)
				} else if xml, e := mb.mutationToDom(); e != nil {
					panic(e)
				} else {
					ret = xml.Object
				}
				return
			}),
			"domToMutation": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
				// workspace block, source xml
				if mb, e := mw.mutableBlock(obj); e != nil {
					panic(e)
				} else if _, e := mb.domToMutation(&dom.Element{Object: parms[0]}); e != nil {
					panic(e)
				}
				return
			}),
			"decompose": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
				// workspace block, mutation workspaace -> mutation container
				if mb, e := mw.mutableBlock(obj); e != nil {
					panic(e)
				} else if muiContainer, e := mb.decompose(&blockly.Workspace{Object: parms[0]}); e != nil {
					panic(e)
				} else {
					ret = muiContainer.Object
				}
				return
			}),
			"compose": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
				if mb, e := mw.mutableBlock(obj); e != nil {
					panic(e)
				} else if e := mb.compose(&blockly.Block{Object: parms[0]}); e != nil {
					panic(e)
				}
				return
			}),
			"saveConnections": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
				// workspace block, mutation ui container
				if mb, e := mw.mutableBlock(obj); e != nil {
					panic(e)
				} else if e := mb.saveConnections(&blockly.Block{Object: parms[0]}); e != nil {
					panic(e)
				}
				return
			}),
		}
	}
	return mw.mixin
}

func (mw *WorkspaceFactory) mutableBlock(obj *js.Object) (ret *mutableBlock, err error) {
	b := &blockly.Block{Object: obj}
	if blockData, ok := mw.blockDataMap[b.Id]; !ok {
		err = errutil.New("block not mutable", b.Type)
	} else {
		ret = &mutableBlock{mw, b, blockData}
	}
	return
}
