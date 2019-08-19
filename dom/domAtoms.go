package dom

type Mutation struct {
	Input string `xml:"name,attr"`
	Atoms Atoms  `xml:"atom,omitempty"`
}

type Atom struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Atoms []*Atom
