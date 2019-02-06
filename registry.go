package gblocks

import (
	"fmt" // for enum key names
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	r "reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Registry struct {
	types     RegisteredTypes
	enums     RegisteredEnums
	mutations RegisteredMutations
}

var blockPtr = r.TypeOf((*Block)(nil))

// NewData - returns a pointer to passed type name
func (reg *Registry) NewData(name TypeName) (ret r.Value, err error) {
	if t, ok := reg.types[name]; !ok {
		err = errutil.New("NewData for unknown type '" + name + "'")
	} else {
		ret = r.New(t)
	}
	return
}

func (reg *Registry) RegisterEnum(n interface{}) error {
	_, e := reg.registerEnum(n)
	return e
}

func (reg *Registry) RegisterBlock(b interface{}, blockDesc Dict) error {
	structType := r.TypeOf(b).Elem()
	typeName := toTypeName(structType)
	return reg.registerType(typeName, structType, blockDesc)
}

func (reg *Registry) RegisterBlocks(blockDesc map[string]Dict, blocks ...interface{}) (err error) {
	for _, b := range blocks {
		structType := r.TypeOf(b).Elem()
		typeName := toTypeName(structType)
		var sub Dict
		if blockDesc, ok := blockDesc[typeName.String()]; ok {
			sub = blockDesc
		}
		if e := reg.registerType(typeName, structType, sub); e != nil {
			e := errutil.New("error registering", typeName, e)
			err = errutil.Append(err, e)
		}
	}
	return
}

func (reg *Registry) RegisterMutation(mutationName string, muiBlocks ...Mutation) (err error) {
	if blockly := GetBlockly(); blockly == nil {
		err = errutil.New("blockly doesnt exist")
	} else {
		// add the "tops" of the prototypes to the pool we pull from to connect "next" blocks.
		types := make(RegisteredTypes)
		for _, el := range muiBlocks {
			structType := r.TypeOf(el.Creates).Elem()
			types.RegisterType(structType)
		}

		// now, walk those prototypes again
		var quarks []TypeName
		blocks := make(map[TypeName]*MutationBlock)

		for _, muiBlock := range muiBlocks {
			prototype, label := muiBlock.Creates, muiBlock.Label
			structType := r.TypeOf(prototype).Elem()
			typeName := toTypeName(structType)

			// scan to the end of the prototype's NextStatement stack
			// lastVal := val
			// var lastField []int
			// for {
			// 	lastType := lastVal.Type()
			// 	if f, ok := lastType.FieldByName(NextField); !ok {
			// 		lastField = nil // there is no next field; clear anything from a previous block in the chain
			// 		break
			// 	} else if nextVal := val.FieldByIndex(f.Index); !nextVal.IsValid() || nextVal.IsNil() {
			// 		lastField = f.Index
			// 		break
			// 	} else {
			// 		lastVal = nextVal.Elem()
			// 		lastType = nextVal.Type()
			// 	}
			// }
			// var constraints Constraints
			// if len(lastField) > 0 {
			// 	if c, e := types.CheckStructField(lastVal.Type().FieldByIndex(lastField)); e != nil {
			// 		err = errutil.Append(err, e)
			// 	} else {
			// 		constraints = c
			// 	}
			// }

			// future: prototype into dom tree
			// xml := ValueToDom(v, true)
			// // does the element have sub-elements (or is it just one block?)
			// var subElements int
			// if shadows := xml.GetElementsByTagName("shadow"); shadows != nil {
			// 	subElements = shadows.Num()
			// }

			if constraints, e := types.CheckField(structType, NextField); e != nil {
				err = errutil.Append(err, e)
			} else {
				muiType := SpecialTypeName("mui", mutationName, typeName.String())
				muiBlock := &MutationBlock{
					MuiLabel:      label,
					MuiType:       muiType,
					WorkspaceType: typeName,
					Constraints:   constraints,
					//		BlockXml:      xml,
				}
				blockly.AddBlock(muiType, muiBlock.BlockFns())
				quarks = append(quarks, muiType)
				blocks[muiType] = muiBlock
			}

			//
			// if !dupes[typeName] {
			// 	// if there are sub-elements; then we have also register the first block
			// 	if isAtom := subElements == 0; !isAtom {
			// 		muiType := SpecialTypeName"mui", name, el.Name, "atom")
			// 		b := &MutationBlock{
			// 			MuiLabel: typeName.Friendly(),
			// 			WorkspaceType: typeName,
			// 			BlockXml:      xml,
			// 			// FIXXXX -- these constraints are probably wrong...
			// 			Constraints: constraints,
			// 		}
			// 		blockly.AddBlock(muiType, b.BlockFns())
			// 		atoms = append(atoms, muiType)

			// 	}
			// 	// regardless, we either have added the atom, or the block itself was an atom.
			// 	dupes[typeName] = true
			// }

		}
		// append the atoms at the end of the other blocks
		// quarkNames, atomNames = append(quarkNames, atomNames...), nil

		// TODO: color code the blocks by the mui's container input -- each input a different set of shades.
		if reg.mutations == nil {
			reg.mutations = make(RegisteredMutations)
		}
		reg.mutations[mutationName] = &RegisteredMutation{blocks: blocks, quarks: quarks}
	}
	return
}

func (reg *Registry) registerType(typeName TypeName, structType r.Type, blockDesc Dict) (err error) {
	if blockly := GetBlockly(); blockly == nil {
		err = errutil.New("blockly doesnt exist")
	} else {
		if reg.types == nil {
			reg.types = make(RegisteredTypes)
		}
		if !reg.types.RegisterType(structType) {
			err = errutil.New("type already exists " + typeName)
		} else {
			if blockDesc == nil {
				blockDesc = make(Dict)
			}
			if mui, e := reg.buildBlockDesc(structType, blockDesc); e != nil {
				err = e
			} else {
				// the ui system has "blocks" -- mapping type names to prototypes; standalone tests do not
				init := js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
					b := &Block{Object: obj}
					// create the blockly block

					if e := b.JsonInit(blockDesc); e != nil {
						panic(e)
					} else if len(mui) > 0 {
						// all of the mutation blocks used by all of the mutable inputs in this block
						var quarkNames []TypeName
						// walk all arguments to find mutations: ( FIX: slow, awkward )
						for _, mi := range mui {
							inputName, mutationName := mi.inputName, mi.mutationName
							if registeredMutation, ok := reg.mutations[mutationName]; !ok {
								panic("unknown mutation " + mutationName)
							} else if in, index := b.InputByName(inputName); index < 0 {
								panic("unknown input " + inputName)
							} else {
								in.ForceMutation(mutationName)
								// add to the palette of types shown by the mutator ui
								quarkNames = append(quarkNames, registeredMutation.quarks...)
							}
						}
						b.SetMutator(NewMutator(quarkNames))
					}
					return
				})
				var fns Dict
				if len(mui) == 0 {
					fns = Dict{
						"init": init,
					}
				} else {
					fns = Dict{
						"init": init,
						"mutationToDom": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
							b := &Block{Object: obj}
							dom := b.mutationToDom()
							return dom.Object
						}),
						"domToMutation": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
							b, xmlElement := &Block{Object: obj}, &XmlElement{Object: parms[0]}
							if _, e := b.domToMutation(reg, xmlElement); e != nil {
								panic(e)
							}
							return
						}),
						"decompose": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
							b, mui := &Block{Object: obj}, &Workspace{Object: parms[0]}
							if muiContainer, e := b.decompose(reg, mui); e != nil {
								panic(e)
							} else {
								ret = muiContainer.Object
							}
							return
						}),
						"compose": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
							b, containerBlock := &Block{Object: obj}, &Block{Object: parms[0]}
							if e := b.compose(reg, containerBlock); e != nil {
								panic(e)
							}
							return
						}),
						"saveConnections": js.MakeFunc(func(obj *js.Object, parms []*js.Object) (ret interface{}) {
							b, containerBlock := &Block{Object: obj}, &Block{Object: parms[0]}
							if e := b.saveConnections(containerBlock); e != nil {
								panic(e)
							}
							return
						}),
					}
				}

				// register things to blockly.
				reg.types[typeName] = structType
				blockly.AddBlock(typeName, fns)

				// create custom mutation container
				// the container has one input pin for every mutatable field in the workspace block
				if len(mui) > 0 {
					fns := Dict{
						"init": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
							muiContainer := &Block{Object: obj}
							for _, mi := range mui {
								inputName, _ := mi.inputName, mi.mutationName
								label := NewFieldLabel(inputName.Friendly(), "")
								in := muiContainer.AppendStatementInput(inputName)
								in.AppendField(label.Field)
								// FIX: check for mutable inputs is whatever matches the field of block
								// in.SetCheck(mutableInput.check)
							}
							return
						}),
					}
					name := SpecialTypeName("mui_container", typeName.String())
					blockly.AddBlock(name, fns)
				}
			}
		}
	}
	return
}

func (reg *Registry) registerEnum(n interface{}) (ret []EnumPair, err error) {
	var pairs []EnumPair
	if src, srcType := r.ValueOf(n), r.TypeOf(n); srcType.Kind() != r.Map {
		err = errutil.New("invalid enum mapping")
	} else if keyType, valueType := srcType.Key(), srcType.Elem(); valueType.Kind() != r.String {
		err = errutil.New("invalid enum value type")
	} else {
		// want to build an array of display to stringer
		// want to sort that array for display
		// want to store that at the v.Type for lookup.
		// eventually a map? probably of stringer to int ( for reverse conversion, setting in response to changes )
		keys := src.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].Int() < keys[j].Int()
		})

		for _, key := range keys {
			unique := fmt.Sprint(key)
			display := src.MapIndex(key).String()
			pair := EnumPair{display, unique}
			pairs = append(pairs, pair)
		}
		if reg.enums == nil {
			reg.enums = make(RegisteredEnums)
		}
		enumName := toTypeName(keyType)
		reg.enums[enumName] = &RegisteredEnum{pairs: pairs}
		ret = pairs
	}
	return
}

func (reg *Registry) buildCheck(t r.Type, field string, blockDesc Dict, key string) (okay bool, err error) {
	if blockDesc.Contains(key) {
		okay = true
	} else if check, e := reg.types.CheckField(t, field); e != nil {
		err = e
	} else if check, ok := check.GetConstraints(); ok {
		blockDesc[key] = check
		okay = true
	}
	return
}

// note: perhaps anonymous structs could be used to separate into args blocks
func (reg *Registry) buildBlockDesc(t r.Type, blockDesc Dict) (retMui []mutationInput, err error) {
	name := toTypeName(t)
	blockDesc.Insert(opt_type, name)
	//zw
	if msg, args, mui, e := reg.buildArgs(t, ""); e != nil {
		err = errutil.Append(err, e)
	} else {
		if len(msg) > 0 {
			blockDesc.Insert(opt_message0, msg)
		}
		if len(args) > 0 {
			blockDesc.Insert(opt_args0, args)
		}
		if !blockDesc.Contains(opt_args0) && !blockDesc.Contains(opt_message0) {
			blockDesc[opt_message0] = name.Friendly()
		}
		retMui = mui
	}
	var hasPrev bool
	if checks, e := reg.buildCheck(t, PreviousField, blockDesc, opt_previous); e != nil {
		err = errutil.Append(err, e)
	} else {
		hasPrev = checks
	}
	if _, e := reg.buildCheck(t, NextField, blockDesc, opt_next); e != nil {
		err = errutil.Append(err, e)
	}
	// block can have prev statement or have output
	hasOutput := blockDesc.Contains(opt_output)
	if !hasPrev && !hasOutput {
		if m, ok := r.PtrTo(t).MethodByName(OutputMethod); ok {
			if cnt := m.Type.NumOut(); cnt != 1 {
				err = errutil.Append(err, errutil.New("unexpected output count", cnt))
			} else {
				out := m.Type.Out(0)
				switch k := out.Kind(); k {
				default:
					err = errutil.Append(err, errutil.New("unexpected output type", k))
				case r.Interface:
					if basicInterface := r.TypeOf((interface{})(nil)); out != basicInterface {
						// basic interface is nil type.
						blockDesc[opt_output] = nil
					} else {
						blockDesc[opt_output] = toTypeName(out)
					}
				case r.Ptr:
					blockDesc[opt_output] = toTypeName(out.Elem())
				}
			}
		}
	}
	return
}

// evaluate the fields of the passed type to generate json usable by blockly initialization
// each field gets its own "options" Dict
func (reg *Registry) buildArgs(t r.Type, path string) (retMsg string, retArgs []Dict, retMui []mutationInput, err error) {
	var msgs []string
	var args []Dict
	var mui []mutationInput
	for i, cnt := 0, t.NumField(); i < cnt && err == nil; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := t.Field(i); len(f.PkgPath) == 0 &&
			f.Name != PreviousField &&
			f.Name != NextField {
			//
			if argDesc, argMui, e := reg.buildArgDesc(f, path); e != nil {
				err = errutil.Append(err, e)
			} else {
				args = append(args, argDesc)
				msgs = append(msgs, "%"+strconv.Itoa(len(args)))
				if argMui != nil {
					mui = append(mui, *argMui)
				}
			}
		}
	}
	if err == nil {
		retMsg = strings.Join(msgs, " ")
		retArgs = args
		retMui = mui
	}
	return
}

// return a description suitable for blockly describing the passed field.
// the format doesnt appear to be documented anywhere; derived through examples.
// additional returns an optional mutation input name/pair
func (reg *Registry) buildArgDesc(f r.StructField, path string) (argDesc Dict, mui *mutationInput, err error) {
	// tags take precedence over type info.
	// ( and some, ex. 'align' come only from tags. )
	argDesc = parseTags(string(f.Tag))
	name := path + pascalToCaps(f.Name)
	argDesc.Insert(opt_name, name)

	// check for some sort of enumerated type first.
	typeName := toTypeName(f.Type)
	if enumType, ok := reg.enums[typeName]; ok {
		argDesc.Insert(opt_type, field_dropdown)
		argDesc.Insert(opt_options, enumType.pairs)

	} else {
		// FIX -- how are we doing variables? always object productions ( ie. a get function )
		switch k := f.Type.Kind(); k {

		// slice of statements
		case r.Slice:
			if cs, e := reg.types.CheckType(f.Type.Elem()); e != nil {
				err = errutil.New(f.Name, e)
			} else {
				argDesc.Insert(opt_type, input_statement)
				if cs, ok := cs.GetConstraints(); ok {
					argDesc.Insert(opt_check, cs)
				}
			}

			// case r.Struct:
			// mutationName := typeName.String()
			// if cs, e := reg.types.CheckType(f.Type); e != nil {
			// 	err = errutil.New(f.Name, e)
			// } else {
			// 	if _, ok := reg.mutations[mutationName]; !ok {
			// 		err = errutil.New("unknown mutation", name, mutationName)
			// 	} else {
			// 		argDesc.Insert(opt_type, input_dummy)
			// 		mui = &mutationInput{InputName(name), mutationName, cs}
			// 	}
			// }

		// input containing another block
		case r.Ptr, r.Interface:
			if cs, e := reg.types.CheckType(f.Type); e != nil {
				err = errutil.New(f.Name, e)
			} else {
				if mutationName, ok := f.Tag.Lookup(tag_mutation); ok {
					if _, ok := reg.mutations[mutationName]; !ok {
						err = errutil.New("unknown mutation", name, mutationName)
					} else {
						argDesc.Insert(opt_type, input_dummy)
						mui = &mutationInput{InputName(name), mutationName, cs}
					}
				} else {
					argDesc.Insert(opt_type, input_value)
					if cs, ok := cs.GetConstraints(); ok {
						argDesc.Insert(opt_check, cs)
					}
				}
			}

		// a field of some sort ( ex. angle, checkbox, colour, date, dropdown, image, label, number, text, variable )
		case r.Bool:
			argDesc.Insert(opt_type, field_checkbox)

		case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
			argDesc.Insert(opt_type, field_number)
			argDesc.Insert(opt_precision, 1)

		case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
			argDesc.Insert(opt_type, field_number)
			argDesc.Insert(opt_precision, 1)
			argDesc.Insert(opt_min, 0)

		case r.Float32, r.Float64:
			argDesc.Insert(opt_type, field_number)

		case r.String:
			var argType string
			if argDesc.Contains(opt_readOnly) {
				argType = field_label
			} else {
				argType = field_input
			}
			argDesc.Insert(opt_type, argType)
			argDesc.Insert(opt_text, f.Name)

		default:
			switch r.PtrTo(f.Type) {
			case r.TypeOf((*time.Time)(nil)).Elem():
				argDesc.Insert(opt_type, field_date)
			default:
				err = errutil.New("field has unknown type", f.Name, f.Type)
				break
			}
		}
		// type FieldVariable string
		// type FieldImageDropdown []FieldImage
		// type FieldImage struct {
		// 	Width, Height int
		// 	Src           string
		// 	Alt           string
		// }
	}
	return
}
