package git

import (
	"fmt"
	"path"
	"regexp"
)

const labelRefNamespace = "stacker-label"

type RefName string

func (rn RefName) String() string { return string(rn) }

type Commit string

func (c Commit) String() string { return string(c) }

type Ref struct {
	name   RefName
	commit Commit
}

var refLineRE = regexp.MustCompile("^([0-9a-f]{40}) (refs/.*)$")

func ParseRef(line string) (Ref, error) {
	groups := refLineRE.FindStringSubmatch(line)
	if len(groups) != 3 {
		return Ref{}, fmt.Errorf("invalid line: %q", line)
	}
	return Ref{
		name:   RefName(groups[2]),
		commit: Commit(groups[1]),
	}, nil
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
