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
	types map[TypeName]r.Type
	// given a golang type, find the corresponding blockly field
	// ex. FieldInput	 -> field_input', Blockly.FieldTextInput);
	enums     map[r.Type]enumInfo
	mutations map[string]mutationMap
}

// var TheRegistry Registry

// func RegisterBlocks(opt map[string]Options, blocks ...interface{}) error {
// 	return TheRegistry.registerBlocks(opt, blocks...)
// }

// // register a mapping of workspace block type to mutation ui block type.
// // name is the struct tag; it describes a mutation input
// func RegisterMutation(name string, pairs ...interface{}) error {
// 	return TheRegistry.registerMutation(name, pairs...)
// }

// func RegisterBlock(b interface{}, opt Options) error {
// 	structType := r.TypeOf(b).Elem()
// 	typeName := toTypeName(structType)
// 	return TheRegistry.registerType(typeName, structType, opt)
// }

// // RegisterEnum - expects a map of intish to string
// func RegisterEnum(n interface{}) error {
// 	// string pairs is returned for the sake of tests; we can ignore it here.
// 	_, err := TheRegistry.registerEnum(n)
// 	return err
// }

type mutationField struct {
	mutationName string
	inputName    InputName
	check        TypeName // fix, can this be an array.
}

type mutationType struct {
	workspaceType TypeName // name of the block in the workspace
	mutationBlock r.Type   // mutation ui block
}

//
type mutationMap []mutationType

func (mm *mutationMap) findMutationType(elType TypeName) (ret r.Type, okay bool) {
	for _, mtype := range *mm {
		if mtype.workspaceType == elType {
			ret = mtype.mutationBlock
			okay = true
			break
		}
	}
	return
}

// findAtomType - what kind of block in the main workspace does the mutation elemeent represent*
func (mm *mutationMap) findAtomType(mutationType TypeName) (ret TypeName, okay bool) {
	for _, mt := range *mm {
		// hmmm... skip the initial nil entry.
		if len(mt.workspaceType) > 0 && toTypeName(mt.mutationBlock) == mutationType {
			ret = mt.workspaceType
			okay = true
			break
		}
	}
	return
}

type stringPair [2]string // display, uniquifier
type enumInfo struct {
	pairs []stringPair
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

func (reg *Registry) registerBlocks(opt map[string]Options, blocks ...interface{}) (err error) {
	for _, b := range blocks {
		structType := r.TypeOf(b).Elem()
		typeName := toTypeName(structType)
		var sub Options
		if opt, ok := opt[typeName.String()]; ok {
			sub = opt
		}
		if e := reg.registerType(typeName, structType, sub); e != nil {
			e := errutil.New("error registering", typeName, e)
			err = errutil.Append(err, e)
		}
	}
	return
}

func (reg *Registry) registerMutation(name string, pairs ...interface{}) (err error) {
	if isEven := len(pairs)%1 == 0; !isEven {
		err = errutil.New("expected pairs of types; found", len(pairs))
	} else {
		if reg.mutations == nil {
			reg.mutations = make(map[string]mutationMap)
		} else if _, exists := reg.mutations[name]; exists {
			err = errutil.New("mutation already exists")
		}
		if err == nil {
			mutationTypes := make(mutationMap, 0, len(pairs)/2)
			for i := 0; i < len(pairs); i += 2 {
				var elType TypeName
				p1, p2 := r.TypeOf(pairs[i]), r.TypeOf(pairs[i+1])
				if p1 != nil {
					elType = toTypeName(p1.Elem())
				}
				if was, exists := mutationTypes.findMutationType(elType); exists {
					err = errutil.New("type mapping already exists", elType, was)
					break
				}
				mutationTypes = append(mutationTypes, mutationType{elType, p2.Elem()})
			}
			reg.mutations[name] = mutationTypes
		}
	}
	return
}

func (reg *Registry) contains(typeName TypeName) bool {
	_, ok := reg.types[typeName]
	return ok
}

func (reg *Registry) registerType(typeName TypeName, structType r.Type, opt Options) (err error) {
	//
	if reg.types == nil {
		reg.types = make(map[TypeName]r.Type)
	} else if _, exists := reg.types[typeName]; exists {
		panic("type already exists " + typeName)
	}
	if opt == nil {
		opt = make(Options)
	}
	// delay the json initialization, so enums, etc. can all be registered first?
	// useful though to get options out.
	if mutableInputs, e := reg.initJson(structType, opt); e != nil {
		err = e
	} else if blockly := js.Global.Get("Blockly"); !blockly.Bool() {
		err = errutil.New("blockly doesnt exist")
	} else if blocks := blockly.Get("Blocks"); !blocks.Bool() {
		err = errutil.New("blockly blocks dont exist")
	} else {
		// the ui system has "blocks" -- mapping type names to prototypes; standalone tests do not
		init := js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
			b := &Block{Object: obj}
			if e := b.JsonInit(opt); e != nil {
				panic(e)
			} else if cnt := len(mutableInputs); cnt > 0 {
				var blockNames []TypeName
				// all of the mutatable inputs
				for _, mutableInput := range mutableInputs {
					mutationName, inputName := mutableInput.mutationName, mutableInput.inputName
					if in, index := b.InputByName(inputName); index >= 0 {
						in.ForceMutation(mutationName)
					} else {
						panic("unexpected missing input" + inputName)
					}
					if mutationBlocks, ok := reg.mutations[mutationName]; !ok {
						panic("unknown mutation" + mutationName)
					} else {
						for _, mblock := range mutationBlocks {
							if len(mblock.workspaceType) > 0 {
								blockNames = append(blockNames, toTypeName(mblock.mutationBlock))
							}
						}
					}

				}
				b.SetMutator(NewMutator(blockNames))
			}
			return
		})
		var fns Options
		if len(mutableInputs) == 0 {
			fns = Options{
				"init": init,
			}
		} else {
			fns = Options{
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
		blocks.Set(string(typeName), fns)

		// create custom mutation block
		// FIX: if the same mutation appears in multiple blocks, this wont work
		if len(mutableInputs) > 0 {
			fns := Options{
				"init": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
					b := &Block{Object: obj}
					for _, mutableInput := range mutableInputs {
						label := NewFieldLabel(mutableInput.inputName.Friendly(), "")
						in := b.AppendStatementInput(mutableInput.inputName)
						in.AppendField(label.Field)
						in.SetCheck(mutableInput.check)
					}
					return
				}),
			}
			mutationName := string(typeName) + "$mutation"
			blocks.Set(mutationName, fns)
		}
	}
	return
}

func (reg *Registry) registerEnum(n interface{}) (ret []stringPair, err error) {
	var pairs []stringPair
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
			pair := stringPair{display, unique}
			pairs = append(pairs, pair)
		}
		if reg.enums == nil {
			reg.enums = make(map[r.Type]enumInfo)
		}
		reg.enums[keyType] = enumInfo{pairs: pairs}
		ret = pairs
	}
	return
}

// note: perhaps anonymous structs could be used to separate into args blocks
func (reg *Registry) initJson(t r.Type, opt Options) (ret []mutationField, err error) {
	name := toTypeName(t)
	if _, ok := opt[opt_type]; !ok {
		opt[opt_type] = name
	}
	if msg, args, mutations, e := reg.makeArgs(t, ""); e != nil {
		err = errutil.Append(err, e)
	} else {
		if len(msg) > 0 && !opt.contains(opt_message0) {
			opt[opt_message0] = msg
		}
		if len(args) > 0 && !opt.contains(opt_args0) {
			opt[opt_args0] = args
		}
		if !opt.contains(opt_args0) && !opt.contains(opt_message0) {
			spaces := name.Friendly()
			// FIX: this is ill advised.
			// if mutations blocks can be derived from the workspace blocks that they create
			// then we could use the label thse blocks with the workspace block names
			if strings.HasSuffix(spaces, " el") {
				spaces = spaces[:len(spaces)-3]
			}
			opt[opt_message0] = spaces
		}

		ret = mutations
	}
	hasPrev, hasNext := opt.contains(opt_previousStatement), opt.contains(opt_nextStatement)
	if !hasPrev {
		if f, ok := t.FieldByName(PreviousField); ok {
			if c, e := constraintsForField(f); e != nil {
				err = errutil.Append(err, e)
			} else {
				opt[opt_previousStatement] = c
				hasPrev = true
			}
		}
	}
	if !hasNext {
		if f, ok := t.FieldByName(NextField); ok {
			if c, e := constraintsForField(f); e != nil {
				err = errutil.Append(err, e)
			} else {
				opt[opt_nextStatement] = c
				hasNext = true
			}
		}
	}
	// block can have prev statement or have output
	hasOutput := opt.contains(opt_output)
	if !hasPrev && !hasOutput {
		if m, ok := r.PtrTo(t).MethodByName(OutputMethod); ok {
			if cnt := m.Type.NumOut(); cnt != 1 {
				err = errutil.Append(err, errutil.New("unexpected output count", cnt))
			} else if out := m.Type.Out(0); out.Kind() != r.Interface {
				err = errutil.Append(err, errutil.New("unexpected output type", out.Kind()))
			} else if basicInterface := r.TypeOf((interface{})(nil)); out != basicInterface {
				// basic interface is nil type.
				opt[opt_output] = nil
			} else {
				opt[opt_output] = toTypeName(out)
			}
		}
	}

	return
}

// return a string, string array, or nil
func constraintsForField(f r.StructField) (ret interface{}, err error) {
	if s, ok := f.Tag.Lookup(opt_check); ok {
		ret = splitCommas(s)
	} else {
		switch f.Type.Kind() {
		case r.Interface:
			ret = nil
		case r.Ptr:
			ret = toTypeName(f.Type.Elem())
		default:
			err = errutil.New("statement has unexpected type", f.Type)
		}
	}
	return
}

// evaluate the fields of the passed type to generate json usable by blockly initialization
// each field gets its own "options" Options
func (reg *Registry) makeArgs(t r.Type, path string) (msg string, args []Options, mutableInputs []mutationField, err error) {
	var msgs []string
	for i, cnt := 0, t.NumField(); i < cnt && err == nil; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := t.Field(i); len(f.PkgPath) == 0 &&
			f.Name != PreviousField &&
			f.Name != NextField {

			// there can be only one.
			if opt, e := reg.makeOpt(f, path); e != nil {
				err = errutil.Append(err, e)
			} else {
				args = append(args, opt)
				msgs = append(msgs, "%"+strconv.Itoa(len(args)))
				//.
				if m, ok := opt[opt_mutation]; ok {
					m := m.(string)
					// ugh.
					input := r.ValueOf(opt[opt_name]).Convert(
						r.TypeOf(((*InputName)(nil))).Elem())

					var check TypeName
					if c, ok := opt[opt_check]; ok {
						check = c.(TypeName)
					} else if types, ok := reg.mutations[m]; !ok {
						panic("couldnt find mutation named " + m)
					} else if inputType, ok := types.findMutationType(""); !ok {
						panic("error " + m)
					} else {
						check = toTypeName(inputType)
					}

					mutableInput := mutationField{
						mutationName: m,
						inputName:    input.Interface().(InputName),
						check:        check,
					}
					mutableInputs = append(mutableInputs, mutableInput)
				}
			}
		}
	}
	msg = strings.Join(msgs, " ")
	return
}

func (reg *Registry) makeOpt(f r.StructField, path string) (opt Options, err error) {
	// tags take precedence over type info.
	// ( and some, ex. 'align' come only from tags. )
	opt = parseTags(string(f.Tag))
	name := pascalToCaps(f.Name)
	opt.add(opt_name, path+name)

	// build 'type', 'check'
	switch k := f.Type.Kind(); k {

	// slice of statements or a mutation.
	case r.Slice:
		if s, ok := f.Tag.Lookup(tag_mutation); ok {
			opt.add(opt_type, input_dummy)
			opt[opt_mutation] = s
			// note: -- the dummy input cant have a check type
			// if !opt.contains(opt_check) {
			// 	if types, ok := reg.mutations[s]; !ok {
			// 		panic("couldnt find mutation named " + s)
			// 	} else if inputType, ok := types[""]; !ok {
			// 		panic("error")
			// 	} else {
			// 		inputType := toTypeName(inputType)
			// 		opt.add(opt_check, inputType)
			// 	}
			// }
		} else {
			opt.add(opt_type, input_statement)
			switch elType := f.Type.Elem(); elType.Kind() {
			case r.Interface:
				if basicInterface := r.TypeOf((interface{})(nil)); elType != basicInterface {
					opt.add(opt_check, toTypeName(elType))
				}
			case r.Ptr:
				opt.add(opt_check, toTypeName(elType))
			default:
				err = errutil.New(f.Name, "has unexpected type", elType)
			}
		}

	// input containing another block.
	case r.Ptr:
		switch elType := f.Type.Elem(); elType.Kind() {
		case r.Struct:
			opt.add(opt_type, input_value)
			opt.add(opt_check, toTypeName(elType))
		default:
			err = errutil.New(f.Name, "has unexpected type", elType)
		}

	// input of one or more block type
	case r.Interface:
		opt.add(opt_type, input_value)

		if !opt.contains(opt_check) {
			// the basic interface mean no check: it accepts everything
			if basicInterface := r.TypeOf((interface{})(nil)); f.Type != basicInterface {
				var check []TypeName
				// run through all the registered types to see if they implement the interface
				for n, t := range reg.types {
					if t.Implements(f.Type) {
						check = append(check, n)
					}
				}
				if len(check) > 0 {
					opt[opt_check] = check
				}
			}
		}

	// a field of some sort ( ex. angle, checkbox, colour, date, dropdown, image, label, number, text, variable )
	// FIX -- how are we doing variables? always object productions ( ie. a get function )
	default:
		if !opt.contains(opt_type) {
			if enumType, ok := reg.enums[f.Type]; ok {
				opt[opt_type] = field_dropdown
				opt.add(opt_options, enumType.pairs)

			} else {
				switch f.Type.Kind() {
				case r.Bool:
					opt[opt_type] = field_checkbox

				case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
					opt[opt_type] = field_number //setConstraints
					opt.add(opt_precision, 1)

				case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
					opt[opt_type] = field_number
					opt.add(opt_precision, 1)
					opt.add(opt_min, 0)

				case r.Float32, r.Float64:
					opt[opt_type] = field_number

				case r.String:
					if opt.contains(opt_readOnly) {
						opt[opt_type] = field_label
					} else {
						opt[opt_type] = field_input
					}
					opt.add(opt_text, f.Name)

				default:
					switch r.PtrTo(f.Type) {
					case r.TypeOf((*time.Time)(nil)).Elem():
						opt[opt_type] = field_date
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
		}
	}
	return
}
