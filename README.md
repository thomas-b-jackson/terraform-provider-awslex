# Terraform Provider AWS Lex

Provider for building aws lexv2 bots via terraform

## Requirements

- wsl or mac
-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Using the provider

See [./examples](./examples)

## Developing the Provider

Development steps:
1. update `VERSION` in `GNUMakefile` with the desired version tag
2. set aws credentials (as env variables or via `.aws/credentials` file)
3. make changes to provider sources
4. build and install the provider locally by running:
   `make install`
5. test the provider against examples in [./examples](./examples) as:
   1. reference the `localhost/va/awslex` version of the provider
   2. remove `.terraform.lock.hcl` between `make install` iterations
6. re-run integrations test(s) as:
   `make test`

## Release

Create a release using GoReleaser. 

**Note:** steps are adapted from [these instructions](https://www.terraform.io/docs/registry/providers/publishing.html#using-goreleaser-locally)

Setup Steps:
* Install GoReleaser
* Obtain fingerprint of GPG private key for signing (key currently controlled by Tom Jackson)
  * fingerprint is 40 chars and is obtained by running this command:
    `gpg --list-secret-keys --keyid-format=long`
* Obtain Personal Access Token for repo (token currently controlled by Tom Jackson)

Release Steps:
* Commit changes locally
* Set GITHUB_TOKEN to a Personal Access Token
* Set GPG_FINGERPRINT to fingerprint
* Tag your release commit to match version in GNUmakefile, e.g.:
  `git tag 0.2.0`
* Build, sign, and upload your release with:
  `goreleaser release --rm-dist`
* Re-run terraform init against the release in the registry (to make sure it has sync'd from github)
* Test the released provider in a pipeline
* Push commit and do pull request
 
## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.