package git

import (
	"bytes"
	"fmt"
	"testing"
)

type fakeGit struct {
	exitCode int
	stdout   string
	stderr   string
	err      error
}

func (g *fakeGit) Exec(args ...string) (Result, error) {
	return Result{
		Args:     args,
		ExitCode: g.exitCode,
		Stdout:   *bytes.NewBufferString(g.stdout),
		Stderr:   *bytes.NewBufferString(g.stderr),
	}, g.err
}

func (g *fakeGit) DumpOutput() {}

func TestSnapshotRepository(t *testing.T) {
	t.Parallel()

	g := &fakeGit{
		stdout: fields{
			Head:       true,
			ObjectName: "dabfcd577644ee74ad10e150720c29130e8dc5ea",
			RefName:    "refs/heads/main",
			ObjectType: TypeCommit,
			Track:      "=",
			Remote:     "origin",
			RemoteRef:  "refs/heads/main",
			SymRef:     "",
			Upstream:   "refs/remotes/origin/main",
		}.String(),
	}
	repo, err := SnapshotRepository(g)
	if err != nil {
		t.Fatal(err)
	}
	assertRepository(t, repo, Repository{
		refs: map[RefName]Ref{
			"refs/heads/main": {
				name:         "refs/heads/main",
				typ:          TypeCommit,
				objectName:   "dabfcd577644ee74ad10e150720c29130e8dc5ea",
				head:         true,
				symRefTarget: "",
				remote:       "origin",
				remoteRef:    "refs/heads/main",
				upstreamRef:  "refs/remotes/origin/main",
			},
		},
	})
}

func assertRepository(t *testing.T, actual, expected Repository) {
	t.Helper()
	if got, want := len(actual.refs), len(expected.refs); got != want {
		t.Errorf(`len(refs) = %d, want %d`, got, want)
		return
	}
	for name := range expected.refs {
		t.Run(fmt.Sprintf("ref: %q", name), func(t *testing.T) {
			t.Helper()
			assertRef(t, actual.refs[name], expected.refs[name])
		})
	}
}

func assertRef(t *testing.T, actual, expected Ref) {
	t.Helper()
	if got, want := actual.name, expected.name; got != want {
		t.Errorf(`name = %q, want %q`, got, want)
	}
	if got, want := actual.typ, expected.typ; got != want {
		t.Errorf(`typ = %v, want %v`, got, want)
	}
	if got, want := actual.objectName, expected.objectName; got != want {
		t.Errorf(`objectName = %q, want %q`, got, want)
	}
	if got, want := actual.head, expected.head; got != want {
		t.Errorf(`head = %t, want %t`, got, want)
	}
	if got, want := actual.symRefTarget, expected.symRefTarget; got != want {
		t.Errorf(`symRefTarget = %q, want %q`, got, want)
	}
	if got, want := actual.remote, expected.remote; got != want {
		t.Errorf(`remote = %q, want %q`, got, want)
	}
	if got, want := actual.remoteRef, expected.remoteRef; got != want {
		t.Errorf(`remoteRef = %q, want %q`, got, want)
	}
	if got, want := actual.upstreamRef, expected.upstreamRef; got != want {
		t.Errorf(`upstreamRef = %q, want %q`, got, want)
	}
}
