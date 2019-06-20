package enum

import (
	r "reflect" // for inspecting enumerated maps
	"sort"

	"github.com/ionous/errutil"
)

type Pair []string // display, uniquifier

type Pairs struct {
	mapping map[string][]Pair
}
type stringer interface{ String() string }

// name is generally pascal cased; returns nil if nothing registered
func (ep *Pairs) GetPairs(name string) []Pair {
	return ep.mapping[name]
}

// add one or more enum; see AddEnum
func (ep *Pairs) AddEnums(mappings ...interface{}) (err error) {
	for _, m := range mappings {
		if _, e := ep.AddEnum(m); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// add a map of enum value to string.
// gblocks will turn fields of the enum's key type into a dropdown box.
// the dropdown box displays the the vaues of the mapping.
// ex.
//  type Example int
//  const (
//  	DefaultChoice Example = iota
//  	AlternativeChoice
//  )
//  ... AddEnum(map[Example]string{
//  	DefaultChoice:     "default",
//	 	AlternativeChoice: "alt",
//  }
func (ep *Pairs) AddEnum(mapping interface{}) (ret []Pair, err error) {
	if src, srcType := r.ValueOf(mapping), r.TypeOf(mapping); srcType.Kind() != r.Map {
		err = errutil.New("expected a mapping of enum key (integer based) to enum value (string based)", srcType)
	} else {
		keyType, valueType := srcType.Key(), srcType.Elem()
		if keyType.Kind() != r.Int {
			err = errutil.New("expected an integer underlying the key type", keyType)
		} else if valueType.Kind() != r.String {
			err = errutil.New("expected a string underlying the value type", valueType)
		} else if pairs, e := makePairs(src); e != nil {
			err = e
		} else {
			ret, err = ep.addPairs(keyType.Name(), pairs)
		}
	}
	return
}

// add one or more simplified enums; see AddList
func (ep *Pairs) AddLists(choices ...interface{}) (err error) {
	for _, cs := range choices {
		if _, e := ep.AddList(cs); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// provides for a simpler way of declaring enums
// enum types are declared as string
// and enum choices are specified as an array of strings.
// ( translation would use the english text as a key )
func (ep *Pairs) AddList(choices interface{}) (ret []Pair, err error) {
	list := r.ValueOf(choices)
	if listType := list.Type(); listType.Kind() != r.Slice {
		err = errutil.New("expected a list of choices", listType)
	} else {
		keyType := listType.Elem()
		if keyType.Kind() != r.String {
			err = errutil.New("expected a string underlying the key type", keyType)
		} else if pairs, e := makePairsFromList(list); e != nil {
			err = e
		} else {
			ret, err = ep.addPairs(keyType.Name(), pairs)
		}
	}
	return
}
func (ep *Pairs) addPairs(name string, pairs []Pair) (ret []Pair, err error) {
	if ep.mapping == nil {
		ep.mapping = map[string][]Pair{name: pairs}
	} else {
		ep.mapping[name] = pairs
		ret = pairs
	}
	return
}

func makePairs(src r.Value) (ret []Pair, err error) {
	if keys := src.MapKeys(); len(keys) == 0 {
		err = errutil.New("expected at least once choice")
	} else {
		var pairs []Pair
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].Int() < keys[j].Int()
		})
		for _, key := range keys {
			var unique string
			if stringer, ok := key.Interface().(stringer); ok {
				//unique := fmt.Sprint(key)
				unique = stringer.String()
			} else {
				unique = key.String()
			}
			display := src.MapIndex(key).String()
			pair := Pair{display, unique}
			pairs = append(pairs, pair)
		}
		ret = pairs
	}
	return
}

func makePairsFromList(list r.Value) (ret []Pair, err error) {
	if cnt := list.Len(); cnt == 0 {
		err = errutil.New("expected at least once choice")
	} else {
		pairs := make([]Pair, cnt, cnt)
		for i := 0; i < cnt; i++ {
			choice := list.Index(i).String()
			pairs[i] = Pair{choice, choice}
		}
		ret = pairs
	}
	return
}
