package stacker

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cloneable/stacker/internal/git"
)

func newRepo(t *testing.T) string {
	t.Helper()

	// tmpDir, err := os.MkdirTemp("/Users/fb/tmp", "stacker-test-*")
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

func copyFile(t *testing.T, src, dst string) {
	t.Helper()

	srcFile, err := os.Open(src)
	if err != nil {
		t.Fatalf("cannot open src file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		t.Fatalf("cannot create dst file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		t.Fatalf("cannot copy file: %v", err)
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

	os.WriteFile(gitDir+"/README.md", []byte("# State 0\n"), 0x600)
	cmd(t, gitDir, "git", "add", "README.md")
	cmd(t, gitDir, "git", "commit", "-m", "state 0")
	cmd(t, gitDir, "git", "push", "-u", "origin")

	stkr := Stacker{
		git: &git.Runner{
			WorkDir: gitDir,
			Env:     cmdEnv,
		},
	}
	if err := stkr.Create(ctx, "test-branch-0"); err != nil {
		t.Fatal(err)
	}
}
