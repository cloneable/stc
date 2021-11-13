package git

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type Runner struct {
	WorkDir       string
	Env           []string
	PrintCommands bool
	output        output
}

var _ Git = (*Runner)(nil)

func (r *Runner) Exec(args ...string) (Result, error) {
	c := exec.Command("git", args...)
	res := Result{
		Args: c.Args,
	}
	c.Env = r.Env
	c.Dir = r.WorkDir
	c.Stdin = nil
	c.Stdout = teeWriter{&res.Stdout, r.output.stdout()}
	c.Stderr = teeWriter{&res.Stderr, r.output.stderr()}

	cmdOut := r.output.stderr()
	if r.PrintCommands {
		cmdOut = teeWriter{os.Stderr, r.output.stdout()}
	}

	fmt.Fprintf(cmdOut, "### %s\n", c.Args)

	err := c.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		res.ExitCode = exitErr.ExitCode()
	}

	if err == nil {
		fmt.Fprintf(cmdOut, "### [OK]\n")
	} else {
		fmt.Fprintf(cmdOut, "### [err %d]: %v\n", res.ExitCode, err)
	}

	return res, err
}

func (r *Runner) DumpOutput() { r.output.dump() }

type teeWriter struct {
	a, b io.Writer
}

func (w teeWriter) Write(p []byte) (n int, err error) {
	aN, aErr := w.a.Write(p)
	if aErr != nil {
		return aN, aErr
	}
	bN, bErr := w.b.Write(p)
	if bErr != nil {
		return bN, bErr
	}
	return bN, nil
}

type output struct {
	mu     sync.Mutex
	writes []write
}

func (o *output) appendWrite(w write) (n int, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.writes = append(o.writes, w)
	return len(w.data), nil
}

type write struct {
	tag  int
	data []byte
}

func (o *output) stdout() io.Writer {
	return writerFunc(func(p []byte) (n int, err error) {
		buf := make([]byte, len(p))
		copy(buf, p)
		return o.appendWrite(write{tag: 0, data: p})
	})
}

func (o *output) stderr() io.Writer {
	return writerFunc(func(p []byte) (n int, err error) {
		buf := make([]byte, len(p))
		copy(buf, p)
		return o.appendWrite(write{tag: 1, data: buf})
	})
}

func (o *output) dump() {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, w := range o.writes {
		switch w.tag {
		case 0:
			os.Stdout.Write(w.data)
		case 1:
			os.Stderr.Write(w.data)
		}
	}
}

type writerFunc func(p []byte) (n int, err error)

func (f writerFunc) Write(p []byte) (n int, err error) { return f(p) }
