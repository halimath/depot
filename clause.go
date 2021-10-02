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

import "strings"

// Clause defines the interface implemented by all clauses used to describe different parts of a query.
// A clause captures optional arguments.
type Clause interface {
	// SQL returns the SQL query part expressed by this clause.
	SQL() string

	// Args returns any positional arguments used by this clause.
	Args() []interface{}
}

// ColsClause defines a Clause used to select columns.
type ColsClause interface {
	Clause

	// Names returns the list of column names to select in query order.
	Names() []string
	cols()
}

type colsClause struct {
	names []string
}

func (c *colsClause) SQL() string {
	return strings.Join(c.names, ", ")
}

func (c *colsClause) Names() []string {
	return c.names
}

func (c *colsClause) Args() []interface{} {
	return nil
}

func (c *colsClause) cols() {}

// Cols implements a factory for a ColsClause.
func Cols(cols ...string) ColsClause {
	return &colsClause{
		names: cols,
	}
}

// func cols(clauses ...Clause) ColsClause {
// 	for _, clause := range clauses {
// 		if cols, ok := clause.(ColsClause); ok {
// 			return cols
// 		}
// 	}

// 	panic("no cols given")
// }

// --

// TableClause implements a clause used to name a table.
type TableClause interface {
	Clause
	tableClause()
}

type tableClause struct {
	table string
}

func (c *tableClause) SQL() string {
	return c.table
}

func (c *tableClause) Args() []interface{} {
	return nil
}

func (c *tableClause) tableClause() {}

// Table creates a TableClause from the single table name.
func Table(table string) TableClause {
	return &tableClause{
		table: table,
	}
}

// From is an alias for Table supporting a more DSL style interface.
func From(table string) TableClause {
	return Table(table)
}

// Into is an alias for Table supporting a more DSL style interface.
func Into(table string) TableClause {
	return Table(table)
}

// func table(clauses ...Clause) TableClause {
// 	for _, clause := range clauses {
// 		if tableClause, ok := clause.(TableClause); ok {
// 			return tableClause
// 		}
// 	}

// 	panic("no table clause given")
// }
