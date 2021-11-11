package git

import (
	"fmt"
	"os"
	"os/exec"
)

type Runner struct {
	WorkDir string
	Env     []string
}

var _ Git = (*Runner)(nil)

func (r *Runner) Exec(args ...string) (Result, error) {
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
