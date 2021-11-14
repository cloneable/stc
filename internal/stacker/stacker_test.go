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

type fileState struct {
	exists  bool
	content string
}

func assertFiles(t *testing.T, workDir string, files map[string]fileState) {
	t.Helper()
	for filePath, want := range files {
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
		if gotExists != want.exists {
			t.Errorf("%s: exists = %t, want %t", filePath, gotExists, want.exists)
		}
		if gotContent != want.content {
			t.Errorf("%s: content = %q, want %q", filePath, gotContent, want.content)
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

	assertFiles(t, workDir, map[string]fileState{
		".git/HEAD":                              {exists: true, content: "ref: refs/heads/test-branch-1\n"},
		".git/refs/heads/main":                   {exists: true, content: "2dcf4e535d10575b237f8fe0c7be220928d1ae6f\n"},
		".git/refs/heads/test-branch-1":          {exists: true, content: "2dcf4e535d10575b237f8fe0c7be220928d1ae6f\n"},
		".git/refs/stacker/base/test-branch-1":   {exists: true, content: "ref: refs/heads/main\n"},
		".git/refs/stacker/start/test-branch-1":  {exists: true, content: baseCommit},
		".git/refs/stacker/remote/test-branch-1": {},
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

	if err := stkr.Push(ctx); err != nil {
		t.Fatal(err)
	}

	assertFiles(t, workDir, map[string]fileState{
		".git/HEAD":                              {exists: true, content: "ref: refs/heads/test-branch-2\n"},
		".git/refs/heads/main":                   {exists: true, content: "2dcf4e535d10575b237f8fe0c7be220928d1ae6f\n"},
		".git/refs/heads/test-branch-2":          {exists: true, content: "cb53d6cb3801cd371346e474e7d4c6d79d6dd56c\n"},
		".git/refs/stacker/base/test-branch-2":   {exists: true, content: "ref: refs/heads/main\n"},
		".git/refs/stacker/start/test-branch-2":  {exists: true, content: "2dcf4e535d10575b237f8fe0c7be220928d1ae6f\n"},
		".git/refs/stacker/remote/test-branch-2": {exists: true, content: "cb53d6cb3801cd371346e474e7d4c6d79d6dd56c\n"},
	})
}
