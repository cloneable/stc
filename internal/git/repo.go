package git

import (
	"encoding"
	"fmt"
	"sort"
	"strings"
)

const (
	StackerRefPrefix       = "refs/stacker/"
	StackerBaseRefPrefix   = StackerRefPrefix + "base/"
	StackerStartRefPrefix  = StackerRefPrefix + "start/"
	StackerRemoteRefPrefix = StackerRefPrefix + "remote/"

	branchRefPrefix = "refs/heads/"
	tagRefPrefix    = "refs/tags/"
)

type Repository struct {
	refs    map[RefName]Ref
	head    BranchName
	hasHead bool
}

func (r Repository) Head() (BranchName, bool) { return r.head, r.head != "" }
func (r Repository) AllRefs() []Ref {
	refs := make([]Ref, 0, len(r.refs))
	for _, ref := range r.refs {
		refs = append(refs, ref)
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].name < refs[j].name })
	return refs
}
func (r Repository) LookupRef(name RefName) (ref Ref, ok bool) { ref, ok = r.refs[name]; return }
func (r Repository) LookupBranch(name string) (BranchName, bool) {
	n := BranchName(name)
	_, ok := r.refs[n.RefName()]
	if !ok {
		return "", false
	}
	return n, true
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

func (n RefName) String() string { assert(n != ""); return string(n) }
func (n RefName) BranchName() BranchName {
	assert(n != "")
	if idx := strings.LastIndexByte(string(n), '/'); idx >= 0 {
		return BranchName(n[idx+1:])
	}
	return BranchName(n)
}

type ObjectName string

const NonExistantObject ObjectName = "0000000000000000000000000000000000000000"

func (n ObjectName) String() string { assert(n != ""); return string(n) }

type Ref struct {
	name         RefName
	typ          RefType
	objectName   ObjectName
	head         bool
	symRefTarget RefName
	remote       RemoteName
	remoteRef    RefName
	upstreamRef  RefName
}

func (r Ref) Name() RefName          { return r.name }
func (r Ref) ObjectName() ObjectName { return r.objectName }
func (r Ref) SymRefTarget() RefName  { return r.symRefTarget }
func (r Ref) Remote() RemoteName     { return r.remote }
func (r Ref) RemoteRefName() RefName { return r.remoteRef }
func (r Ref) UpstreamRef() RefName   { return r.upstreamRef }
func (r Ref) Stacker() bool          { return strings.HasPrefix(string(r.name), StackerRefPrefix) }

type BranchName string

func (n BranchName) String() string { assert(n != ""); return string(n) }

func (n BranchName) RefName() RefName              { return RefName(branchRefPrefix + n) }
func (n BranchName) StackerBaseRefName() RefName   { return RefName(StackerBaseRefPrefix + n) }
func (n BranchName) StackerStartRefName() RefName  { return RefName(StackerStartRefPrefix + n) }
func (n BranchName) StackerRemoteRefName() RefName { return RefName(StackerRemoteRefPrefix + n) }

type RemoteName string

func (n RemoteName) String() string { assert(n != ""); return string(n) }

func assert(cond bool) {
	if !cond {
		panic("bad assumption")
	}
}
