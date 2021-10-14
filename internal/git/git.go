package git

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

var protectedBranches = map[string]bool{
	"main":       true,
	"master":     true,
	"production": true,
	"release":    true,
}

type Git struct {
	path   string
	runner func(args ...string) (stdout bytes.Buffer, stderr bytes.Buffer, exitCode int, err error)
}

func (g Git) ListRefs() ([]Ref, error) {
	stdout, _, _, err := g.runner(
		"show-ref",
	)
	if err != nil {
		return nil, err
	}

	var refs []Ref
	scan := bufio.NewScanner(&stdout)
	for scan.Scan() {
		ref, err := ParseRef(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("cannot parse ref: %w", err)
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

func (g Git) RebaseOnto(onto, branchedCommit, branchName string) error {
	_, _, _, err := g.runner(
		"rebase",
		"--committer-date-is-author-date",
		"--onto",
		onto,
		branchedCommit,
		branchName,
	)
	if err != nil {
		return err
	}

	return nil
}

func (g Git) Push(branchName, remoteName, expectedCommit string) error {
	_, _, _, err := g.runner(
		"push",
		"--porcelain",
		fmt.Sprintf("--force-with-lease=%s:%s", branchName, expectedCommit),
		fmt.Sprintf("refs/heads/%s:refs/remotes/%s/%s", branchName, remoteName, branchName),
	)
	if err != nil {
		return err
	}

	return nil
}

func (g Git) GetRef(refName RefName) (Ref, error) {
	_, _, _, err := g.runner(
		"show-ref",
		"--verify",
		refName.String(),
	)
	if err != nil {
		return Ref{}, err
	}

	return Ref{}, nil
}

func (g Git) CreateRef(refName, commit string) error {
	_, _, _, err := g.runner(
		"update-ref",
		"--create-reflog",
		refName,
		commit,
		strings.Repeat("0", 40),
	)
	if err != nil {
		return err
	}

	return nil
}

func (g Git) UpdateRef(refName, newCommit, oldCommit string) error {
	_, _, _, err := g.runner(
		"update-ref",
		"--create-reflog",
		refName,
		newCommit,
		oldCommit,
	)
	if err != nil {
		return err
	}

	return nil
}

func (g Git) DeleteRef(refName, oldCommit string) error {
	_, _, _, err := g.runner(
		"update-ref",
		"-d",
		refName,
		oldCommit,
	)
	if err != nil {
		return err
	}

	return nil
}
