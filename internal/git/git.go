package git

import "bytes"

type Result struct {
	Args     []string
	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	ExitCode int
}

type Git interface {
	Exec(args ...string) (Result, error)
}
