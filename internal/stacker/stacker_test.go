package stacker

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"testing"
)

func newRepo(t *testing.T) string {
	t.Helper()

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

func cmd(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	c := exec.Command(
		name,
		args...,
	)
	c.Dir = dir
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		t.Log("STDOUT:\n", stdout.String())
		t.Log("STDERR:\n", stderr.String())
		t.Fatalf("command failed: %v", err)
	}
}

func TestInit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	gitDir := newRepo(t)
	cmd(t, gitDir, "touch", "README.md")
	cmd(t, gitDir, "git", "add", "README.md")
	cmd(t, gitDir, "git", "commit", "-m", "initial")
	cmd(t, gitDir, "git", "push", "-u", "origin")

	stkr := New(gitDir)
	stkr.Init(ctx, false)
}
