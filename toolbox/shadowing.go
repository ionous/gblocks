package toolbox

// Shadowing - when creating xml from golang types should we create shadow blocks
// https://developers.google.com/blockly/guides/configure/web/toolbox#shadow_blocks
type Shadowing int

const (
	NoShadow Shadowing = iota
	IsShadow
	SubShadow
)

// Children of shadows or subshadows are shadows
func (s Shadowing) Children() Shadowing {
	out := s
	if out == SubShadow {
		out = IsShadow // upgrade shadowing; otherwise no change
	}
	return out
}
