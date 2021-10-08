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

// ClauseWriter defines the interface for types that can be used to generate SQL and bound variables for a
// query.
type ClauseWriter interface {
	// WriteString appends s to the current query.
	WriteString(s string)

	// WriteRune appends r to the current query.
	WriteRune(r rune)

	// BindParameter binds a new parameter. This method adds a parameter placeholder to the current query
	// and adds arg as the bound value for the parameter.
	BindParameter(arg interface{})
}

// Clause defines the interface implemented by all clauses used to describe different parts of a query.
type Clause interface {
	// Write writes this clause to the given writer.
	Write(w ClauseWriter)

	clause()
}

// ColsClause defines a Clause used to select columns.
type ColsClause interface {
	Clause
	Names() []string
	cols()
}

type colsClause struct {
	names []string
}

func (c *colsClause) clause() {}
func (c *colsClause) cols()   {}

func (c *colsClause) Names() []string {
	return c.names
}

func (c *colsClause) Write(w ClauseWriter) {
	for i, n := range c.names {
		if i > 0 {
			w.WriteRune(',')
		}
		w.WriteString(n)
	}
}

// Cols implements a factory for a ColsClause.
func Cols(cols ...string) ColsClause {
	return &colsClause{
		names: cols,
	}
}

// --

// TableClause implements a clause used to name a table.
type TableClause interface {
	Clause
	table()
}

type tableClause struct {
	name string
}

func (t *tableClause) clause() {}
func (t *tableClause) table()  {}

func (t *tableClause) Write(w ClauseWriter) {
	w.WriteString(t.name)
}

// Table creates a TableClause from the single table name.
func Table(name string) TableClause {
	return &tableClause{
		name: name,
	}
}

// From is an alias for Table supporting a more DSL style interface.
func From(name string) TableClause {
	return Table(name)
}

// Into is an alias for Table supporting a more DSL style interface.
func Into(name string) TableClause {
	return Table(name)
}

// --

// SelectClause is an interface for clauses that can be used in select queries, such as where, order by,
// group by, having, ...
type SelectClause interface {
	Clause
	sel()
}

// --

// OrderByClause defines an order by clause.
type OrderByClause interface {
	SelectClause
	orderBy()
}

type OrderByCol struct {
	Name      string
	Ascending bool
}

type orderByClause struct {
	cols []OrderByCol
}

func (o *orderByClause) clause()  {}
func (o *orderByClause) sel()     {}
func (o *orderByClause) orderBy() {}

func (o *orderByClause) Write(w ClauseWriter) {
	if len(o.cols) == 0 {
		return
	}

	for i, c := range o.cols {
		if i > 0 {
			w.WriteString(", ")
		}

		w.WriteString(c.Name)
		if c.Ascending {
			w.WriteString(" asc")
		} else {
			w.WriteString(" desc")
		}
	}
}

// OrderBy constructs a new OrderByClause.
func OrderBy(cols ...OrderByCol) OrderByClause {
	return &orderByClause{
		cols: cols,
	}
}

// Asc returns an OrderByCol sorting by the given column in ascending order.
func Asc(column string) OrderByCol {
	return OrderByCol{
		Name:      column,
		Ascending: true,
	}
}

// Desc returns an OrderByCol sorting by the given column in descending order.
func Desc(column string) OrderByCol {
	return OrderByCol{
		Name: column,
	}
}

// --

type WhereClause interface {
	SelectClause
	where()
}

type SearchCondition interface {
	Write(w ClauseWriter)
}

type whereClause struct {
	conditions []SearchCondition
}

func (wc *whereClause) clause() {}
func (wc *whereClause) sel()    {}
func (wc *whereClause) where()  {}
func (wc *whereClause) Write(w ClauseWriter) {
	if len(wc.conditions) == 0 {
		return
	}

	for i, c := range wc.conditions {
		if i > 0 {
			w.WriteString(" and ")
		}

		w.WriteRune('(')
		c.Write(w)
		w.WriteRune(')')
	}
}

func Where(conditions ...SearchCondition) WhereClause {
	return &whereClause{
		conditions: conditions,
	}
}

// --

type OperatorSearchCondition struct {
	Column   string
	Operator string
	Value    interface{}
}

func (o OperatorSearchCondition) Write(w ClauseWriter) {
	w.WriteString(o.Column)
	w.WriteRune(' ')
	w.WriteString(o.Operator)
	w.WriteRune(' ')
	w.BindParameter(o.Value)
}

func Eq(column string, val interface{}) SearchCondition {
	return OperatorSearchCondition{
		Column:   column,
		Operator: "=",
		Value:    val,
	}
}

func GT(column string, val interface{}) SearchCondition {
	return OperatorSearchCondition{
		Column:   column,
		Operator: ">",
		Value:    val,
	}
}

func GE(column string, val interface{}) SearchCondition {
	return OperatorSearchCondition{
		Column:   column,
		Operator: ">=",
		Value:    val,
	}
}

func LT(column string, val interface{}) SearchCondition {
	return OperatorSearchCondition{
		Column:   column,
		Operator: "<",
		Value:    val,
	}
}

func LE(column string, val interface{}) SearchCondition {
	return OperatorSearchCondition{
		Column:   column,
		Operator: "<=",
		Value:    val,
	}
}

type nullClause struct {
	col string
	not bool
}

func (c nullClause) Write(w ClauseWriter) {
	w.WriteString(c.col)
	w.WriteString(" is ")
	if c.not {
		w.WriteString("not ")
	}
	w.WriteString("null")
}

func IsNotNull(col string) SearchCondition {
	return &nullClause{
		col: col,
		not: true,
	}
}

func IsNull(col string) SearchCondition {
	return &nullClause{
		col: col,
	}
}

type inClause struct {
	column string
	values []interface{}
}

func (c inClause) Write(w ClauseWriter) {
	w.WriteString(c.column)
	w.WriteString(" in (")

	for i, v := range c.values {
		if i > 0 {
			w.WriteString(", ")
		}

		w.BindParameter(v)
	}

	w.WriteRune(')')
}

// In creates a WhereClause using the `in` operator.
func In(column string, values ...interface{}) SearchCondition {
	return &inClause{
		column: column,
		values: values,
	}
}
