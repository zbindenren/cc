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



