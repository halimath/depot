# depot

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]

`depot` is a thin abstraction layer for accessing relational databases using Golang. In addition, `depot`
provides a code generator which generates object-relational mappings (ORM) and repository types that easily
map Go types (most notably `struct`s) to database tables and vice versa.

`depot` is implemented to provide a more convenient API to applications while stil remaining what I consider
to be _idiomatic go_.

**`depot` is under heavy development and _not_ ready for production systems.**

# Usage

See the [usage guide](./docs/usage.md) for a detailed description.

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