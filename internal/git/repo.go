package git

import (
	"encoding"
	"fmt"
	"path"
	"regexp"
	"strings"
)

type Repository struct {
	refs []Ref
}

type RefType int

const (
	_ RefType = iota

	TypeCommit
	TypeTree
	TypeBlob
	TypeTag
)

var _ (encoding.TextUnmarshaler) = (*RefType)(nil)
var _ (encoding.TextMarshaler) = (*RefType)(nil)

func (t *RefType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "commit":
		*t = TypeCommit
	case "tree":
		*t = TypeTree
	case "blob":
		*t = TypeBlob
	case "tag":
		*t = TypeTag
	default:
		return fmt.Errorf("unknown ref type: %s", text)
	}
	return nil
}

func (t RefType) MarshalText() ([]byte, error) {
	switch t {
	case TypeCommit:
		return []byte("commit"), nil
	case TypeTree:
		return []byte("tree"), nil
	case TypeBlob:
		return []byte("blob"), nil
	case TypeTag:
		return []byte("tag"), nil
	default:
		return nil, fmt.Errorf("unknown ref type: %d", t)
	}
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
	typ          RefType
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

type Tag struct {
	name string
}

func (t Tag) RefName() RefName {
	return RefName(path.Join("refs/tags", t.name))
}

type Head struct {
	name string
}

func (h Head) RefName() RefName {
	return RefName(path.Join("refs/heads", h.name))
}

type Remote string

type Branch struct {
	head   Head
	remote Remote
}

type RemoteBranch struct {
}

type Label struct {
	ref Ref
}

type BranchName string
