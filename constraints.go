package gblocks

type Constraints struct {
	constraints   []TypeName
	hasConnection bool
}

func (c *Constraints) GetConstraints() ([]TypeName, bool) {
	return c.constraints, c.hasConnection
}

// unlike in blockly proper, a single empty typename here means "any connection", rather than "no connection".
func (c *Constraints) AddConstraint(typeName TypeName) {
	if len(typeName) > 0 {
		c.constraints = append(c.constraints, typeName)
	}
	c.hasConnection = true
}

// does the next of this and the prev of that overlap?
// func (c *Constraints) ConnectsTo(c *Constraints) bool {
// }
