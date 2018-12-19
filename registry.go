package gblocks

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/errutil"
	r "reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Registry struct {
	types map[string]r.Type
	// given a golang type, find the corresponding blockly field
	// ex. FieldInput	 -> field_input', Blockly.FieldTextInput);
	fields map[r.Type]string
	enums  map[r.Type]enumInfo
}

type stringPair [2]string // display, uniquifier
type enumInfo struct {
	pairs []stringPair
}

var blockPtr = r.TypeOf((*Block)(nil))

func (reg *Registry) RegisterField(name string, f interface{}) {
	if reg.fields == nil {
		reg.fields = make(map[r.Type]string)
	}
	t := r.TypeOf(f).Elem()
	reg.fields[t] = name
}

// New - returns a pointer to passed type name
func (reg *Registry) New(name string) (ret r.Value, err error) {
	if t, ok := reg.types[name]; !ok {
		err = errutil.New("unknown type", name)
	} else {
		ret = r.New(t)
	}
	return
}

func (reg *Registry) RegisterBlocks(opt Options, blocks ...interface{}) (err error) {
	for _, b := range blocks {
		structType := r.TypeOf(b).Elem()
		structName := toTypeName(structType)
		var sub Options
		if opt, ok := opt[structName]; ok {
			if opt, ok := opt.(Options); ok {
				sub = opt
			}
		}
		if e := reg.registerBlock(structName, structType, sub); e != nil {
			e := errutil.New("error registering", structName, e)
			err = errutil.Append(err, e)
		}
	}
	return
}

func (reg *Registry) RegisterBlock(b interface{}, opt Options) (err error) {
	structType := r.TypeOf(b).Elem()
	structName := toTypeName(structType)
	return reg.registerBlock(structName, structType, opt)

}
func (reg *Registry) registerBlock(structName string, structType r.Type, opt Options) (err error) {
	//
	if reg.types == nil {
		reg.types = make(map[string]r.Type)
	} else if _, exists := reg.types[structName]; exists {
		panic("type already exists " + structName)
	}
	if opt == nil {
		opt = make(Options)
	}
	// delay the json initialization, so enums, etc. can all be registered first?
	// useful though to get options out.
	if e := reg.initJson(structType, opt); e != nil {
		err = e
	} else {
		fns := Options{
			"init": js.MakeFunc(func(obj *js.Object, _ []*js.Object) (ret interface{}) {
				b := &Block{Object: obj}

				if e := b.JsonInit(opt); e != nil {
					// hmmm...
				} else {
					// ideally, there'd be a custom input type for mutations;
					// but there's not.
					if ms, ok := opt[opt_mutation]; ok {
						for _, name := range ms.([]string) {
							in := b.GetInput(name)
							in.ForceMutation()
						}
					}
				}
				return // init has no return
			}),
		}
		reg.types[structName] = structType

		// the ui system has a mapping of block type names to block prototypes
		// standalone tests do not
		if blocks := js.Global.Get("Blockly").Get("Blocks"); blocks.Bool() {
			blocks.Set(structName, fns)
		}
	}
	return
}

// expects a map of intish to string
func (reg *Registry) RegisterEnum(n interface{}) (ret []stringPair, err error) {
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
func (reg *Registry) initJson(t r.Type, opt Options) (err error) {
	name := toTypeName(t)
	if _, ok := opt[opt_type]; !ok {
		opt[opt_type] = name
	}
	if msg, args, mutations, e := reg.makeArgs(t); e != nil {
		err = errutil.Append(err, e)
	} else {
		// stores the names of all mutations
		if len(mutations) > 0 {
			opt[opt_mutation] = mutations
		}
		if len(msg) > 0 && !opt.contains(opt_message0) {
			opt[opt_message0] = msg
		}
		if len(args) > 0 && !opt.contains(opt_args0) {
			opt[opt_args0] = args
		}
	}
	hasPrev, hasNext := opt.contains(opt_previousStatement), opt.contains(opt_nextStatement)
	if !hasPrev {
		if f, ok := t.FieldByName(previousStatement); ok {
			if c, e := constraintsForField(f); e != nil {
				err = errutil.Append(err, e)
			} else {
				opt[opt_previousStatement] = c
				hasPrev = true
			}
		}
	}
	if !hasNext {
		if f, ok := t.FieldByName(nextStatement); ok {
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
		if m, ok := r.PtrTo(t).MethodByName("Output"); ok {
			if cnt := m.Type.NumOut(); cnt != 1 {
				err = errutil.Append(err, errutil.New("unexpected output count", cnt))
			} else if out := m.Type.Out(0); out.Kind() != r.Interface {
				err = errutil.Append(err, errutil.New("unexpected output type", out.Kind()))
			} else if basicInterface := r.TypeOf((interface{})(nil)); out != basicInterface {
				// basic interface is nil type.
				opt[opt_output] = nil
			} else {
				outName := toTypeName(out)
				if _, ok := reg.types[outName]; !ok {
					err = errutil.Append(err, errutil.New("unknown output type", outName))
				} else {
					opt[opt_output] = outName
				}
			}
		}
	}

	return
}

// return a string, string array, or nil
func constraintsForField(f r.StructField) (ret interface{}, err error) {
	if s, ok := f.Tag.Lookup("check"); ok {
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
func (reg *Registry) makeArgs(t r.Type) (msg string, args []Options, muts []string, err error) {
	var msgs []string
	for i := 0; i < t.NumField() && err == nil; i++ {
		// skip unexpected symbols ( only unexported symbols have a pkg path )
		if f := t.Field(i); len(f.PkgPath) == 0 &&
			f.Name != previousStatement &&
			f.Name != nextStatement {

			// there can be only one.
			if opt, e := reg.makeOpt(f); e != nil {
				err = errutil.Append(err, e)
			} else {
				args = append(args, opt)
				msgs = append(msgs, "%"+strconv.Itoa(len(args)))
				// the only dummy inputs are mutations right now.
				if opt.contains(opt_mutation) {
					muts = append(muts, opt[opt_name].(string))
				}
			}
		}
	}
	msg = strings.Join(msgs, " ")
	return
}

func (reg *Registry) makeOpt(f r.StructField) (opt Options, err error) {
	// tags take precedence over type info.
	// ( and some, ex. 'align' come only from tags. )
	opt = parseTags(string(f.Tag))
	name := pascalToCaps(f.Name)
	opt.add(opt_name, name)

	if r.PtrTo(f.Type).Implements(r.TypeOf((*Mutation)(nil)).Elem()) {
		opt.add(opt_type, input_dummy)
		opt[opt_mutation] = true
	} else {
		// build 'type', 'check'
		switch k := f.Type.Kind(); k {

		// slice of statements.
		case r.Slice:
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
					var check []string
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
				if fieldType, ok := reg.fields[f.Type]; ok {
					opt[opt_type] = fieldType

				} else if enumType, ok := reg.enums[f.Type]; ok {
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
					// type FieldDropdown []string
					// type FieldImageDropdown []FieldImage
					// type FieldImage struct {
					// 	Width, Height int
					// 	Src           string
					// 	Alt           string
					// }
				}
			}
		}
	}
	return
}
