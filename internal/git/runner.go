package git

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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

	quotedArgs := make([]string, 0, len(args))
	for _, arg := range c.Args {
		quotedArgs = append(quotedArgs, shellQuote(arg))
	}
	cmdLine := strings.Join(quotedArgs, " ")

	fmt.Fprintf(r.output.stderr(), "\n### %s\n", cmdLine)

	err := c.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		res.ExitCode = exitErr.ExitCode()
	}

	if err == nil {
		fmt.Fprintf(os.Stderr, "[OK] %s\n", cmdLine)
		fmt.Fprintf(r.output.stderr(), "### [OK]\n")
	} else {
		fmt.Fprintf(os.Stderr, "[err %d] %s: %v\n", res.ExitCode, cmdLine, err)
		fmt.Fprintf(r.output.stderr(), "### [err %d]: %v\n", res.ExitCode, err)
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

func (o *output) appendWrite(tag int, buf []byte) (n int, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	data := make([]byte, len(buf))
	copy(data, buf)
	o.writes = append(o.writes, write{tag: tag, data: data})
	return len(data), nil
}

type write struct {
	tag  int
	data []byte
}

func (o *output) stdout() io.Writer {
	return writerFunc(func(p []byte) (n int, err error) {
		return o.appendWrite(1, p)
	})
}

func (o *output) stderr() io.Writer {
	return writerFunc(func(p []byte) (n int, err error) {
		return o.appendWrite(2, p)
	})
}

func (o *output) dump() {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, w := range o.writes {
		switch w.tag {
		case 1:
			os.Stdout.Write(w.data)
		case 2:
			os.Stderr.Write(w.data)
		}
	}
}

type writerFunc func(p []byte) (n int, err error)

func (f writerFunc) Write(p []byte) (n int, err error) { return f(p) }
