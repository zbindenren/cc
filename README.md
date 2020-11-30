[![Go Report Card](https://goreportcard.com/badge/github.com/zbindenren/cc)](https://goreportcard.com/report/github.com/zbindenren/cc)
[![Coverage Status](https://coveralls.io/repos/github/zbindenren/cc/badge.svg)](https://coveralls.io/github/zbindenren/cc)
[![Build Status](https://github.com/zbindenren/cc/workflows/build/badge.svg)](https://github.com/zbindenren/cc/actions)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/zbindenren/cc)](https://pkg.go.dev/github.com/zbindenren/cc)

# cc
A small go library to parse [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/). Instead of regular expressions,
this package uses a lexer, that functions similarly to Rob Pike's discussion about lexer design in this [talk](https://www.youtube.com/watch?v=HxaD_trXwRE).
