# Requirements

`depot` requires at **Go >= 1.14** in order to compile the runtime dependencies as well as the code generator.

# Installation

To use the module, run

```
$ go get github.com/halimath/depot
```

To also use the code generator, run

```
$ go install github.com/halimath/depot/cmd/depot
```

# Connecting to a Database

`depot` uses an abstraction similar to the one offered by the `database/sql` package. You open a connection to
a `depot.DB` which you can use to acquire a `depot.Tx` which in turn allows you to execute queries. You may 
create a `DB` either by calling the `depot.Open` function passing in the same arguments you would pass to 
`sql.Open`.

```go
db := depot.Open("sqlite3", "./test.db", depot.Options{})
```

You can also create a `sql.DB` value yourself (i.e. if you want to configure connection pooling) and pass
this to `depot.New`. Using `depot.New` with a preconfigured database is recommended for production and 
especially for performance critical code.

```go
pool := sql.Open("sqlite3", "./test.db")
// Configure pool size
db := depot.New(pool, depot.Options{})
```

## Connection Options

Both calls to `depot.Open` and `depot.New` require to pass in a `depot.Options` value which can be used to
further customize the database usage. The most important aspect is the `Dialect` which defines how the 
generated SQL will look like. `depot` provides a `DefaultDialect` which is used by default, but it is 
recommended to use a dedicated dialect for the database you are connecting to. `depot` provides dialects for
SQLite, MariaDB and PostgreSQL out of the box.

In addition, you may choose other options (i.e. logging of generated SQL). Check the code to see the available
options.

# Interacting with the Database

Once you have a `DB`, call its `BeginTx` method to begin a transaction. A `Tx` is always bound to a `Context`
in the same way that `sql.BeginTx` binds a transaction to a `Context`. There is no way to work around this
requirement: You definetly need to use contexts.

## Handling Commits and Rollbacks

The `Tx` provides methods to commit or rollback the transaction. Make sure you call one of these 
methods to finish the transaction. The `Rollback` method does nothing, if the transaction has already been 
committed. Thus, you can safely call this method using `defer`. This way, the transaction gets rolled back
in case of any error, but a successful call to `Commit` will commit the transaction if no error occured
along the way.

```go
func InsertSomething (db *depot.DB) error {
	ctx := context.Background()
	tx, ctx, err := db.BeginTx(ctx)
	defer tx.Rollback()

	// Execute queries

	return tx.Commit()
}
```

## Executing Queries

`Tx` provides an interface to issue queries to the database. You can use bare SQL strings but the 
preferred way to issue queries with `depot` is using the abstraction API that uses `depot.Values` and 
`query.Clause`s. 

`Values` are essentially a map of column names to values and are used to read and write data from/to the 
database. `Clause`s are used to dynamically build queries that pay attention to the specific SQL dialect of 
the target database engine.

```go
db := depot.Open("sqlite3", "./test.db", depot.Options{})

ctx := context.Background()
session, ctx, err := db.Session(ctx)

err := session.Insert(depot.Into("test"), depot.Values{
    "id": 1,
    "message": "hello, world",
})
if err != nil {
    log.Fatal(err)
}

if err := session.Commit(); err != nil {
    log.Fatal(err)
}
```

See [`depot_test.go`](./depot_test.go) for an almost complete API example. 


In addition, you may look at
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
the `Context`. Under the hood all of the methods use the `Tx` described above.

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