package fakegit

import (
	"bytes"
)

type FakeGit struct{}

func (g FakeGit) Run(args ...string) (stdout bytes.Buffer, stderr bytes.Buffer, exitCode int, err error) {
	return bytes.Buffer{}, bytes.Buffer{}, 0, nil
}
