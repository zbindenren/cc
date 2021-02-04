## 0.4.1 (2021-02-04)


### Bug Fixes

* **changelog**: create changelog file in the top level git directory instead of current directory ([#2](https://github.com/zbindenren/cc/issues/2), [ccfe76e1](https://github.com/zbindenren/cc/commit/ccfe76e1))



## 0.4.0 (2021-01-13)


### Bug Fixes

* **changelog**: include first release (tag) to history ([f37a1f1d](https://github.com/zbindenren/cc/commit/f37a1f1d))
* **changelog**: issue link URL in github markdown ([4df9cbb8](https://github.com/zbindenren/cc/commit/4df9cbb8))
  > Before this fix, the issue link pointed incorrectly to `https://github.com/org/project/issues/#<issueNR>`.
  > Now the URL points to `https://github.com/org/project/issues/<issueNR>`
  > (without '#').


### New Features

* **changelog**: add `-num` option ([638ea456](https://github.com/zbindenren/cc/commit/638ea456))
  > With this option, it is possible to limit the number of releases (tags)
  > that are included in history output.



## 0.3.1 (2021-01-12)


### Bug Fixes

* **changelog**: command does not fail when there was no commit since last tag ([#1](https://github.com/zbindenren/cc/issues/#1), [adbe8e1a](https://github.com/zbindenren/cc/commit/adbe8e1a))
  > Release failed with a cryptic git error, when there was no commit since last
  > tag. Now command exits with `no commits since last tag` error message.



## 0.3.0 (2021-01-05)


### Bug Fixes

* **changelog**: github issue and commit markdown URLs ([cdeb3531](https://github.com/zbindenren/cc/commit/cdeb3531))
* **changelog**: stage buildinfo go files ([e6f7a8c6](https://github.com/zbindenren/cc/commit/e6f7a8c6))
* **changelog**: verify that new version is greater than current ([85b778f4](https://github.com/zbindenren/cc/commit/85b778f4))
* **changelog**: when overriding release version corresponding tag is also overridden ([62c99ab1](https://github.com/zbindenren/cc/commit/62c99ab1))


### New Features

* **changelog**: add version flag `-v` to show version information ([4b5090c2](https://github.com/zbindenren/cc/commit/4b5090c2))



## 0.2.1 (2021-01-04)


### Bug Fixes

* **changelog**: add `BREAKING CHANGE` footer token value to changelog ([2e070cd5](https://github.com/zbindenren/cc/commit/2e070cd5))
  > Before this change, footer token values for `BREAKING CHANGE` and
  > `BREAKING_CHANGE` were ignored.
* **changelog**: version is no longer prefixed with `v` for `-history` flag ([4a8804f1](https://github.com/zbindenren/cc/commit/4a8804f1))
* **common**: documentation for BreakingMessage() method. ([9cd0522e](https://github.com/zbindenren/cc/commit/9cd0522e))



## 0.2.0 (2021-01-04)


### Bug Fixes

* **changelog**: raise no error when output is set to stdout and `CHANGELOG.md` does not exist ([ba54c827](https://github.com/zbindenren/cc/commit/ba54c827))
* **changelog**: stage changelog file if necessary ([0bdd1106](https://github.com/zbindenren/cc/commit/0bdd1106))
* **changelog**: typo in error message ([30cc4d88](https://github.com/zbindenren/cc/commit/30cc4d88))
* **changelog**: when overriding release version created tag was not overridden ([e185c32b](https://github.com/zbindenren/cc/commit/e185c32b))


### New Features

* **changelog**: add github markdown support ([a1f6009e](https://github.com/zbindenren/cc/commit/a1f6009e))
* **common**: add `config.Read` method to unmarshal config from `io.Reader` ([aecbc18b](https://github.com/zbindenren/cc/commit/aecbc18b))
* **common**: add command to create a markdown changelog file ([b874c814](https://github.com/zbindenren/cc/commit/b874c814))



