package depot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrNoResult is returned when queries execpted to match (at least) on row
	// match no row at all.
	ErrNoResult = errors.New("no result")
	// ErrMarkedForRollback is returned when trying to commit a Session which is
	// already marked for rollback only.
	ErrMarkedForRollback = errors.New("session has been marked for rollback")
)

// Values contains the persistent column values for an entity either after reading
// the values from the database to re-create the entity value or to persist the
// entity's values in the database (either for insertion or update).
type Values map[string]interface{}

// A session defines an interaction session with the database.
// A session uses a single transaction and is bound to a single
// Context. A session provides an abstract interface built around
// Values and Clauses.
type Session struct {
	factory *Factory
	tx      *sql.Tx
	ctx     context.Context
	err     error
}

// Commit commits the session's transaction and returns an error
// if the commit fails.
func (s *Session) Commit() error {
	if s.err != nil {
		return ErrMarkedForRollback
	}
	return s.tx.Commit()
}

// Rollback rolls the session's transaction back and returns any
// error raised during the rollback.
func (s *Session) Rollback() error {
	return s.tx.Rollback()
}

// Error marks the transaction as failed so it cannot be committed
// later on.
// Calling Error with a nil error clears the error state of the transaction.
func (s *Session) Error(err error) {
	s.err = err
}

// CommitIfNoError tries to commit the transaction but performs a
// rollback in case an error has been logged before.
func (s *Session) CommitIfNoError() error {
	if s.err != nil {
		return s.Rollback()
	}
	return s.Commit()
}

// QueryOne executes a query that is expected to return a single result.
// The columns, table and selection criteria are given as Clauses.
// QueryOne returns the selected values which is never nil. ErrNoResult is
// returned when the query did not match any rows.
func (s *Session) QueryOne(cols ColsClause, from TableClause, where ...Clause) (Values, error) {
	whereClause, params := buildWhereClause(where...)
	query := fmt.Sprintf("select %s from %s %s", cols.SQL(), from.SQL(), whereClause)

	row := s.tx.QueryRowContext(s.ctx, query, params...)
	values, err := collectValues(cols.Names(), row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResult
		}
		return nil, err
	}

	return values, nil
}

// QueryMany executes a query that is expected to match any number of rows. The rows are
// returned as Values. If the query did not match any row an empty slice is returned.
func (s *Session) QueryMany(cols ColsClause, from TableClause, clauses ...Clause) ([]Values, error) {
	whereClause, params := buildWhereClause(clauses...)
	// TODO: limit, ...
	query := fmt.Sprintf("select %s from %s %s %s", cols.SQL(), from.SQL(), whereClause, buildOrderByClause(clauses))

	rows, err := s.tx.QueryContext(s.ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Values, 0)
	for rows.Next() {
		vals, err := collectValues(cols.Names(), rows)
		if err != nil {
			return nil, err
		}
		result = append(result, vals)
	}

	return result, nil
}

// QueryCount executes a counting query and returns the number of matching rows.
func (s *Session) QueryCount(from TableClause, clauses ...Clause) (count int, err error) {
	whereClause, params := buildWhereClause(clauses...)
	query := fmt.Sprintf("select count(*) from %s %s", from.SQL(), whereClause)

	row := s.tx.QueryRowContext(s.ctx, query, params...)

	err = row.Scan(&count)
	return
}

// InsertOne inserts a single row.
func (s *Session) InsertOne(into TableClause, values Values) error {
	args := make([]interface{}, 0, len(values))

	var insert strings.Builder
	insert.WriteString("insert into ")
	insert.WriteString(into.SQL())
	insert.WriteString(" (")

	firstCol := true
	for col, arg := range values {
		if !firstCol {
			insert.WriteString(", ")
		} else {
			firstCol = false
		}
		insert.WriteString(col)

		args = append(args, arg)
	}

	insert.WriteString(") values (")

	firstCol = true
	for range values {
		if !firstCol {
			insert.WriteString(", ")
		} else {
			firstCol = false
		}
		insert.WriteRune('?')
	}

	insert.WriteString(")")

	_, err := s.tx.Exec(insert.String(), args...)
	// TODO: What about the result?
	return err
}

// UpdateMany updates all matching rows with the same values given.
func (s *Session) UpdateMany(table TableClause, values Values, where ...Clause) error {
	args := make([]interface{}, len(values)+len(where))

	var update strings.Builder
	update.WriteString("update ")
	update.WriteString(table.SQL())
	update.WriteString(" set ")

	count := 0
	for col, arg := range values {
		if count > 0 {
			update.WriteString(", ")
		}
		update.WriteString(col)
		update.WriteString(" = ?")

		args[count] = arg
		count++
	}
	update.WriteRune(' ')

	whereClause, whereArgs := buildWhereClause(where...)
	update.WriteString(whereClause)
	copy(args[len(values):], whereArgs)

	_, err := s.tx.Exec(update.String(), args...)
	// TODO: What about the result?
	return err
}

// DeleteMany deletes all matching rows from the database.
func (s *Session) DeleteMany(from TableClause, where ...Clause) error {
	whereClause, whereArgs := buildWhereClause(where...)
	delete := fmt.Sprintf("delete from %s %s", from.SQL(), whereClause)

	_, err := s.tx.Exec(delete, whereArgs...)
	return err
}

// captureScanner implements the sql package's Scanner interface and captures
// the passed value. This struct is used to collect the column values for
// storing them into Values.
type captureScanner struct {
	val interface{}
}

func (c *captureScanner) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		// Copy the slice here as the reference is only valid
		// until the end of the call.
		val := make([]byte, len(v))
		copy(val, v)
		c.val = val
	default:
		c.val = v
	}
	return nil
}

// Scanner defines the Scan method provided by sql.Rows and sql.Row, as
// the sql package does not define such an interface.
type Scanner interface {
	Scan(vals ...interface{}) error
}

var _ sql.Scanner = &captureScanner{}

// collectValues collects the single row values from the given scanner and
// returns them as a Values value. Names define the column names which must
// be in the same order as they appeared in the query.
func collectValues(names []string, scanner Scanner) (Values, error) {
	scanners := make([]interface{}, 0, len(names))
	for range names {
		scanners = append(scanners, &captureScanner{})
	}

	err := scanner.Scan(scanners...)
	if err != nil {
		return nil, err
	}

	values := make(Values, len(names))
	for idx, name := range names {
		values[name] = (scanners[idx]).(*captureScanner).val
	}

	return values, nil
}
