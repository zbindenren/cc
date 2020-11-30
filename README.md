[![Go Report Card](https://goreportcard.com/badge/github.com/zbindenren/cc)](https://goreportcard.com/report/github.com/zbindenren/cc)
[![Coverage Status](https://coveralls.io/repos/github/zbindenren/cc/badge.svg)](https://coveralls.io/github/zbindenren/cc)
[![Build Status](https://github.com/zbindenren/cc/workflows/build/badge.svg)](https://github.com/zbindenren/cc/actions)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/zbindenren/cc)](https://pkg.go.dev/github.com/zbindenren/cc)

# cc
A small go library to parse [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/). Instead of regular expressions,
this package uses a lexer, that functions similarly to Rob Pike's discussion about lexer design in this [talk](https://www.youtube.com/watch?v=HxaD_trXwRE).

This library, creates whit a commit of the form:

```
fix: correct minor typos in code

see the issue for details

on typos fixed.

Reviewed-by: Z
Refs #133
```

a struct like following:

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
