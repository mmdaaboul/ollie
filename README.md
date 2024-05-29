# Ollie

## The TC Helper Tool

I (Jake) made this tool to help streamline common tasks at T-Cetra that take longer than I think they should due to lots of dead time waiting for things to open or load that I shouldn't need to wait on.

Everything for this project is supposed to be a common task and easier to complete when using it.

---

## Installation

### Install with Go

This is a Go project to make it easily transferrable between different OS's and shells as it can be compiled down to a single binary.

1. Clone the repo
2. Install by `cd`'ing into the repo
3. Run `go install .`

If you have Go already installed and your `$GOPATH` already in your environment `$PATH`, you can just run it anywhere with `ollie`

### Download the binary

I am not including the binary in the repo...yet. But if you just want the single binary code, reach out to me and I will give you the built version for whatever OS you need it for.

I will already have a BASH/ZSH (WSL) and Powershell (Windows) version built and ready.

Then, you will need to put the binary into a folder that is in your path, then reload your shell and run it with `ollie`

---

## What does it do?

Currently it has 3 main functions, Deploying (via git tags), Zookeeper (Reading only for now, see [Roadmap](#Roadmap)), and Pretty Printing the Git Tags

All three of these functions are presented as a choice to you when you run `ollie`, however you can skip the first prompt if you pass what you want to do

```bash
$ ollie
```

**OR**

```bash
$ ollie [deploy|zookeeper|printTags|db]
```

You can then go through the tool answering it's questions and it should help you through the tasks

---

## Roadmap

This tool is far from complete in what I want it to do, here's my upcoming list of todo's for what I want it to be.

> **Also of note, I do accept PR's**

- [ ] Make changes within Zookeeper, it's currently read only for Tracfone Specifically
- [ ] Allow you to pick other paths within Zookeeper, not just Tracfone
- [ ] Deploy multiple Repos at once, You have to do deploy tasks from within the repo you want to deploy. I want it to be smart enough that if you aren't in a repo already, it gives you a multiselect list and deploys all the repos you select to the env you select
- [ ] Stack management, it stores all the stacks, release docs, and repos you use in a local DB, but there's no process to delete those like if a stack is destroyed. The list will expand fast and it needs a way to trim it
- [ ] Auto DB management. Say, you haven't selected a specific stack within the last month, maybe the tool just removes it for you. Release docs should be gone after a couple of days
- [ ] Wrap the logger that `go-zookeeper/zk` uses so that it follows the log level set by the `-l` flag
- [ ] Add Database commands to update the DB's on stacks and staging using a mssql server driver. Think things like updating an order to be a certain type, or adding funds to an account

I also accept suggestions into the roadmap, but they won't be as considered if they don't come with an implementation plan or PR.

---

Thanks for reading
