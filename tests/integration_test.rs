use ::phf::{phf_map, Map};
use ::std::{
    env,
    fs::File,
    io::Write,
    path::{Path, PathBuf},
    process::{Command, Output, Stdio},
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
    tempdir: ::tempfile::TempDir,
    origin_dir: PathBuf,
    clone_dir: PathBuf,
}

fn new_test_repo() -> TestRepo {
    let tempdir = ::tempfile::tempdir().expect("cannot create test tempdir");
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
    let repo = new_test_repo();
    let mut stc = new_stc_cmd(&repo.clone_dir);

    stc.arg("init").assert().success();

    // TODO: check init settings
    // TODO: run again, check idempotency
}

#[test]
fn test_stc_start() {
    let repo = new_test_repo();
    let mut stc = new_stc_cmd(&repo.clone_dir);

    stc.arg("start").arg("test-branch").assert().success();

    // TODO: check test-branch refs
}

#[test]
fn test_stc_push() {
    let repo = new_test_repo();

    {
        let mut stc = new_stc_cmd(&repo.clone_dir);
        stc.arg("start").arg("test-branch").assert().success();
    }

    {
        let mut stc = new_stc_cmd(&repo.clone_dir);
        stc.arg("push").assert().success();
    }

    // TODO: check test-branch refs + remote
}
