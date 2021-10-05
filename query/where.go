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

package query

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
	arg      interface{}
}

func (c *operatorWhereClause) Write(w Writer) {
	w.WriteString(c.column)
	w.WriteRune(' ')
	w.WriteString(c.operator)
	w.WriteRune(' ')
	w.BindParameter(c.arg)
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
		arg:      value,
	}
}

// GT creates a WhereClause comparing a column's value for greater than.
func GT(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: ">",
		arg:      value,
	}
}

// GE creates a WhereClause comparing a column's value for greater or equal.
func GE(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: ">=",
		arg:      value,
	}
}

// LT creates a WhereClause comparing a column's value for less than.
func LT(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: "<",
		arg:      value,
	}
}

// LE creates a WhereClause comparing a column's value for less equal.
func LE(column string, value interface{}) WhereClause {
	return &operatorWhereClause{
		column:   column,
		operator: "<=",
		arg:      value,
	}
}

// inClause implements a WhereClause that uses the `in` operator.
type inClause struct {
	column string
	values []interface{}
}

func (c *inClause) Write(w Writer) {
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

func (c *inClause) where() {}

// In creates a WhereClause using the `in` operator.
func In(column string, values ...interface{}) WhereClause {
	return &inClause{
		column: column,
		values: values,
	}
}

type nullClause struct {
	col string
	not bool
}

var _ WhereClause = &nullClause{}

func (c *nullClause) where() {}

func (c *nullClause) Write(w Writer) {
	w.WriteString(c.col)
	w.WriteString(" is ")
	if c.not {
		w.WriteString("not ")
	}
	w.WriteString("null")
}

func IsNotNull(col string) WhereClause {
	return &nullClause{
		col: col,
		not: true,
	}
}

func IsNull(col string) WhereClause {
	return &nullClause{
		col: col,
	}
}
