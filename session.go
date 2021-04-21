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
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	// ErrNoResult is returned when queries execpted to match (at least) on row
	// match no row at all.
	ErrNoResult = errors.New("no result")

	// ErrRollback is returned when trying to commit a session that has already been
	// rolled back.
	ErrRollback = errors.New("rolled back")
)

// Session defines an interaction session with the database.
// A session uses a single transaction and is bound to a single
// Context. A session provides an abstract interface built around
// Values and Clauses.
type Session struct {
	options           *FactoryOptions
	txCount           int
	tx                *sql.Tx
	ctx               context.Context
	err               error
	alreadyRolledback bool
}

// Commit commits the session's transaction and returns an error
// if the commit fails.
func (s *Session) Commit() error {
	if s.err != nil {
		return s.err
	}
	if s.txCount > 1 {
		s.txCount--
		return nil
	}
	return s.tx.Commit()
}

// Rollback rolls the session's transaction back and returns any
// error raised during the rollback.
func (s *Session) Rollback() (err error) {
	if !s.alreadyRolledback {
		err = s.tx.Rollback()
	}

	if s.txCount > 1 {
		s.txCount--
		s.err = ErrRollback
	}

	return err
}

// Error marks the transaction as failed so it cannot be committed
// later on.
// Calling Error with a nil error clears the error state of the transaction.
func (s *Session) Error(err error) {
	s.err = err
}

// QueryOne executes a query that is expected to return a single result.
// The columns, table and selection criteria are given as Clauses.
// QueryOne returns the selected values which is never nil. ErrNoResult is
// returned when the query did not match any rows.
func (s *Session) QueryOne(cols ColsClause, from TableClause, where ...Clause) (Values, error) {
	whereClause, params := buildWhereClause(where...)
	query := fmt.Sprintf("select %s from %s %s", cols.SQL(), from.SQL(), whereClause)

	if s.options.LogSQL {
		log.Printf("QueryOne: '%s'", query)
	}

	row := s.tx.QueryRowContext(s.ctx, query, params...)
	values, err := collectValues(cols.Names(), row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoResult
		}
		return nil, fmt.Errorf("failed to execute '%s': %w", query, err)
	}

	return values, nil
}

// QueryMany executes a query that is expected to match any number of rows. The rows are
// returned as Values. If the query did not match any row an empty slice is returned.
func (s *Session) QueryMany(cols ColsClause, from TableClause, clauses ...Clause) ([]Values, error) {
	whereClause, params := buildWhereClause(clauses...)
	// TODO: limit, ...
	query := fmt.Sprintf("select %s from %s %s %s", cols.SQL(), from.SQL(), whereClause, buildOrderByClause(clauses))
	if s.options.LogSQL {
		log.Printf("QueryMany: '%s'", query)
	}

	rows, err := s.tx.QueryContext(s.ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute '%s': %w", query, err)
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
	if s.options.LogSQL {
		log.Printf("QueryCount: '%s'", query)
	}

	row := s.tx.QueryRowContext(s.ctx, query, params...)

	err = row.Scan(&count)
	if err != nil {
		err = fmt.Errorf("failed to execute '%s': %w", query, err)
	}
	return
}

// Exec executes the given query passing the given args and
// returns the resulting error or nil.
// This is just a wrapper for calling ExexContext on the wrapped
// transaction.
func (s *Session) Exec(query string, args ...interface{}) error {
	_, err := s.tx.ExecContext(s.ctx, query, args...)
	return err
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
	query := insert.String()
	if s.options.LogSQL {
		log.Printf("InsertOne: '%s'", query)
	}

	return s.Exec(query, args...)
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

	query := update.String()
	if s.options.LogSQL {
		log.Printf("UpdateMany: '%s'", query)
	}

	return s.Exec(query, args...)
}

// DeleteMany deletes all matching rows from the database.
func (s *Session) DeleteMany(from TableClause, where ...Clause) error {
	whereClause, whereArgs := buildWhereClause(where...)

	query := fmt.Sprintf("delete from %s %s", from.SQL(), whereClause)
	if s.options.LogSQL {
		log.Printf("DeleteMany: '%s'", query)
	}

	return s.Exec(query, whereArgs...)
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
