package depot

import (
	"context"
	"database/sql"
)

// contextKeySessionType defines the type used to store a Session in a Context.
type contextKeySessionType string

// contextKeySession defines the value used as the key to store a Session in a Context.
const contextKeySession contextKeySessionType = "session"

// Factory provides functions to create new Sessions.
type Factory struct {
	pool *sql.DB
}

// NewSessionFactory creates a new Factory using connections from
// the given pool.
func NewSessionFactory(pool *sql.DB) *Factory {
	return &Factory{
		pool: pool,
	}
}

// Open opens a new database pool and wraps it as a SessionFactory.
func Open(driver, dsn string) (*Factory, error) {
	pool, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return NewSessionFactory(pool), nil
}

// Close closes the factory and the underlying pool
func (f *Factory) Close() {
	f.pool.Close()
}

// Session creates a new Session and binds it to the given context.
func (f *Factory) Session(ctx context.Context) (*Session, context.Context, error) {
	// TODO: How to provide options?
	tx, err := f.pool.BeginTx(ctx, nil)
	if err != nil {
		return nil, ctx, err
	}

	s := &Session{
		factory: f,
		tx:      tx,
		ctx:     ctx,
	}

	return s, context.WithValue(ctx, contextKeySession, s), nil
}

// GetSession returns the Session associated with the given Context.
// This function panics when no Session has been stored in the Context.
func GetSession(ctx context.Context) *Session {
	s, ok := ctx.Value(contextKeySession).(*Session)
	if !ok {
		panic("no session in context")
	}
	return s
}
