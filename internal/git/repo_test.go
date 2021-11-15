package git

import "testing"

func TestRefNameBranchName(t *testing.T) {
	if got, want := RefName("refs/foo/bar/moo").BranchName(), BranchName("moo"); got != want {
		t.Errorf(`RefName("refs/foo/bar/moo").BranchName() = %q, want %q`, got, want)
	}
}
