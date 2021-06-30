<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [cc](#cc)
  - [Changelog CLI](#changelog-cli)
    - [Installation](#installation)
    - [Configuration](#configuration)
    - [Usage](#usage)
    - [Markdown](#markdown)
  - [Library](#library)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

[![Go Report Card](https://goreportcard.com/badge/github.com/zbindenren/cc)](https://goreportcard.com/report/github.com/zbindenren/cc)
[![Coverage Status](https://coveralls.io/repos/github/zbindenren/cc/badge.svg)](https://coveralls.io/github/zbindenren/cc)
[![Build Status](https://github.com/zbindenren/cc/workflows/build/badge.svg)](https://github.com/zbindenren/cc/actions)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/zbindenren/cc)](https://pkg.go.dev/github.com/zbindenren/cc)

# cc
A small go library to parse [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) and a cli to create a changelogs.

## Changelog CLI
The changelog cli creates and updates [CHANGELOG.md](./CHANGELOG.md) markdown files.

### Installation
You can download a precompiled binary from the releases page or install it with go:

```console
$ go install github.com/zbindenren/cc/cmd/changelog@latest
```


### Configuration
You can create a default changelog configuration `.cc.yml` with `changelog -init-config`. This results in the following configuration:

```yaml
sections:
    - type: build
      title: Build System
      hidden: true
    - type: docs
      title: Documentation
      hidden: true
    - type: feat
      title: New Features
      hidden: false
    - type: fix
      title: Bug Fixes
      hidden: false
    - type: refactor
      title: Code Refactoring
      hidden: true
    - type: test
      title: Test
      hidden: true
    - type: chore
      title: Tasks
      hidden: true
github_project_path: ""
```

Hidden sections will not show up in the resulting changelog. The default configuration creates [Gitlab](https://gitlab.com) Markdown.
If your our project is on Github, you have to add the project path to `.cc.yml`:

```yaml
github_project_path: zbindenren/cc
```

### Usage
To create a new release run:

```
$ changelog
last version: 0.2.1
next version: 0.3.0
create release 0.3.0 (press enter to continue with this version or enter version):
```

The proposed version corresponds to [Semantic Versioning](https://semver.org), but you can override the version, by entering a different one. The entered version can
not be below the current version.

The above command then performs the following tasks:

* creates or update `CHANGELOG.md` file
* stages (if necessary) and commits the changes
* create a new version tag
* and pushes everthing to remote

If you just want to see what happens, you can run `changelog -stdout`. With this option, no changes are applied to the git repository.

If you have already release tags in your project, you can create the old changelog with: `changelog -history > CHANGELOG.md`. The history command always
prints to stdout and performs no commits.

To see all available options run: `changelog -h`.

### Markdown
Commits of the form:

```
feat(compiler): add 'comments' option
```

```
feat: add a new feature
```

lead to following Markdown:

```markdown
### New Features

* **common** add a new feature ([a1f6009e](https://github.com/zbindenren/cc/commit/a1f6009e))
* **compiler** add 'comments' option  ([aecbc18b](https://github.com/zbindenren/cc/commit/aecbc18b))
```

If you add a footer with a issue reference `Closes: #1` or `Fixes: #1` like:

```
fix: a bug

Closes: #1
```

you get the following markdown:


```markdown
### Bug Fixes

* **common** fix a bug ([#1](https://github.com/zbindenren/cc/issues/#1), [a1f6009e](https://github.com/zbindenren/cc/commit/a1f6009e))
```

If you commit breaking changes:

```
feat(server-cmd)!: remove option -a
```

```
feat(client-cmd): remove option -b

BREAKING CHANGE: this change is introduced because ...
```

you get the following markdown:

```markdown
### Breaking Changes

* **server-cmd** remove option -a ([a1f6009e](https://github.com/zbindenren/cc/commit/a1f6009e))
* **client-cmd** remove option -b ([a1f6009e](https://github.com/zbindenren/cc/commit/a1f6009e))
  > this change is introduced because ...
```

An example can be found [here](./CHANGELOG.md).

## Library
Instead of regular expressions, this package uses a lexer, that functions similarly to Rob Pike's discussion about lexer
design in this [talk](https://www.youtube.com/watch?v=HxaD_trXwRE).

This library, parses a commit of the form:

```
fix: correct minor typos in code

see the issue for details

on typos fixed.

Reviewed-by: Z
Refs #133
```

into a struct:

```go
&cc.Commit{
  Header: cc.Header{
    Type: "fix",
    Scope: "",
    Description: "correct minor typos in code",
  },
  Body: "see the issue for details\n\non typos fixed.",
  Footer: cc.Footers{
    cc.Footer{
      Token: "Reviewed-by",
      Value: "Z",
    },
    cc.Footer{
      Token: "Refs",
      Value: "#133",
    },
  },
}
```
