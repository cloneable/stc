package stacker

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/cloneable/stacker/internal/git"
)

func newRepo(t *testing.T) string {
	t.Helper()

	// tmpDir, err := os.MkdirTemp(os.TempDir(), "stacker-test-*")
	// if err != nil {
	// 	t.Fatalf("cannot create tmp dir: %v", err)
	// }
	tmpDir := t.TempDir()

	pathOrigin := filepath.Join(tmpDir, "test-origin.git")
	git := exec.Command(
		"git",
		"init",
		"--bare",
		"--initial-branch=main",
		pathOrigin,
	)
	if err := git.Run(); err != nil {
		t.Fatalf("cannot init origin repo: %v", err)
	}
	pathRepo := filepath.Join(tmpDir, "test-repo")
	git = exec.Command(
		"git",
		"clone",
		"--origin=origin",
		pathOrigin,
		pathRepo,
	)
	if err := git.Run(); err != nil {
		t.Fatalf("cannot clone origin repo: %v", err)
	}

	return pathRepo
}

var cmdEnv = []string{
	// Make commit reproducible.
	"GIT_CONFIG_GLOBAL=/dev/null",
	"GIT_CONFIG_SYSTEM=/dev/null",
	"GIT_AUTHOR_NAME=tester",
	"GIT_AUTHOR_EMAIL=tester@example.com",
	"GIT_AUTHOR_DATE=1600000000 +0000",
	"GIT_COMMITTER_NAME=tester",
	"GIT_COMMITTER_EMAIL=tester@example.com",
	"GIT_COMMITTER_DATE=1600000000 +0000",
}

func cmd(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	c := exec.Command(
		name,
		args...,
	)
	c.Env = cmdEnv
	c.Dir = dir
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		t.Log("STDOUT:\n", stdout.String())
		t.Log("STDERR:\n", stderr.String())
		t.Fatalf("command failed: %v", err)
	}
}

func writeFile(t *testing.T, workDir, filePath, content string) {
	t.Helper()
	if err := os.WriteFile(path.Join(workDir, filePath), []byte(content), 0o600); err != nil {
		t.Fatalf("cannot write file %s: %v", filePath, err)
	}
}

func readFile(t *testing.T, workDir, filePath string) string {
	t.Helper()
	data, err := os.ReadFile(path.Join(workDir, filePath))
	if err != nil {
		t.Fatalf("cannot read file %s: %v", filePath, err)
	}
	return string(data)
}

type fileState string

func file(s string) *fileState    { f := fileState(s); return &f }
func (f *fileState) Exists() bool { return f != nil }
func (f *fileState) Content() string {
	if f == nil {
		return ""
	}
	return string(*f)
}

func assertFiles(t *testing.T, workDir string, files map[string]*fileState) {
	t.Helper()
	for filePath, wantFile := range files {
		fullPath := path.Join(workDir, filePath)
		_, err := os.Stat(fullPath)
		gotExists := true
		var gotContent string
		if err != nil {
			if !os.IsNotExist(err) {
				t.Fatalf("os.IsNotExist(%q): unexpected error: %v", fullPath, err)
			}
			gotExists = false
		} else {
			gotContent = readFile(t, workDir, filePath)
		}
		if got, want := gotExists, wantFile.Exists(); got != want {
			t.Errorf("%s: exists = %t, want %t", filePath, got, want)
		}
		if got, want := gotContent, wantFile.Content(); got != want {
			t.Errorf("%s: content = %q, want %q", filePath, gotContent, want)
		}
	}
}

func TestStackerStart(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	workDir := newRepo(t)

	writeFile(t, workDir, "main.txt", "0\n")
	cmd(t, workDir, "git", "add", "main.txt")
	cmd(t, workDir, "git", "commit", "-m", "main: state 0")
	cmd(t, workDir, "git", "push")

	stkr := Stacker{
		git: &git.Runner{
			WorkDir: workDir,
			Env:     cmdEnv,
		},
	}
	if err := stkr.Start(ctx, "test-branch-1"); err != nil {
		t.Fatal(err)
	}

	baseCommit := readFile(t, workDir, ".git/refs/heads/main")

	assertFiles(t, workDir, map[string]*fileState{
		".git/HEAD":                              file("ref: refs/heads/test-branch-1\n"),
		".git/refs/heads/test-branch-1":          file(baseCommit),
		".git/refs/stacker/base/test-branch-1":   file("ref: refs/heads/main\n"),
		".git/refs/stacker/start/test-branch-1":  file(baseCommit),
		".git/refs/stacker/remote/test-branch-1": nil,
	})
}

func TestStackerPush(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	workDir := newRepo(t)

	writeFile(t, workDir, "main.txt", "0\n")
	cmd(t, workDir, "git", "add", "main.txt")
	cmd(t, workDir, "git", "commit", "-m", "main: state 0")
	cmd(t, workDir, "git", "push")

	stkr := Stacker{
		git: &git.Runner{
			WorkDir: workDir,
			Env:     cmdEnv,
		},
	}

	if err := stkr.Start(ctx, "test-branch-2"); err != nil {
		t.Fatal(err)
	}

	writeFile(t, workDir, "test-branch-2.txt", "0\n")
	cmd(t, workDir, "git", "add", "test-branch-2.txt")
	cmd(t, workDir, "git", "commit", "-m", "test-branch-2: state 0")

	branchCommit := readFile(t, workDir, ".git/refs/heads/test-branch-2")

	assertFiles(t, workDir, map[string]*fileState{
		".git/refs/remotes/origin/test-branch-2": nil,
		".git/refs/stacker/remote/test-branch-2": nil,
	})

	if err := stkr.Push(ctx); err != nil {
		t.Fatal(err)
	}

	baseCommit := readFile(t, workDir, ".git/refs/heads/main")

	assertFiles(t, workDir, map[string]*fileState{
		".git/HEAD":                              file("ref: refs/heads/test-branch-2\n"),
		".git/refs/remotes/origin/test-branch-2": file(branchCommit),
		".git/refs/stacker/base/test-branch-2":   file("ref: refs/heads/main\n"),
		".git/refs/stacker/start/test-branch-2":  file(baseCommit),
		".git/refs/stacker/remote/test-branch-2": file(branchCommit),
	})
}
