package gblocks

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/blockly"
	"github.com/ionous/gblocks/inspect"
	"github.com/ionous/gblocks/option"
	r "reflect"
)

type Registry struct {
	inspect.EnumPairs
	inspect.DependencyPool
	//mutations Palettes
	//decor decor.Decorators
}

func (reg *Registry) RegisterBlock(b interface{}, blockDesc block.Dict) error {
	ptrType := r.TypeOf(b)
	name := block.TypeFromStruct(ptrType.Elem())
	return reg.registerType(name, ptrType, blockDesc)
}

func (reg *Registry) RegisterBlocks(blockDesc map[string]block.Dict, blocks ...interface{}) (err error) {
	for _, b := range blocks {
		ptrType := r.TypeOf(b)
		name := block.TypeFromStruct(ptrType.Elem())
		var sub block.Dict
		if blockDesc, ok := blockDesc[name.String()]; ok {
			sub = blockDesc
		}
		if e := reg.registerType(name, ptrType, sub); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// func (reg *Registry) RegisterMutation(m interface{}, muiBlocks ...Mutation) (err error) {
// 	structType := r.TypeOf(m).Elem()
// 	return reg.mutations.RegisterMutation(structType, muiBlocks...)
// }

func (reg *Registry) registerType(name block.Type, ptrType r.Type, blockDesc block.Dict) (err error) {
	if _, exists := reg.types[name]; exists {
		err = errutil.New("already registered", name)
	} else {
		// add the type, even if there's an error: we block it out from future attempts
		if reg.types == nil {
			reg.types = map[block.Type]r.Type{name: ptrType}
		} else {
			reg.types[name] = ptrType
		}
		//
		if blockDesc == nil {
			blockDesc = make(block.Dict)
		}
		if e := reg.buildBlockDesc(name, ptrType, blockDesc); e != nil {
			err = e
		} else {
			blockly.DefineBlock(name, blockDesc)
		}
	}
	return
}

var linkOptions = map[inspect.Class]string{inspect.NextLink: option.Next, inspect.PrevLink: option.Prev}

func (reg *Registry) buildBlockDesc(name block.Type, ptrType r.Type, blockDesc block.Dict) (err error) {
	args := inspect.NewArgs("", reg.EnumPairs, reg.DependencyPool)
	//
	inspect.VisitItems(ptrType.Elem(), func(item *inspect.Item, e error) bool {
		if e != nil {
			err = errutil.Append(err, e)
		} else if e := args.AddItem(item); e != nil {
			err = errutil.Append(err, e)
		} else if opt, ok := linkOptions[item.Class]; ok {
			if types, ok := check.GetConstraints(item.Type); !ok {
				err = errutil.Append(err, errutil.New("block", name, "has link with no matching types", item))
			} else {
				blockDesc.Insert(opt, types)
			}
		}
		return true
	})
	if err == nil {
		blockDesc.Insert(option.Type, name)

		// add arg and message content to the block
		var hasContent bool
		if msg := args.Message(); len(msg) > 0 {
			blockDesc.Insert(option.Message(0), msg)
			hasContent = true
		}
		if list := args.List(); len(list) > 0 {
			blockDesc.Insert(option.Args(0), list)
			hasContent = true
		}
		if !hasContent {
			blockDesc[option.Message(0)] = name.Friendly()
		}

		// block can have either prev statement or output()
		if _, hasPrev := blockDesc[option.Prev]; !hasPrev {
			if m, ok := ptrType.MethodByName(block.OutputMethod); ok {
				if cnt := m.Type.NumOut(); cnt != 1 {
					err = errutil.Append(err, errutil.New("unexpected output count", cnt))
				} else {
					outType := m.Type.Out(0)
					switch k := outType.Kind(); k {
					default:
						err = errutil.Append(err, errutil.New("unexpected output type", k))
					case r.Interface:
						if basicInterface := r.TypeOf((interface{})(nil)); outType != basicInterface {
							// basic interface is nil type.
							blockDesc[option.Output] = nil
						} else {
							blockDesc[option.Output] = block.TypeFromStruct(outType)
						}
					case r.Ptr:
						blockDesc[option.Output] = block.TypeFromStruct(outType.Elem())
					}
				}
			}
		}
	}
	return
}
