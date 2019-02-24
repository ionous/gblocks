package inspect

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	"github.com/ionous/gblocks/option"
	r "reflect"
)

type TypePool struct {
	EnumPairs
	DependencyPool
}

// note: to handle dependencies properly, all types need to be registered before block descriptions are extracted.
func (tp *TypePool) BuildDesc(ptrType r.Type, blockDesc block.Dict) (ret block.Dict, err error) {
	typeName := block.TypeFromStruct(ptrType.Elem())
	if match, ok := tp.Types[typeName]; !ok || match != ptrType {
		err = errutil.New("type", typeName, "not registered", ptrType)
	} else {
		if blockDesc == nil {
			blockDesc = make(block.Dict)
		}
		if e := tp.buildDesc(typeName, ptrType, blockDesc); e != nil {
			err = e
		} else {
			ret = blockDesc
		}
	}
	return
}

func (tp *TypePool) buildDesc(typeName block.Type, ptrType r.Type, blockDesc block.Dict) (err error) {
	args := NewArgs(blockDesc, tp.EnumPairs, tp.DependencyPool)
	if e := args.AddMembers("", ptrType.Elem()); e != nil {
		err = e
	} else {
		blockDesc.Insert(option.Type, typeName)

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
			blockDesc[option.Message(0)] = typeName.Friendly()
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
						if basicInterface := r.TypeOf((*interface{})(nil)).Elem(); outType != basicInterface {
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
