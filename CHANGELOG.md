# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html),
and is generated by [Changie](https://github.com/miniscruff/changie).

## v0.1.0-alpha3 - 2024-05-23
### Changed
- branch delete: Report hash of deleted branch. Use this to recover the deleted branch before the next `git gc`.
- Rename 'gs complete' to 'gs completion'.
### Fixed
- repo sync: Fix deleting merged branches after a manual 'git pull'.

## v0.1.0-alpha2 - 2024-05-23
### Fixed
- branch submit: Fix default PR title and body for branches with multiple commits.
- repo sync: Delete remote tracking branches for merged PRs.

## v0.1.0-alpha1 - 2024-05-22

Initial alpha release.