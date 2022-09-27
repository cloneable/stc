use assert_fs::{assert::PathAssert, fixture::PathChild};
use phf::{phf_map, Map};
use predicates::prelude::*;
use std::{
    env,
    fs::File,
    io::{Read, Write},
    path::{Path, PathBuf},
    process::{Command, Output, Stdio},
};

// Make commits reproducible.
static TEST_ENV: Map<&'static str, &'static str> = phf_map! {
    "GIT_CONFIG_GLOBAL" => "/dev/null",
    "GIT_CONFIG_SYSTEM" => "/dev/null",
    "GIT_AUTHOR_NAME" => "tester",
    "GIT_AUTHOR_EMAIL" => "tester@example.com",
    "GIT_AUTHOR_DATE" => "1600000000 +0000",
    "GIT_COMMITTER_NAME" => "tester",
    "GIT_COMMITTER_EMAIL" => "tester@example.com",
    "GIT_COMMITTER_DATE" => "1600000000 +0000",
};

fn run_git<I, S>(repo_dir: &Path, args: I) -> Output
where
    I: IntoIterator<Item = S>,
    S: AsRef<::std::ffi::OsStr>,
{
    let output = Command::new("git")
        .args(args)
        .current_dir(&repo_dir)
        .env_clear()
        .env("PATH", env::var("PATH").expect("$PATH not defined"))
        .envs(TEST_ENV.entries())
        .stdin(Stdio::null())
        .output()
        .expect("failed to run git command");
    assert!(output.status.success());
    output
}

fn write_file(repo_dir: &Path, file_path: &str, content: &str) {
    let mut f = File::create(repo_dir.join(file_path)).expect("cannot create file");
    f.write_all(content.as_bytes())
        .expect("cannot write content");
    f.flush().expect("cannot flush content");
}

fn read_file(p: &Path) -> String {
    let mut buf = String::new();
    ::std::fs::File::open(p)
        .expect("cannot open file")
        .read_to_string(&mut buf)
        .expect("cannot read file");
    buf
}

fn new_stc_cmd(repo_dir: &Path) -> ::assert_cmd::Command {
    let mut cmd = ::assert_cmd::Command::cargo_bin("stc").expect("cannot resolve stc binary");
    cmd.current_dir(repo_dir.to_path_buf())
        .env_clear()
        .env("PATH", env::var("PATH").expect("$PATH not defined"))
        .envs(TEST_ENV.entries());
    cmd
}

#[allow(dead_code)]
struct TestRepo {
    tempdir: ::assert_fs::TempDir,
    origin_dir: PathBuf,
    clone_dir: PathBuf,
}

fn new_test_repo(persistent: bool) -> TestRepo {
    let tempdir = ::assert_fs::TempDir::new()
        .expect("cannot create test tempdir")
        .into_persistent_if(persistent);
    let origin_dir = tempdir.path().join("origin.git");
    let clone_dir = tempdir.path().join("clone");

    run_git(
        tempdir.path(),
        [
            "init",
            "--bare",
            "--initial-branch=main",
            origin_dir.as_os_str().to_str().unwrap(),
        ],
    );
    run_git(
        tempdir.path(),
        [
            "clone",
            "--origin=origin",
            origin_dir.as_os_str().to_str().unwrap(),
            clone_dir.as_os_str().to_str().unwrap(),
        ],
    );

    write_file(&clone_dir, "README.md", "# test repo\n");
    run_git(&clone_dir, ["add", "README.md"]);
    run_git(&clone_dir, ["commit", "-m", "Initial commit"]);
    run_git(&clone_dir, ["push", "origin"]);

    TestRepo {
        tempdir,
        origin_dir,
        clone_dir,
    }
}

#[test]
fn test_stc_init() {
    let repo = new_test_repo(false);
    let mut stc = new_stc_cmd(&repo.clone_dir);

    stc.arg("init").assert().success();

    // TODO: check init settings
    // TODO: run again, check idempotency
}

#[test]
fn test_stc_start() {
    let repo = new_test_repo(false);
    let mut stc = new_stc_cmd(&repo.clone_dir);

    repo.tempdir
        .child("clone/.git/refs/stc")
        .assert(predicate::path::missing());
    let main_ref = read_file(repo.tempdir.child("clone/.git/refs/heads/main").path());

    stc.arg("start").arg("feat/test-branch").assert().success();

    repo.tempdir
        .child("clone/.git/refs/stc/start/feat/test-branch")
        .assert(main_ref);
    repo.tempdir
        .child("clone/.git/refs/stc/base/feat/test-branch")
        .assert("ref: refs/heads/main\n");
    repo.tempdir
        .child("clone/.git/refs/stc/remote/feat/test-branch")
        .assert(predicate::path::missing());
}

#[test]
fn test_stc_push() {
    let repo = new_test_repo(false);

    {
        let mut stc = new_stc_cmd(&repo.clone_dir);
        stc.arg("start").arg("feat/test-branch").assert().success();
    }

    write_file(&repo.clone_dir, "test-branch.txt", "test-branch #1\n");
    run_git(&repo.clone_dir, ["add", "test-branch.txt"]);
    run_git(&repo.clone_dir, ["commit", "-m", "test-branch.txt #1"]);

    let branch_ref = read_file(
        repo.tempdir
            .child("clone/.git/refs/heads/feat/test-branch")
            .path(),
    );

    repo.tempdir
        .child("clone/.git/refs/stc/remote/feat/test-branch")
        .assert(predicate::path::missing());

    {
        let mut stc = new_stc_cmd(&repo.clone_dir);
        stc.arg("push").assert().success();
    }

    repo.tempdir
        .child("clone/.git/refs/stc/remote/feat/test-branch")
        .assert(&branch_ref);

    {
        // No-op. Test idempotency.
        let mut stc = new_stc_cmd(&repo.clone_dir);
        stc.arg("push").assert().success();
    }

    repo.tempdir
        .child("clone/.git/refs/stc/remote/feat/test-branch")
        .assert(&branch_ref);
}
