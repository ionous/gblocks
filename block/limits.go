package block

// lists types permitted to linking statements or terms
// becomes the "check" field in blockly block and item descriptions.
type Limits struct {
	Types    []string
	Connects bool
}

// return the types in a blockly friendly format:
// a nil pointer, a string, or an array of strings.
func (l *Limits) Check() (ret interface{}) {
	if l.Types != nil {
		if cnt := len(l.Types); cnt == 1 {
			ret = l.Types[0]
		} else {
			ret = l.Types
		}
	}
	return
}

func (l *Limits) IsUnlimited() bool {
	return l.Types == nil
}

// returns emtpy array
func MakeOffLimits() Limits {
	return Limits{[]string{}, true}
}

func MakeUnlimited() Limits {
	return Limits{nil, true}
}

func MakeLimits(types []string) Limits {
	return Limits{types, true}
}
