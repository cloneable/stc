package git

import (
	"bytes"
	"os"
	"os/exec"
)

type realGit struct {
	path string
}

func (g realGit) Run(args ...string) (stdout bytes.Buffer, stderr bytes.Buffer, exitCode int, err error) {
	git := exec.Cmd{
		Path: g.path,
		Args: append([]string{"git"}, args...),
		Env: []string{
			"HOME=" + os.Getenv("HOME"),
			"SSH_AUTH_SOCK=" + os.Getenv("SSH_AUTH_SOCK"),
		},
		Dir:    "",  // inherit
		Stdin:  nil, // /dev/null
		Stdout: &stdout,
		Stderr: &stderr,
	}

	err = git.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
		// err = nil
	}

	return
}
