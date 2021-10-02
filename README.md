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

`depot` requires at least Go 1.14.

## Installation

To use the module, run

```
$ go get github.com/halimath/depot
```

To also use the code generator, run

```
$ go install github.com/halimath/depot/cmd/depot
```

## API

The fundamental type to interact with `depot` is the `Session`. A `Session` is bound to a `Context` and  
wraps a database transaction. To obtain a session, you use a `SessionFactory` which in turn wraps a 
`sql.DB`. You may create a `SessionFactory` either by calling the `depot.Open` function passing in the same 
arguments you would pass to `sql.Open`, or you create a `sql.DB` value yourself (i.e. if you want to 
configure connection pooling) and pass this to `depot.NewSessionFactory`.

Once you have a factory, call its `Session` method to create a session.

The `Session` provides methods to commit or rollback the wrapped transaction. Make sure you call one of these 
methods to finish the transaction.

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

See [`depot_test.go`](./depot_test.go) for an almost complete API example. In addition, you may look at
the [`example`](./example) which uses the code generator.

## Code Generator

The code generator provided by `depot` can be used to generate data access types - 
called _repositories_ - for Go-`struct` types.

In order to generate a repository the struct's fields must be tagged using standard
Go tags with the key `depot`. The value is the column name to map the field to.

The following struct shows the generation process:

```go
type (
	// Message demonstrates a persistent struct showing several mapped fields.
	Message struct {
		ID         string    `depot:"id,id"`
		Text       string    `depot:"text"`
		OrderIndex int       `depot:"order_index"`
		Length     float32   `depot:"len"`
		Attachment []byte    `depot:"attachment"`
		Created    time.Time `depot:"created"`
	}
)
```

The struct `Message` should be mapped to a table `messages` with each struct field being mapped to
a column. Every field tagged with a `depot` tag will be part of the mapping. All other fields are
ignored. 

The column name is given as the tag's value. The field `ID` being used as the primary key is further
marked with `id` directive separated by a comma.

To generate a repository for this model, invoke `depot` with the following command line:

```
$ depot generate-repo --table=messages --repo-package repo --out ../repo/gen-messagerepo.go models.go Message
```

You can also use `go:generate` to do the same thing. Simply place a comment like

```go
//go:generate depot generate-repo --table=messages --repo-package repo --out ../repo/gen-messagerepo.go $GOFILE Message
```

in the source file containing `Message`. You can place multiple comments of that type in a single source file containing
multiple model structs.

The generated repository provides the following methods:

```go
type MessageRepo struct {
	factory *depot.Factory
}

func NewMessageRepo(factory *depot.Factory) *MessageRepo
func (r *MessageRepo) Begin(ctx context.Context) (context.Context, error)
func (r *MessageRepo) Commit(ctx context.Context) error
func (r *MessageRepo) Rollback(ctx context.Context) error
func (r *MessageRepo) fromValues(vals depot.Values) (*models.Message, error)
func (r *MessageRepo) find(ctx context.Context, clauses ...depot.Clause) ([]*models.Message, error)
func (r *MessageRepo) count(ctx context.Context, clauses ...depot.Clause) (int, error)
func (r *MessageRepo) LoadByID(ctx context.Context, ID string) (*models.Message, error)
func (r *MessageRepo) toValues(entity *models.Message) depot.Values
func (r *MessageRepo) Insert(ctx context.Context, entity *models.Message) error
func (r *MessageRepo) delete(ctx context.Context, clauses ...depot.Clause) error
func (r *MessageRepo) Update(ctx context.Context, entity *models.Message) error
func (r *MessageRepo) DeleteByID(ctx context.Context, ID string) error
func (r *MessageRepo) Delete(ctx context.Context, entity *models.Message) error
```

_Note that no public method uses a parameter or return type exported from `depot`. This means that a generated
repository's public interface does not expose the use of `depot`. This allows the interface to be used in
architecture styles (such as _clean architecture_ or _ports and adapters_) that require the business types
to have no dependency to a persistence framework._

Use `Begin`, `Commit` and `Rollback` to control the transaction scope. The transaction is stored as part of 
the `Context`. Under the hood all of the methods use the `Session` described above.

`find` and `count` are methods that can be used by custom finder methods. They execute `select` queries for 
the entity. `LoadByID` uses `find` to load a single message by `ID`.

The mutation methods all handle single instances of `Message`. `delete` is provided similar to `find` to do 
batch deletes.

Note that all of these methods contain simple to read and debug go code with no reflection being used at all.
The code is idiomatic and in most cases looks like being written by a human go programmer.

You can easily extend the repo with custom finder methods by writing a custom method being part of 
`MessageRepo` to another non-generated file and use `find` to do the actual work. Here is an example for a
method to load `Message`s based on the value of the `Text` field.

```go
func (r *MessageRepo) FindByText(ctx context.Context, text string) ([]*models.Message, error) {
	return r.find(ctx, depot.Where("text", text))
}
```

### `null` values

If a mapped field should support SQL `null` values, you have to add the `nullable` directive
to the field's mapping tag:

```go
type Entity struct {
	// ...
	Foo *string `depot:"foo,nullable"`
}
```

You can either use a pointer type as in the example above or the plain type (`string` in this case).
Using a pointer is recommended, as SQL `null` values will be represented as `nil`. When using the
plain type, `null` is represented with the value's default type (`""` in this case) which only
works when reading `null` from the datase. If you wish to `insert` or `update` a `null` value
you are required to use a pointer type.

### List of directives

The following table lists all supported directives for field mappings.

Directive | Used for | Example | Description
-- | -- | -- | --
`id` | Mark a field as the entity's ID. | `ID string "depot:\"id,id\""` | Only a single field may be tagged with `id`. If one is given, the generated repo will contain the methods `LoadByID` and `DeleteByID` which are not generated when no ID is declared.
`nullable` | Mark a field as being able to store a `null` value. | `Message *string "depot:\"msg,nullable\""` | See the section above for `null` values.

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