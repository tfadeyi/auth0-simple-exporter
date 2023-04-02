# Contributing

Thanks for taking the time to contribute! The following is a set of guidelines for contributing to our project.
We encourage everyone to follow them with their best judgement.


## Prerequisites

- [Go 1.19+](https://go.dev/) (minimum version required is 1.19):
    - Install on macOS with `brew install go`.
    - Install on Ubuntu with `sudo apt install golang`.
    - Install on Windows with [this link](https://go.dev/doc/install) or `choco install go`
- Git
- Make (optional, if nix is present)
- Nix (optional, if make is present)

## Setting Up Your Environment

1. Fork the repository on GitHub.
2. Clone your forked repository to your local machine.

```shell
 git clone https://github.com/tfadeyi/auth0-simple-exporter.git
```
3. Change directory to the cloned repository.

```shell
cd auth0-simple-exporter
```
4. Install dependencies.

```shell
go get ./...
```

## Making Changes

1. Create a new branch for your changes.

```shell
git checkout -b <issue number>-<branch name>
```

2. Make your changes and commit them.

```shell
git commit --signoff
```

3. Push your changes to your forked repository.

```shell
git push origin <issue number>-<branch name>
```

4. Open a pull request on GitHub from your forked repository to the original repository.

## Code Review Process

All contributions will be reviewed by the maintainers of [Project Name]. Here are a few things to keep in mind:

## How to prepare and merge a pull request

Pull requests let you tell others about changes you have pushed to a branch in the repository. They are a dedicated forum for discussing the implementation of the proposed feature or bugfix between committer and reviewer(s).
This is an essential mechanism to maintain or improve the quality of our codebase, so let's see what we look for in a pull request.

The following are all formal requirements to be met before considering the content of the proposed changes.

### Naming Conventions

For pull requests, branches, and commits, a standard naming convention will help with automatically linking the development work with a ticket management program (JIRA). Our JIRA tickets follow a format of `VC-XXXXX` found in the story/task ID of the ticket. For this reason, please follow the following naming conventions:

* Pull Requests: `<pull request title>`
* Branches: `<issue number>-<branch-name>`
* Commits: `Related <[#issue number](issue URL)>`
