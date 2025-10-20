# ₿ go-bitcoin
> A library for working with Bitcoin (BSV) transactions, addresses, keys, encryption, and more.

<table>
  <thead>
    <tr>
      <th>CI&nbsp;/&nbsp;CD</th>
      <th>Quality&nbsp;&amp;&nbsp;Security</th>
      <th>Docs&nbsp;&amp;&nbsp;Meta</th>
      <th>Community</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td valign="top" align="left">
        <a href="https://github.com/BitcoinSchema/go-bitcoin/releases">
          <img src="https://img.shields.io/github/release-pre/BitcoinSchema/go-bitcoin?logo=github&style=flat" alt="Latest Release">
        </a><br/>
        <a href="https://github.com/BitcoinSchema/go-bitcoin/actions">
          <img src="https://img.shields.io/github/actions/workflow/status/BitcoinSchema/go-bitcoin/fortress.yml?branch=master&logo=github&style=flat" alt="Build Status">
        </a><br/>
		<a href="https://github.com/BitcoinSchema/go-bitcoin/actions">
          <img src="https://github.com/BitcoinSchema/go-bitcoin/actions/workflows/codeql-analysis.yml/badge.svg?style=flat" alt="CodeQL">
        </a><br/>
        <a href="https://github.com/BitcoinSchema/go-bitcoin/commits/master">
		  <img src="https://img.shields.io/github/last-commit/BitcoinSchema/go-bitcoin?style=flat&logo=clockify&logoColor=white" alt="Last commit">
		</a>
      </td>
      <td valign="top" align="left">
        <a href="https://goreportcard.com/report/github.com/BitcoinSchema/go-bitcoin">
          <img src="https://goreportcard.com/badge/github.com/BitcoinSchema/go-bitcoin?style=flat" alt="Go Report Card">
        </a><br/>
		<a href="https://codecov.io/gh/BitcoinSchema/go-bitcoin">
          <img src="https://codecov.io/gh/BitcoinSchema/go-bitcoin/branch/master/graph/badge.svg?style=flat" alt="Code Coverage">
        </a><br/>
		<a href="https://scorecard.dev/viewer/?uri=github.com/BitcoinSchema/go-bitcoin">
          <img src="https://api.scorecard.dev/projects/github.com/BitcoinSchema/go-bitcoin/badge?logo=springsecurity&logoColor=white" alt="OpenSSF Scorecard">
        </a><br/>
		<a href=".github/SECURITY.md">
          <img src="https://img.shields.io/badge/security-policy-blue?style=flat&logo=springsecurity&logoColor=white" alt="Security policy">
        </a>
      </td>
      <td valign="top" align="left">
        <a href="https://golang.org/">
          <img src="https://img.shields.io/github/go-mod/go-version/BitcoinSchema/go-bitcoin?style=flat" alt="Go version">
        </a><br/>
        <a href="https://pkg.go.dev/github.com/BitcoinSchema/go-bitcoin?tab=doc">
          <img src="https://pkg.go.dev/badge/github.com/BitcoinSchema/go-bitcoin.svg?style=flat" alt="Go docs">
        </a><br/>
        <a href=".github/AGENTS.md">
          <img src="https://img.shields.io/badge/AGENTS.md-found-40b814?style=flat&logo=openai" alt="AGENTS.md rules">
        </a><br/>
        <a href="https://github.com/mrz1836/mage-x">
          <img src="https://img.shields.io/badge/Mage-supported-brightgreen?style=flat&logo=go&logoColor=white" alt="MAGE-X Supported">
        </a><br/>
		<a href=".github/dependabot.yml">
          <img src="https://img.shields.io/badge/dependencies-automatic-blue?logo=dependabot&style=flat" alt="Dependabot">
        </a>
      </td>
      <td valign="top" align="left">
        <a href="https://github.com/BitcoinSchema/go-bitcoin/graphs/contributors">
          <img src="https://img.shields.io/github/contributors/BitcoinSchema/go-bitcoin?style=flat&logo=contentful&logoColor=white" alt="Contributors">
        </a><br/>
        <a href="https://github.com/sponsors/BitcoinSchema">
          <img src="https://img.shields.io/badge/sponsor-BitcoinSchema-181717.svg?logo=github&style=flat" alt="Sponsor">
        </a><br/>
        <a href="https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-bitcoin&utm_term=go-bitcoin&utm_content=go-bitcoin">
          <img src="https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat" alt="Donate Bitcoin">
        </a>
      </td>
    </tr>
  </tbody>
</table>

<br/>

## 🗂️ Table of Contents

- [Installation](#-installation)
- [Documentation](#-documentation)
- [Examples & Tests](#-examples--tests)
- [Benchmarks](#-benchmarks)
- [Code Standards](#-code-standards)
- [Maintainers](#-maintainers)
- [Contributing](#-contributing)
- [License](#-license)

<br/>

## 📦 Installation

**go-bitcoin** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).

```shell script
go get -u github.com/bitcoinschema/go-bitcoin/v2
```

<br/>

## 📚 Documentation

View the generated [documentation](https://pkg.go.dev/github.com/bitcoinschema/go-bitcoin)

### Features

- **Addresses**
  - [Address from PrivateKey (bec.PrivateKey)](address.go)
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
  - [Create WIF](private_key.go)
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
- [libsv/go-bk](https://github.com/libsv/go-bk)
- [libsv/go-bt](https://github.com/libsv/go-bt)
</details>

<details>
<summary><strong><code>Development Setup (Getting Started)</code></strong></summary>
<br/>

Install [MAGE-X](https://github.com/mrz1836/mage-x) build tool for development:

```bash
# Install MAGE-X for development and building
go install github.com/mrz1836/mage-x/cmd/magex@latest
magex update:install
```
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) for streamlined binary and library deployment to GitHub. To get started, install it via:

```bash
brew install goreleaser
```

The release process is defined in the [.goreleaser.yml](.goreleaser.yml) configuration file.

Then create and push a new Git tag using:

```bash
magex version:bump bump=patch push
```

This process ensures consistent, repeatable releases with properly versioned artifacts and citation metadata.

</details>

<details>
<summary><strong><code>Build Commands</code></strong></summary>
<br/>

View all build commands

```bash script
magex help
```

</details>

<details>
<summary><strong><code>GitHub Workflows</code></strong></summary>
<br/>


### 🎛️ The Workflow Control Center

All GitHub Actions workflows in this repository are powered by configuration files: [**.env.base**](.github/.env.base) (default configuration) and optionally **.env.custom** (project-specific overrides) – your one-stop shop for tweaking CI/CD behavior without touching a single YAML file! 🎯

**Configuration Files:**
- **[.env.base](.github/.env.base)** – Default configuration that works for most Go projects
- **[.env.custom](.github/.env.custom)** – Optional project-specific overrides

This magical file controls everything from:
- **🚀 Go version matrix** (test on multiple versions or just one)
- **🏃 Runner selection** (Ubuntu or macOS, your wallet decides)
- **🔬 Feature toggles** (coverage, fuzzing, linting, race detection, benchmarks)
- **🛡️ Security tool versions** (gitleaks, nancy, govulncheck)
- **🤖 Auto-merge behaviors** (how aggressive should the bots be?)
- **🏷️ PR management rules** (size labels, auto-assignment, welcome messages)

> **Pro tip:** Want to disable code coverage? Just add `ENABLE_CODE_COVERAGE=false` to your .env.custom to override the default in .env.base and push. No YAML archaeology required!

<br/>

| Workflow Name                                                                      | Description                                                                                                            |
|------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------|
| [auto-merge-on-approval.yml](.github/workflows/auto-merge-on-approval.yml)         | Automatically merges PRs after approval and all required checks, following strict rules.                               |
| [codeql-analysis.yml](.github/workflows/codeql-analysis.yml)                       | Analyzes code for security vulnerabilities using [GitHub CodeQL](https://codeql.github.com/).                          |
| [dependabot-auto-merge.yml](.github/workflows/dependabot-auto-merge.yml)           | Automatically merges [Dependabot](https://github.com/dependabot) PRs that meet all requirements.                       |
| [fortress.yml](.github/workflows/fortress.yml)                                     | Runs the GoFortress security and testing workflow, including linting, testing, releasing, and vulnerability checks.    |
| [pull-request-management.yml](.github/workflows/pull-request-management.yml)       | Labels PRs by branch prefix, assigns a default user if none is assigned, and welcomes new contributors with a comment. |
| [scorecard.yml](.github/workflows/scorecard.yml)                                   | Runs [OpenSSF](https://openssf.org/) Scorecard to assess supply chain security.                                        |
| [stale.yml](.github/workflows/stale-check.yml)                                     | Warns about (and optionally closes) inactive issues and PRs on a schedule or manual trigger.                           |
| [sync-labels.yml](.github/workflows/sync-labels.yml)                               | Keeps GitHub labels in sync with the declarative manifest at [`.github/labels.yml`](./.github/labels.yml).             |

</details>

<details>
<summary><strong><code>Updating Dependencies</code></strong></summary>
<br/>

To update all dependencies (Go modules, linters, and related tools), run:

```bash
magex deps:update
```

This command ensures all dependencies are brought up to date in a single step, including Go modules and any managed tools. It is the recommended way to keep your development environment and CI in sync with the latest versions.

</details>

<br/>

## 🧪 Examples & Tests
All unit tests and [examples](examples) run via [GitHub Actions](https://github.com/BitcoinSchema/go-bitcoin/actions) and
uses [Go version 1.23.x](https://golang.org/doc/go1.23). View the [configuration file](.github/workflows/fortress.yml) for more details.

Run all tests (fast):

```bash script
magex test
```

Run all tests with race detector (slower):
```bash script
magex test:race
```

<br/>

## ⚡ Benchmarks
Run the Go benchmarks:

```bash script
magex bench
```

<br/>

## 🛠️ Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## 👥 Maintainers

| [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) | [<img src="https://github.com/rohenaz.png" height="50" alt="MrZ" />](https://github.com/rohenaz) |
|:------------------------------------------------------------------------------------------------:|:------------------------------------------------------------------------------------------------:|
|                                [MrZ](https://github.com/mrz1836)                                 |                              [Satchmo](https://github.com/rohenaz)                               |

<br/>

## 🤝 Contributing

View the [contributing guidelines](.github/CONTRIBUTING.md) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?

All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/BitcoinSchema) :clap:
or by making a [**bitcoin donation**](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-bitcoin&utm_term=go-bitcoin&utm_content=go-bitcoin) to ensure this journey continues indefinitely! :rocket:

[![Stars](https://img.shields.io/github/stars/BitcoinSchema/go-bitcoin?label=Please%20like%20us&style=social)](https://github.com/BitcoinSchema/go-bitcoin/stargazers)

<br/>

## 📝 License

[![License](https://img.shields.io/github/license/BitcoinSchema/go-bitcoin.svg?style=flat)](LICENSE)
