package git

import "path"

type Repo struct {
	branches []Branch
	labels   []Label
}

const labelRefNamespace = "stacker-label"

func NewLabel(name string) Label {
	return Label{
		ref: Ref{
			name: RefName(path.Join("refs", labelRefNamespace, name)),
		},
	}
}

func (l Label) Ref() Ref {
	return l.ref
}
