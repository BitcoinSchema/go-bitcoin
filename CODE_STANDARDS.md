# Code Standards

This project uses the following code standards and specifications from:
- [effective go](https://golang.org/doc/effective_go.html)
- [go benchmarks](https://golang.org/pkg/testing/#hdr-Benchmarks)
- [go examples](https://golang.org/pkg/testing/#hdr-Examples)
- [go tests](https://golang.org/pkg/testing/)
- [godoc](https://godoc.org/golang.org/x/tools/cmd/godoc)
- [gofmt](https://golang.org/cmd/gofmt/)
- [golangci-lint](https://golangci-lint.run/usage/quick-start/)
- [golint](https://github.com/golang/lint)
- [report card](https://goreportcard.com/)
- [vet](https://golang.org/cmd/vet/)

### *effective go* standards
View the [effective go](https://golang.org/doc/effective_go.html) standards documentation.

### *golint* specifications
The package [golint](https://github.com/golang/lint) differs from [gofmt](https://golang.org/cmd/gofmt/). The package [gofmt](https://golang.org/cmd/gofmt/) formats Go source code, whereas [golint](https://github.com/golang/lint) prints out style mistakes. The package [golint](https://github.com/golang/lint) differs from [vet](https://golang.org/cmd/vet/). The package [vet](https://golang.org/cmd/vet/) is concerned with correctness, whereas [golint](https://github.com/golang/lint) is concerned with coding style. The package [golint](https://github.com/golang/lint) is in use at Google, and it seeks to match the accepted style of the open source [Go project](https://golang.org/).

How to install [golint](https://github.com/golang/lint):
```shell script
go get -u golang.org/x/lint/golint
cd ../go-bitcoin
golint
```

### *golangci-lint* specifications
The package [golangci-lint](https://golangci-lint.run/usage/quick-start) runs several linters in one package/cmd.

How to install [golangci-lint](https://golangci-lint.run/):
```shell script
brew install golangci-lint
```

### *go vet* specifications
[Vet](https://golang.org/cmd/vet/) examines Go source code and reports suspicious constructs. [Vet](https://golang.org/cmd/vet/) uses heuristics that do not guarantee all reports are genuine problems, but it can find errors not caught by the compilers.

How to run [vet](https://golang.org/cmd/vet/):
```shell script
cd ../go-bitcoin
go vet -v
```

### *godoc* specifications
All code is written with documentation in mind. Follow the best practices with naming, examples and function descriptions.