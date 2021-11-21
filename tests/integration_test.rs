#![allow(dead_code)]

use ::phf::{phf_map, Map};
use ::std::{
    fs::File,
    io::Write,
    path::{Path, PathBuf},
    process::{Command, Stdio},
};

// Make commits reproducible.
static TEST_ENV: Map<&'static str, &'static str> = phf_map! {
    "GIT_CONFIG_GLOBAL" => "/dev/null",
    "GIT_CONFIG_SYSTEM"=> "/dev/null",
    "GIT_AUTHOR_NAME"=> "tester",
    "GIT_AUTHOR_EMAIL" =>"tester@example.com",
    "GIT_AUTHOR_DATE"=> "1600000000 +0000",
    "GIT_COMMITTER_NAME"=> "tester",
    "GIT_COMMITTER_EMAIL"=> "tester@example.com",
    "GIT_COMMITTER_DATE"=> "1600000000 +0000",
};

fn run_git(repo_dir: &Path, args: &[&str]) {
    ::assert_cmd::Command::new("git")
        .current_dir(&repo_dir)
        .envs(TEST_ENV.entries())
        .args(args)
        .assert()
        .success();
}

fn write_file(repo_dir: &Path, file_path: &str, content: &str) {
    let mut f = File::create(repo_dir.join(file_path)).expect("cannot create file");
    f.write_all(content.as_bytes())
        .expect("cannot write content");
    f.flush().expect("cannot flush content");
}

fn new_stc_cmd(repo_dir: &Path) -> ::assert_cmd::Command {
    let mut cmd = ::assert_cmd::Command::cargo_bin("stc").expect("cannot resolve stc binary");
    cmd.current_dir(repo_dir.to_path_buf())
        .envs(TEST_ENV.entries());
    cmd
}

struct TestRepo {
    tempdir: ::tempfile::TempDir,
    origin_dir: PathBuf,
    clone_dir: PathBuf,
}

fn new_test_repo() -> TestRepo {
    let tempdir = ::tempfile::tempdir().expect("cannot create test tempdir");
    let origin_dir = tempdir.path().join("origin.git");
    let clone_dir = tempdir.path().join("clone");

    let status = Command::new("git")
        .args(&["init", "--bare", "--initial-branch=main"])
        .arg(&origin_dir)
        .current_dir(&tempdir)
        .envs(TEST_ENV.entries())
        .stdin(Stdio::null())
        .status()
        .expect("failed to run git command");
    assert!(status.success());

    let status = Command::new("git")
        .args(&["clone", "--origin=origin"])
        .arg(&origin_dir)
        .arg(&clone_dir)
        .current_dir(&tempdir)
        .envs(TEST_ENV.entries())
        .stdin(Stdio::null())
        .status()
        .expect("failed to run git command");
    assert!(status.success());

    write_file(&clone_dir, "README.md", "# test repo\n");
    run_git(&clone_dir, &["add", "README.md"]);
    run_git(&clone_dir, &["commit", "-m", "Initial commit"]);
    run_git(&clone_dir, &["push", "origin"]);

    TestRepo {
        tempdir,
        origin_dir,
        clone_dir,
    }
}

#[test]
fn test_stc_init() {
    let repo = new_test_repo();
    let mut stc = new_stc_cmd(&repo.clone_dir);

    stc.arg("init").assert().success();

    // TODO: check init settings
}
