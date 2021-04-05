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

// FactoryOptions defines the additional options for a Factory.
type FactoryOptions struct {
	// When set to true all SQL statements will be logged using the log package.
	LogSQL bool
}

// Factory provides functions to create new Sessions.
type Factory struct {
	pool    *sql.DB
	options *FactoryOptions
}

// NewSessionFactory creates a new Factory using connections from
// the given pool. When providing nil for the options, default options
// are used.
func NewSessionFactory(pool *sql.DB, options *FactoryOptions) *Factory {
	if options == nil {
		options = &FactoryOptions{}
	}

	return &Factory{
		pool:    pool,
		options: options,
	}
}

// Open opens a new database pool and wraps it as a SessionFactory.
func Open(driver, dsn string, options *FactoryOptions) (*Factory, error) {
	pool, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return NewSessionFactory(pool, options), nil
}

// Close closes the factory and the underlying pool
func (f *Factory) Close() {
	f.pool.Close()
}

// Session creates a new Session and binds it to the given context.
func (f *Factory) Session(ctx context.Context) (*Session, context.Context, error) {
	if s, ok := GetSession(ctx); ok {
		s.txCount++
		return s, ctx, nil
	}

	// TODO: How to provide transaction options?
	tx, err := f.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, ctx, err
	}

	s := &Session{
		options: f.options,
		tx:      tx,
		txCount: 1,
		ctx:     ctx,
	}

	return s, context.WithValue(ctx, contextKeySession, s), nil
}

// GetSession returns the Session associated with the given Context and
// a boolean flag (ok) indicating if a session has been registered with
// the given context.
func GetSession(ctx context.Context) (*Session, bool) {
	s, ok := ctx.Value(contextKeySession).(*Session)
	return s, ok
}

// MustGetSession returns the Session associated with the given Context.
// This function panics when no Session has been stored in the Context.
func MustGetSession(ctx context.Context) *Session {
	s, ok := GetSession(ctx)
	if !ok {
		panic("no session in context")
	}
	return s
}
