#![no_implicit_prelude]
#![allow(missing_docs)] // TODO: change to warn/deny
#![allow(dead_code)] // TODO: remove

use ::anyhow::Result;
use ::clap::{self, Parser, Subcommand};
use ::std::{
    format,
    option::Option::{self, None, Some},
    string::String,
};

mod git;
mod runner;
mod stc;

#[derive(Parser)]
#[clap(about, version)]
struct Root {
    #[clap(subcommand)]
    subcommand: Command,
}

#[derive(Subcommand)]
#[clap()]
enum Command {
    #[clap(
        name = "clean",
        about = "Cleans any stc related refs and settings from repo.",
        override_usage = "stc clean"
    )]
    Clean,

    #[clap(
        name = "fix",
        about = "Adds, updates and deletes tracking refs if needed.",
        override_usage = "stc fix [<branch> [<base]]"
    )]
    Fix {
        /// name of the branch to fix
        #[clap(name = "branch")]
        branch: Option<String>,
        /// name of the base branch
        #[clap(name = "base")]
        base: Option<String>,
    },

    #[clap(
        name = "init",
        about = "Initializes the repo and tries to set stc refs for any non-default branches.",
        override_usage = "stc init"
    )]
    Init,

    #[clap(
        name = "push",
        about = "Sets remote branch head to what local branch head points to.",
        override_usage = "stc push"
    )]
    Push,

    #[clap(
        name = "rebase",
        about = "Rebases current branch on top of its base branch.",
        override_usage = "stc rebase"
    )]
    Rebase,

    #[clap(
        name = "start",
        about = "Starts a new branch off of current branch.",
        override_usage = "stc start <branch>"
    )]
    Start {
        /// name of the new branch to create
        #[clap(name = "branch")]
        branch: String,
    },

    #[clap(
        name = "sync",
        about = "Fetches all branches and tags and prunes deleted ones.",
        override_usage = "stc sync"
    )]
    Sync,
}

fn main() -> Result<()> {
    let root = Root::parse();
    let runner = runner::Runner::new("git")?;
    let stc = stc::Stc::new(runner);
    match root.subcommand {
        Command::Clean => stc.clean(),
        Command::Fix { branch, base } => stc.fix(branch, base),
        Command::Init => stc.init(),
        Command::Push => stc.push(),
        Command::Rebase => stc.rebase(),
        Command::Start { branch } => stc.start(branch),
        Command::Sync => stc.sync(),
    }
}
