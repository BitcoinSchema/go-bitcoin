# go-bitcoin
> A library for working with Bitcoin (BSV) transactions, addresses, keys, encryption, and more.

[![Release](https://img.shields.io/github/release-pre/BitcoinSchema/go-bitcoin.svg?logo=github&style=flat&v=4)](https://github.com/BitcoinSchema/go-bitcoin/releases)
[![Build Status](https://img.shields.io/github/workflow/status/BitcoinSchema/go-bitcoin/run-go-tests?logo=github&v=3)](https://github.com/BitcoinSchema/go-bitcoin/actions)
[![Report](https://goreportcard.com/badge/github.com/BitcoinSchema/go-bitcoin?style=flat&v=4)](https://goreportcard.com/report/github.com/BitcoinSchema/go-bitcoin)
[![codecov](https://codecov.io/gh/BitcoinSchema/go-bitcoin/branch/master/graph/badge.svg?v=4)](https://codecov.io/gh/BitcoinSchema/go-bitcoin)
[![Go](https://img.shields.io/github/go-mod/go-version/BitcoinSchema/go-bitcoin?v=4)](https://golang.org/)
[![Sponsor](https://img.shields.io/badge/sponsor-BitcoinSchema-181717.svg?logo=github&style=flat&v=4)](https://github.com/sponsors/BitcoinSchema)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=4)](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-bitcoin&utm_term=go-bitcoin&utm_content=go-bitcoin)

<br/>

## Table of Contents

- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

<br/>

## Installation

**go-bitcoin** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).

```shell script
go get -u github.com/bitcoinschema/go-bitcoin
```

<br/>

## Documentation

View the generated [documentation](https://pkg.go.dev/github.com/bitcoinschema/go-bitcoin)

[![GoDoc](https://godoc.org/github.com/bitcoinschema/go-bitcoin?status.svg&style=flat)](https://pkg.go.dev/github.com/bitcoinschema/go-bitcoin)

### Features

- **Addresses**
  - [Address from PrivateKey (bsvec.PrivateKey)](address.go)
  - [Address from Script](address.go)
- **Encryption**
  - [Encrypt With Private Key](encryption.go)
  - [Decrypt With Private Key](encryption.go)
  - [Encrypt Shared](encryption.go)
- **HD Keys** _(Master / xPub)_
  - [Generate HD Keys](hd_key.go)
  - [Generate HD Key from string](hd_key.go)
  - [Get HD Key by Path](hd_key.go)
  - [Get PrivateKey by Path](hd_key.go)
  - [Get HD Child Key](hd_key.go)
  - [Get Addresses from HD Key](hd_key.go)
  - [Get XPub from HD Key](hd_key.go)
  - [Get HD Key from XPub](hd_key.go)
  - [Get PublicKeys for Path](hd_key.go)
  - [Get Addresses for Path](hd_key.go)
- **PubKeys**
  - [Create PubKey from PrivateKey](pubkey.go)
  - [PubKey from String](pubkey.go)
- **Private Keys**
  - [Create PrivateKey](private_key.go)
  - [PrivateKey (string) to Address (string)](address.go)
  - [PrivateKey from string](private_key.go)
  - [Generate Shared Keypair](private_key.go)
  - [Get Private and Public keys](private_key.go)
  - [WIF to PrivateKey](private_key.go)
  - [PrivateKey to WIF](private_key.go)
- **Scripts**
  - [Script from Address](script.go)
- **Signatures**
  - [Sign](sign.go) & [Verify a Bitcoin Message](verify.go)
  - [Verify a DER Signature](verify.go)
- **Transactions**
  - [Calculate Fee](transaction.go)
  - [Create Tx](transaction.go)
  - [Create Tx using WIF](transaction.go)
  - [Create Tx with Change](transaction.go)
  - [Tx from Hex](transaction.go)

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

- [bitcoinsv/bsvd](https://github.com/bitcoinsv/bsvd)
- [bitcoinsv/bsvutil](https://github.com/bitcoinsv/bsvutil)
- [libsv/go-bt](https://github.com/libsv/go-bt)
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed via: `brew install goreleaser`.

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production.

</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands

```shell script
make help
```

List of all current commands:

```text
all                  Runs multiple commands
clean                Remove previous builds and any test cache data
clean-mods           Remove all the Go mod cache
coverage             Shows the test coverage
godocs               Sync the latest tag with GoDocs
help                 Show this help message
install              Install the application
install-go           Install the application (Using Native Go)
lint                 Run the golangci-lint application (install if not found)
release              Full production release (creates release in Github)
release              Runs common.release then runs godocs
release-snap         Test the full release (build binaries)
release-test         Full production test release (everything except deploy)
replace-version      Replaces the version in HTML/JS (pre-deploy)
tag                  Generate a new tag and push (tag version=0.0.0)
tag-remove           Remove a tag if found (tag-remove version=0.0.0)
tag-update           Update an existing tag to current commit (tag-update version=0.0.0)
test                 Runs vet, lint and ALL tests
test-ci              Runs all tests via CI (exports coverage)
test-ci-no-race      Runs all tests via CI (no race) (exports coverage)
test-ci-short        Runs unit tests via CI (exports coverage)
test-short           Runs vet, lint and tests (excludes integration tests)
uninstall            Uninstall the application (and remove files)
update-linter        Update the golangci-lint package (macOS only)
vet                  Run the Go vet application
```

</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [Github Actions](https://github.com/BitcoinSchema/go-bitcoin/actions) and
uses [Go version 1.14.x, 1.15.x and 1.16.x](https://golang.org/doc/go1.16). View the [configuration file](.github/workflows/run-tests.yml).

Run all tests (including integration tests)

```shell script
make test
```

Run tests (excluding integration tests)

```shell script
make test-short
```

<br/>

## Benchmarks

Run the Go benchmarks:

```shell script
make bench
```

<br/>

## Code Standards

Read more about this Go project's [code standards](CODE_STANDARDS.md).

<br/>

## Usage

Checkout all the [examples](examples)!

<br/>

## Maintainers

| [<img src="https://github.com/rohenaz.png" height="50" alt="MrZ" />](https://github.com/rohenaz) | [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
| :----------------------------------------------------------------------------------------------: | :----------------------------------------------------------------------------------------------: |
|                              [Satchmo](https://github.com/rohenaz)                               |                                [MrZ](https://github.com/mrz1836)                                 |

<br/>

## Contributing

View the [contributing guidelines](CONTRIBUTING.md) and follow the [code of conduct](CODE_OF_CONDUCT.md).

### How can I help?

All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/BitcoinSchema) :clap:
or by making a [**bitcoin donation**](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-bitcoin&utm_term=go-bitcoin&utm_content=go-bitcoin) to ensure this journey continues indefinitely! :rocket:

<br/>

## License

![License](https://img.shields.io/github/license/BitcoinSchema/go-bitcoin.svg?style=flat&v=4)
