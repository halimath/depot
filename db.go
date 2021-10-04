// Copyright 2021 Alexander Metzner.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package depot

import (
	"context"
	"database/sql"
)

// contextKeySessionType defines the type used to store a Session in a Context.
type contextKeySessionType string

// contextKeySession defines the value used as the key to store a Session in a Context.
const contextKeySession contextKeySessionType = "session"

// Options defines the options for a DB.
type Options struct {
	// Dialect defines the SQL dialect to generate. If not set a default dialect will be used.
	Dialect Dialect

	// When set to true all SQL statements will be logged using the log package.
	LogSQL bool
}

// DB provides the primary interface to interact with the persistence framework. It provides functions to
// create new Sessions which encapsulate a database transaction.
type DB struct {
	pool    *sql.DB
	options Options
}

// New creates a new DB using connections from the given pool. Options may be empty in which case defaults
// are used (which might not work with the database connecting to). It is recommended to at least provide a
// matching Dialect.
func New(pool *sql.DB, options Options) *DB {
	if options.Dialect == nil {
		options.Dialect = &DefaultDialect{}
	}

	return &DB{
		pool:    pool,
		options: options,
	}
}

// Open opens a new database pool and wraps it in a DB. This function resembles sql.Open (which is called)
// from this function and uses defaults to connect to the database. For performance-critical code it is
// recommended to use New with a preconstructed database pool.
func Open(driver, dsn string, options Options) (*DB, error) {
	pool, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return New(pool, options), nil
}

// Close closes the depot and the underlying pool.
func (f *DB) Close() {
	f.pool.Close()
}

// BeginTx creates a begins a new transaction and binds it to ctx. If ctx already contains a transaction this
// it is returned instead with it's txCount incremented.
func (f *DB) BeginTx(ctx context.Context) (*Tx, context.Context, error) {
	// TODO: Add support to call this function in parallel.

	if s, ok := GetTx(ctx); ok {
		s.txCount++
		return s, ctx, nil
	}

	// TODO: How to provide transaction options?
	tx, err := f.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, ctx, err
	}

	s := &Tx{
		options: &f.options,
		tx:      tx,
		txCount: 1,
		ctx:     ctx,
	}

	return s, context.WithValue(ctx, contextKeySession, s), nil
}

// GetTx returns the transaction associated with the given Context and a boolean flag (ok) indicating if a
// transaction has been registered with the given context.
func GetTx(ctx context.Context) (*Tx, bool) {
	s, ok := ctx.Value(contextKeySession).(*Tx)
	return s, ok
}

// MustGetTx returns the transaction associated with the given Context. This function panics when no
// transaction has been stored in the Context.
func MustGetTx(ctx context.Context) *Tx {
	s, ok := GetTx(ctx)
	if !ok {
		panic("no tx in context")
	}
	return s
}
