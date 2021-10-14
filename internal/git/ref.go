package git

import "path"

const labelRefNamespace = "stacker-label"

type RefName string

func (rn RefName) String() string { return string(rn) }

type Commit string

func (c Commit) String() string { return string(c) }

type Ref struct {
	name   RefName
	commit Commit
}

type Label struct {
	name string
	ref  Ref
}

func NewLabel(name string) Label {
	return Label{
		name: name,
		ref: Ref{
			name: RefName(path.Join("refs", labelRefNamespace, name)),
		},
	}
}

func (l Label) Ref() Ref {
	return l.ref
}
