package depot

import (
	"fmt"
	"strings"
)

// WhereClause defines the interface implemented by all
// clauses that contribute a "where" condition.
type WhereClause interface {
	Clause
	where()
}

// operatorWhereClause implements a WhereClause connecting a
// single column to a literal value using an SQL operator.
type operatorWhereClause struct {
	column   string
	operator string
	args     []interface{}
}

func (c *operatorWhereClause) SQL() string {
	return fmt.Sprintf("%s %s ?", c.column, c.operator)
}

func (c *operatorWhereClause) Args() []interface{} {
	return c.args
}

func (c *operatorWhereClause) where() {}

var _ WhereClause = &operatorWhereClause{}

// Where is an alias for EQ.
func Where(column string, value interface{}) WhereClause {
	return EQ(column, value)
}

// EQ creates a WhereClause comparing a column's value for equality.
func EQ(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: "=",
		args:     []interface{}{value},
	}
}

// GT creates a WhereClause comparing a column's value for greater than.
func GT(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: ">",
		args:     []interface{}{value},
	}
}

// GE creates a WhereClause comparing a column's value for greater or equal.
func GE(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: ">=",
		args:     []interface{}{value},
	}
}

// LT creates a WhereClause comparing a column's value for less than.
func LT(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: "<",
		args:     []interface{}{value},
	}
}

// LE creates a WhereClause comparing a column's value for less equal.
func LE(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: "<=",
		args:     []interface{}{value},
	}
}

// inClause implements a WhereClause that uses the `in` operator.
type inClause struct {
	column string
	values []interface{}
}

func (c *inClause) SQL() string {
	var s strings.Builder
	s.WriteString(c.column)
	s.WriteString(" in (")

	first := true
	for range c.values {
		if !first {
			s.WriteString(", ")
		} else {
			first = false
		}

		s.WriteRune('?')
	}

	s.WriteRune(')')

	return s.String()
}

func (c *inClause) Args() []interface{} {
	return c.values
}

func (c *inClause) where() {}

// In creates a WhereClause using the `in` operator.
func In(column string, values ...interface{}) WhereClause {
	return &inClause{
		column: column,
		values: values,
	}
}

// buildWhereClause selects all where clauses from the given clauses
// and joins them together with AND. If the result contains at least
// one clause, the keyword `where` is put in front. The function also
// returns the collected arguments.
func buildWhereClause(clauses ...Clause) (string, []interface{}) {
	var b strings.Builder
	args := make([]interface{}, 0, len(clauses))

	for _, clause := range clauses {
		if w, ok := clause.(WhereClause); ok {
			if b.Len() > 0 {
				b.WriteString(" and ")
			}
			b.WriteString(w.SQL())
			args = append(args, w.Args()...)
		}
	}

	if b.Len() > 0 {
		return fmt.Sprintf("where %s", b.String()), args
	}

	return b.String(), args
}

// --

// OrderByClause defines the interface used to sort rows.
type OrderByClause interface {
	Clause
	orderBy()
}

type orderByClause struct {
	column string
	asc    bool
}

func (c *orderByClause) SQL() string {
	direction := "asc"
	if !c.asc {
		direction = "desc"
	}

	return fmt.Sprintf("%s %s", c.column, direction)
}

func (c *orderByClause) Args() []interface{} {
	return nil
}

func (c *orderByClause) orderBy() {}

var _ OrderByClause = &orderByClause{}

// OrderBy constructs a new OrderByClause.
func OrderBy(column string, asc bool) OrderByClause {
	return &orderByClause{
		column: column,
		asc:    asc,
	}
}

// Asc returns an OrderByClause sorting by the given column in
// ascending order.
func Asc(column string) OrderByClause {
	return OrderBy(column, true)
}

// Desc returns an OrderByClause sorting by the given column in
// descending order.
func Desc(column string) OrderByClause {
	return OrderBy(column, false)
}

// buildOrderByClause selects all OrderByClauses and joins them together.
// If at least on clause is selected, the keyphrase `order by` is put in front.
func buildOrderByClause(clauses []Clause) string {
	var b strings.Builder

	for _, clause := range clauses {
		if w, ok := clause.(OrderByClause); ok {
			if b.Len() > 0 {
				b.WriteString(", ")
			}
			b.WriteString(w.SQL())
		}
	}

	if b.Len() > 0 {
		return fmt.Sprintf("order by %s", b.String())
	}

	return b.String()
}
