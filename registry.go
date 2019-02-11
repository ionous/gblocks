package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	r "reflect"
	"time"
)

type Registry struct {
	types     RegisteredTypes
	enums     RegisteredEnums
	mutations RegisteredMutations
}

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
	_, e := reg.enums.registerEnum(n)
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

func (reg *Registry) RegisterMutation(m interface{}, muiBlocks ...Mutation) (err error) {
	structType := r.TypeOf(m).Elem()
	mutationName := toTypeName(structType)
	return reg.mutations.RegisterMutation(mutationName, muiBlocks...)
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
							if registeredMutation, ok := reg.mutations.GetMutation(mutationName); !ok {
								e := errutil.New("unknown mutation", mutationName, "at", inputName)
								panic(e)
							} else if in, index := b.InputByName(inputName); index < 0 {
								e := errutil.New("unknown input", inputName)
								panic(e)
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

// could be a map, except maps arent ordered.
type mutationInput struct {
	mutationName TypeName  // the name as per RegisterMutation
	inputName    InputName //
	constraints  Constraints
}

// note: perhaps anonymous structs could be used to separate into args blocks
func (reg *Registry) buildBlockDesc(t r.Type, blockDesc Dict) (retMui []*mutationInput, err error) {
	name := toTypeName(t)
	blockDesc.Insert(opt_type, name)
	//
	if args, e := reg.buildArgs(t, ""); e != nil {
		err = errutil.Append(err, e)
	} else {
		if msg := args.message(); len(msg) > 0 {
			blockDesc.Insert(opt_message0, msg)
		}
		if len(args.list) > 0 {
			blockDesc.Insert(opt_args0, args.list)
		}
		if !blockDesc.Contains(opt_args0) && !blockDesc.Contains(opt_message0) {
			blockDesc[opt_message0] = name.Friendly()
		}
		retMui = args.mutations
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
func (reg *Registry) buildArgs(t r.Type, path string) (ret *argsOut, err error) {
	var args argsOut
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		if f := t.Field(i); len(f.PkgPath) == 0 && f.Name != PreviousField && f.Name != NextField {
			if e := reg.buildArgDesc(f, path, &args); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	if err == nil {
		ret = &args
	}
	return
}

// return a description suitable for blockly describing the passed field.
// the format doesnt appear to be documented anywhere; derived through examples.
// additional returns an optional mutation input name/pair
func (reg *Registry) buildArgDesc(f r.StructField, path string, out *argsOut) (err error) {
	argDesc := makeArg(f, path)

	// check for some sort of enumerated type first.
	if enumType, ok := reg.enums.GetEnum(argDesc.TypeName()); ok {
		argDesc.Insert(opt_type, field_dropdown)
		argDesc.Insert(opt_options, enumType.pairs)
		out.addArg(argDesc)
	} else {
		// FIX -- how are we doing variables? always object productions ( ie. a get function )
		switch k := f.Type.Kind(); k {

		// slice of statements
		case r.Slice:
			if cs, e := reg.types.CheckType(f.Type.Elem()); e != nil {
				err = errutil.New("invalid slice", argDesc, e)
			} else {
				argDesc.Insert(opt_type, input_statement)
				if cs, ok := cs.GetConstraints(); ok {
					argDesc.Insert(opt_check, cs)
				}
				out.addArg(argDesc)
			}

		// mutation
		case r.Struct:
			if e := reg.buildMutation(argDesc, out); e != nil {
				err = errutil.New("invalid mutation", argDesc, e)
			}

		// input containing another block
		case r.Ptr, r.Interface:
			if cs, e := reg.types.CheckType(f.Type); e != nil {
				err = errutil.New("invalid reference", argDesc, e)
			} else {
				argDesc.Insert(opt_type, input_value)
				if cs, ok := cs.GetConstraints(); ok {
					argDesc.Insert(opt_check, cs)
				}
				out.addArg(argDesc)
			}

		// a field of some sort ( ex. angle, checkbox, colour, date, dropdown, image, label, number, text, variable )
		case r.Bool:
			argDesc.Insert(opt_type, field_checkbox)
			out.addArg(argDesc)

		case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
			argDesc.Insert(opt_type, field_number)
			argDesc.Insert(opt_precision, 1)
			out.addArg(argDesc)

		case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
			argDesc.Insert(opt_type, field_number)
			argDesc.Insert(opt_precision, 1)
			argDesc.Insert(opt_min, 0)
			out.addArg(argDesc)

		case r.Float32, r.Float64:
			argDesc.Insert(opt_type, field_number)
			out.addArg(argDesc)

		case r.String:
			var argType string
			if argDesc.Contains(opt_readOnly) {
				argType = field_label
			} else {
				argType = field_input
			}
			argDesc.Insert(opt_type, argType)
			argDesc.Insert(opt_text, f.Name)
			out.addArg(argDesc)

		default:
			switch r.PtrTo(f.Type) {
			case r.TypeOf((*time.Time)(nil)).Elem():
				argDesc.Insert(opt_type, field_date)
				out.addArg(argDesc)

			default:
				err = errutil.New("unknown type", argDesc)
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

// t is a structType describing the mutation
// we expand all fields of the into block inputs, etc.
// if/when the field is NextStatement insert the mutation dummy element
func (reg *Registry) buildMutation(argDesc argDesc, out *argsOut) (err error) {
	inputName, mutationName := argDesc.InputName(), argDesc.TypeName()
	if !reg.mutations.Contains(mutationName) {
		err = errutil.New("unknown mutation")
	} else {
		t := argDesc.Type
		for i, cnt := 0, t.NumField(); i < cnt; i++ {
			switch f := t.Field(i); {
			// skip unexpected symbols ( only unexported symbols have a pkg path )
			case len(f.PkgPath) > 0:
			case f.Name == PreviousField:
				// error?
			case f.Name == NextField:
				fieldType := f.Type
				if cs, e := reg.types.CheckType(fieldType); e != nil {
					err = errutil.New("invalid reference", e)
				} else {
					argDesc.Insert(opt_type, input_dummy)
					argDesc.Insert(opt_mutation, mutationName)
					out.addMutation(argDesc, &mutationInput{mutationName, inputName, cs})
				}
			default:
				if e := reg.buildArgDesc(f, argDesc.Path, out); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
	}
	return
}
