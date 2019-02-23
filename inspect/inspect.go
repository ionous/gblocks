package inspect

import (
	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
	r "reflect"
	"strings"
	"time"
)

type Class int

//go:generate stringer -type=Class
const (
	Unknown Class = iota
	Bool
	Int
	Uint
	Float
	Label
	Text
	Date
	Group
	// connection types:
	Input
	Statements
	NextLink
	PrevLink
)

const (
	NextStatement     = "NextStatement"
	PreviousStatement = "PreviousStatement"
)

func (cls Class) Connects() bool {
	return cls >= Input
}
func (cls Class) IsLink() bool {
	return cls == NextLink || cls == PrevLink
}

type Item struct {
	Name    block.Item
	Type    r.Type
	Class   Class
	Options block.Dict
}

type Items []*Item

func (it *Item) String() string {
	return strings.Join([]string{
		it.Name.String(),
		it.Type.Name(),
		it.Class.String(),
	}, ":")
}

func MakeItems(t r.Type) (out Items, err error) {
	VisitItems(t, func(it *Item, e error) bool {
		if e != nil {
			err = errutil.Append(err, e)
		} else {
			out = append(out, it)
		}
		return true // keep going regardless of error
	})
	return
}

func VisitItems(t r.Type, visit func(it *Item, e error) (keepGoing bool)) {
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		if f := t.Field(i); len(f.PkgPath) == 0 {
			it, e := MakeItem(f)
			if keepGoing := visit(it, e); !keepGoing {
				break
			}
		}
	}
}

func MakeItem(f r.StructField) (ret *Item, err error) {
	itemName, itemType := block.ItemFromField(f), f.Type
	options := parseTags(string(f.Tag))
	if cls := classify(itemType); cls == Unknown {
		err = errutil.New("unknown item at", f)
	} else {
		if cls == Text {
			if options.Contains("readonly") {
				cls = Label
			}
		} else if cls == Input {
			if f.Name == block.PreviousStatement {
				cls = PrevLink
			} else if f.Name == block.NextStatement {
				cls = NextLink
			}
		}
		ret = &Item{itemName, itemType, cls, options}
	}
	return
}

func classify(t r.Type) (retClass Class) {
	switch k := t.Kind(); k {
	case r.Slice:
		retClass = Statements
	case r.Struct:
		retClass = Group
	case r.Interface:
		retClass = Input
	case r.Ptr:
		if elType := t.Elem(); elType.Kind() == r.Struct {
			retClass = Input
		}
	case r.Bool:
		retClass = Bool
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		retClass = Int
	case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
		retClass = Uint
	case r.Float32, r.Float64:
		retClass = Float
	case r.String:
		retClass = Text
	default:
		switch t {
		case r.TypeOf((*time.Time)(nil)).Elem():
			retClass = Date
		}
	}
	// type FieldVariable string
	// type FieldImageDropdown []FieldImage
	// type FieldImage struct {
	// 	Width, Height int
	// 	Src           string
	// 	Alt           string
	// }
	return
}
