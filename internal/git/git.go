package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var protectedBranches = map[string]bool{
	"main":       true,
	"master":     true,
	"production": true,
	"release":    true,
}

type Git struct {
	path string
}

func (g Git) ListRefs() ([]Ref, error) {
	stdout, _, _, err := g.Run(
		"show-ref",
	)
	if err != nil {
		return nil, err
	}

	var refs []Ref
	scan := bufio.NewScanner(&stdout)
	for scan.Scan() {
		line := scan.Text()
		frags := strings.SplitN(line, " ", 2) // XXX regexp parse, check len
		refs = append(refs, Ref{
			name:   RefName(frags[1]),
			commit: Commit(frags[0]),
		})
	}

	return refs, nil
}

func (g Git) RebaseOnto(onto, branchedCommit, branchName string) error {
	_, _, _, err := g.Run(
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
	_, _, _, err := g.Run(
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
	_, _, _, err := g.Run(
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
	_, _, _, err := g.Run(
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
	_, _, _, err := g.Run(
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
	_, _, _, err := g.Run(
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

func (g Git) Run(args ...string) (stdout bytes.Buffer, stderr bytes.Buffer, exitCode int, err error) {
	git := exec.Cmd{
		Path: g.path,
		Args: append([]string{"git"}, args...),
		Env: []string{
			"HOME=" + os.Getenv("HOME"),
		},
		Dir:    "",  // inherit
		Stdin:  nil, // /dev/null
		Stdout: &stdout,
		Stderr: &stderr,
	}

	err = git.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
		err = nil
	}

	return
}
