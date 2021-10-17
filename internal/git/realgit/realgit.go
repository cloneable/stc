package realgit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Run(dir string, args ...string) (stdout bytes.Buffer, stderr bytes.Buffer, exitCode int, err error) {
	git := exec.Command("git", args...)
	// TODO: filter env? HOME, PATH, SSH_AUTH_SOCK, GIT_*
	git.Env = nil
	git.Dir = dir
	git.Stdin = nil
	git.Stdout = &stdout
	git.Stderr = &stderr

	err = git.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
		// err = nil
	}

	fmt.Fprintf(os.Stderr, "GIT(%d): %v\n", exitCode, git.Args)

	return
}
