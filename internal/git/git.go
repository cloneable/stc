package git

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Result struct {
	Args     []string
	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	ExitCode int
}

type Git interface {
	Exec(args ...string) (Result, error)
}

func SnapshotRepository(g Git) (Repository, error) {
	res, err := g.Exec(
		"for-each-ref",
		"--format="+formatFields.formatFlagValue(),
	)
	if err != nil {
		return Repository{}, fmt.Errorf("cannot list refs: %w", err)
	}
	refs := make(map[RefName]Ref)
	var hasHead bool
	var head BranchName
	scan := bufio.NewScanner(&res.Stdout)
	for scan.Scan() {
		var f fields
		if err := json.Unmarshal(scan.Bytes(), &f); err != nil {
			return Repository{}, fmt.Errorf("cannot read JSON output: %w", err)
		}
		if err := f.validate(); err != nil {
			return Repository{}, fmt.Errorf("fields validation failed: %w", err)
		}
		r := f.ref()
		if r.head {
			if hasHead {
				return Repository{}, fmt.Errorf("found a second HEAD: %v and %v", head, r.name)
			}
			if !strings.HasPrefix(r.name.String(), branchRefPrefix) {
				return Repository{}, fmt.Errorf("found a non-branch HEAD: %v", r.name)
			}
			hasHead = true
			head = BranchName(strings.TrimPrefix(r.name.String(), branchRefPrefix))
		}
		refs[r.name] = r
	}
	return Repository{
		refs:    refs,
		head:    head,
		hasHead: hasHead,
	}, nil
}

type forEachRefField struct {
	jsonValue string
	re        *regexp.Regexp
}

type forEachRefSpec map[string]forEachRefField

func (s forEachRefSpec) formatFlagValue() string {
	if len(s) == 0 {
		return ""
	}
	var buf strings.Builder
	first := true
	buf.WriteByte('{')
	for name, spec := range s {
		if first {
			first = false
		} else {
			buf.WriteByte(',')
		}
		buf.WriteByte('"')
		buf.WriteString(name)
		buf.WriteByte('"')
		buf.WriteByte(':')
		buf.WriteString(spec.jsonValue)
	}
	buf.WriteByte('}')
	return buf.String()
}

type fields struct {
	Head       bool    `json:"head"`
	ObjectName string  `json:"objectname"`
	RefName    string  `json:"refname"`
	ObjectType RefType `json:"objecttype"`
	Track      string  `json:"track"`
	Remote     string  `json:"remote"`
	RemoteRef  string  `json:"remoteref"`
	SymRef     string  `json:"symref"`
}

func (f fields) String() string {
	var buf strings.Builder
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	if err := enc.Encode(f); err != nil {
		panic(err)
	}
	return buf.String()
}

var formatFields = forEachRefSpec{
	"head": {
		`%(if)%(HEAD)%(then)true%(else)false%(end)`,
		regexp.MustCompile("^(?:true|false)?$"),
	},
	"objectname": {
		`"%(objectname)"`,
		regexp.MustCompile("^[a-fA-F0-9]{40}$"),
	},
	"refname": {
		`"%(refname)"`,
		regexp.MustCompile("^refs/.+$"),
	},
	"objecttype": {
		`"%(objecttype)"`,
		regexp.MustCompile("^(?:commit|blob|tree|tag)$"),
	},
	"track": {
		`"%(upstream:trackshort)"`,
		regexp.MustCompile("^(?:=|>|<|<>)?$"),
	},
	"remote": {
		`"%(upstream:remotename)"`,
		regexp.MustCompile("^[^/]*$"),
	},
	"remoteref": {
		`"%(upstream:remoteref)"`,
		regexp.MustCompile("^(?:|refs/.+)$"),
	},
	"symref": {
		`"%(symref)"`,
		regexp.MustCompile("^(?:|refs/.+)$"),
	},
}

func (f fields) ref() Ref {
	return Ref{
		name:         RefName(f.RefName),
		objectName:   ObjectName(f.ObjectName),
		head:         f.Head,
		symRefTarget: RefName(f.SymRef),
	}
}

func (s forEachRefSpec) validate(name, value string) error {
	if !s[name].re.MatchString(value) {
		return fmt.Errorf("%s does not match pattern %v: %q", name, s[name].re.String(), value)
	}
	return nil
}

func (f fields) validate() error {
	// if err := formatFields.validate("head", f.Head); err != nil {
	// 	return err
	// }
	if err := formatFields.validate("objectname", f.ObjectName); err != nil {
		return err
	}
	if err := formatFields.validate("objectname", f.ObjectName); err != nil {
		return err
	}
	if err := formatFields.validate("refname", f.RefName); err != nil {
		return err
	}
	// if err := formatFields.validate("objecttype", f.ObjectType); err != nil {
	// 	return err
	// }
	if err := formatFields.validate("track", f.Track); err != nil {
		return err
	}
	if err := formatFields.validate("remote", f.Remote); err != nil {
		return err
	}
	if err := formatFields.validate("remoteref", f.RemoteRef); err != nil {
		return err
	}
	if err := formatFields.validate("symref", f.SymRef); err != nil {
		return err
	}
	return nil
}
