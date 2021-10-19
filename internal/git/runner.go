package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	WorkDir string
	Env     []string
}

var _ Git = (*Runner)(nil)

type Result struct {
	Args     []string
	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	ExitCode int
}

func (r *Runner) exec(args ...string) (Result, error) {
	c := exec.Command("git", args...)
	res := Result{
		Args: c.Args,
	}
	// TODO: filter env? HOME, PATH, SSH_AUTH_SOCK, GIT_*
	c.Env = r.Env
	c.Dir = r.WorkDir
	c.Stdin = nil
	c.Stdout = &res.Stdout
	c.Stderr = &res.Stderr

	err := c.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		res.ExitCode = exitErr.ExitCode()
	} else if err != nil {
		return Result{}, err
	}

	fmt.Fprintf(os.Stderr, "GIT: %v (%d)\n", c.Args, res.ExitCode)

	return res, nil
}

func (r *Runner) Bare() (bool, error) {
	return false, nil
}

func (r *Runner) RepoRoot() (string, error) {
	return "", nil
}

func (r *Runner) ValidBranchName(name string) (bool, error) {
	return false, nil
}

func (g Runner) ListRefs() ([]Ref, error) {
	res, err := g.exec(
		"show-ref",
	)
	if err != nil {
		return nil, fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("command failed with %d", res.ExitCode)
	}

	var refs []Ref
	scan := bufio.NewScanner(&res.Stdout)
	for scan.Scan() {
		ref, err := ParseRef(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("cannot parse ref: %w", err)
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

func (g Runner) RebaseOnto(onto, branchedCommit, branchName string) error {
	res, err := g.exec(
		"rebase",
		"--committer-date-is-author-date",
		"--onto",
		onto,
		branchedCommit,
		branchName,
	)
	if err != nil {
		return fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return nil
}

func (g Runner) Push(branchName, remoteName, expectedCommit string) error {
	res, err := g.exec(
		"push",
		fmt.Sprintf("--force-with-lease=%s:%s", branchName, expectedCommit),
		fmt.Sprintf("refs/heads/%s:refs/remotes/%s/%s", branchName, remoteName, branchName),
	)
	if err != nil {
		return fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return nil
}

func (g Runner) GetRef(refName RefName) (Ref, error) {
	res, err := g.exec(
		"show-ref",
		"--verify",
		refName.String(),
	)
	if err != nil {
		return Ref{}, fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return Ref{}, fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return Ref{}, nil
}

func (g Runner) CreateRef(refName, commit string) error {
	res, err := g.exec(
		"update-ref",
		"--create-reflog",
		refName,
		commit,
		strings.Repeat("0", 40),
	)
	if err != nil {
		return fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return nil
}

func (g Runner) UpdateRef(refName, newCommit, oldCommit string) error {
	res, err := g.exec(
		"update-ref",
		"--create-reflog",
		refName,
		newCommit,
		oldCommit,
	)
	if err != nil {
		return fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return nil
}

func (g Runner) DeleteRef(refName, oldCommit string) error {
	res, err := g.exec(
		"update-ref",
		"-d",
		refName,
		oldCommit,
	)
	if err != nil {
		return fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return nil
}

func (g Runner) CheckBranchName(name string) error {
	res, err := g.exec(
		"check-ref-format",
		"--branch",
		name,
	)
	if err != nil {
		return fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return nil
}

func (g Runner) CurrentBranch() (string, error) {
	res, err := g.exec(
		"branch",
		"--show-current",
	)
	if err != nil {
		return "", fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return "", fmt.Errorf("command failed with %d", res.ExitCode)
	}
	return strings.TrimSpace(res.Stdout.String()), nil
}

func (g Runner) CreateBranch(name, baseName string) error {
	res, err := g.exec(
		"branch",
		name,
		baseName,
	)
	if err != nil {
		return fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return nil
}

func (g Runner) SwitchBranch(name string) error {
	res, err := g.exec(
		"switch",
		"--no-guess",
		name,
	)
	if err != nil {
		return fmt.Errorf("failure running command: %w", err)
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("command failed with %d", res.ExitCode)
	}

	return nil
}
