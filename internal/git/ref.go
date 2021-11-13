package git

import (
	"fmt"
	"regexp"
	"strings"
)

type Repository struct {
	refs []Ref
}

type RefName string

func ParseRefName(name string) (RefName, error) {
	if !strings.HasPrefix(name, "refs/") {
		return "", fmt.Errorf("invalid ref name: %q", name)
	}
	return RefName(name), nil
}

func (rn RefName) String() string { return string(rn) }

type ObjectName string

func (c ObjectName) String() string { return string(c) }

type Ref struct {
	name         RefName
	objectName   ObjectName
	head         bool
	symRefTarget RefName
}

func (r Ref) Name() RefName          { return r.name }
func (r Ref) ObjectName() ObjectName { return r.objectName }
func (r Ref) SymRefTarget() RefName  { return r.symRefTarget }

var refLineRE = regexp.MustCompile("^([0-9a-f]{40}) (refs/.*)$")

func ParseRef(line string) (Ref, error) {
	groups := refLineRE.FindStringSubmatch(line)
	if len(groups) != 3 {
		return Ref{}, fmt.Errorf("invalid line: %q", line)
	}
	return Ref{
		name:       RefName(groups[2]),
		objectName: ObjectName(groups[1]),
	}, nil
}
