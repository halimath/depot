# depot

![CI Status][ci-img-url] [![Go Report Card][go-report-card-img-url]][go-report-card-url] [![Package Doc][package-doc-img-url]][package-doc-url] [![Releases][release-img-url]][release-url]

`depot` is a thin abstraction layer for accessing relational databases using Golang (technically, 
the concepts used by `depot` should be applicable to other databases as well).

`depot` is implemented to provide a more convenient API to applications while stil remaining
what I consider to be _idiomatic go_.

**`depot` is under heavy development and _not_ ready for production systems.**

# Usage

`depot` requires at least Go 1.14.

## Installation

```
$ go get github.com/halimath/depot
```

## API

The fundamental type to interact with `depot` is the `Session`. A `Session` is bound
to a `Context` and wraps a database transaction. To obtain a session, you use a 
`SessionFactory` which in turn wraps a `sql.DB`. You may create a `SessionFactory`
either by calling the `depot.Open` function passing in the same arguments you would
pass to `sql.Open`, or you create a `sql.DB` value yourself (i.e. if you want to 
configure connection pooling) and pass this to `depot.NewSessionFactory`.

Once you have a factory, call its `Session` method to create a session.

The `Session` provides methods to commit or rollback the wrapped transaction. Make
sure you call one of these methods to finish the transaction.

```go

factory := depot.Open("sqlite3", "./test.db", nil)

ctx := context.Background()
session, ctx, err := factory.Session(ctx)

err := session.Insert(depot.Into("test"), depot.Values{
    "id": 1,
    "message": "hello, world",
})
if err != nil {
    log.Fatal(err)
}

if err := session.CommitIfNoError(); err != nil {
    log.Fatal(err)
}
```

See [`acceptance-test.go`](./acceptance-test.go) for an almost complete API example.

## Code Generator

The code generator provided by `depot` can be used to generate data access types - 
called _repositories_ - for Go-`struct` types.

In order to generate a repository the struct's fields must be tagged using standard
Go tags with the key `depot`. The value is the column name to map the field to.

See the [example app](./example) for a working example.

# Open Issues

`depot` is under heavy development. Expect a lot of bugs. A list of open features 
can be found in [`TODO.md`](./TODO.md).

# License

```
Copyright 2021 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[ci-img-url]: https://github.com/halimath/depot/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/depot
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/depot
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/depot
[release-img-url]: https://img.shields.io/github/v/release/halimath/depot.svg
[release-url]: https://github.com/halimath/depot/releases