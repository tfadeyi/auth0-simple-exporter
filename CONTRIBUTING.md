# Contributing

Thanks for taking the time to contribute! The following is a set of guidelines for contributing to our project.
We encourage everyone to follow them with their best judgement.


## Prerequisites

- [Go 1.19+](https://go.dev/) (minimum version required is 1.19):
    - Install on macOS with `brew install go`.
    - Install on Ubuntu with `sudo apt install golang`.
    - Install on Windows with [this link](https://go.dev/doc/install) or `choco install go`
- Git
- Make
- Nix (optional)

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

All contributions will be reviewed by the maintainers of the project. Here are a few things to keep in mind:
  * Please fill the given Pull Request template to the best of your abilities.
  * Opening an issue before starting a work pieces improves the chances of the work being approved.
### Naming Conventions

For pull requests and branches a standard naming convention will help with automatically linking the development work with the related issue(s).
For this reason, please follow the following naming conventions:

* Branches: When creating a new branch the issue number should be added as a prefix `<issue number>-<branch-name>`
* Commits: The commit body should reference the issue `Related <[#issue number](issue URL)>`
