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

// Writer defines the interface for types that can be used to generate SQL and bound variables for a
// query.
type Writer interface {
	// WriteString appends s to the current query.
	WriteString(s string)

	// WriteRune appends r to the current query.
	WriteRune(r rune)

	// BindParameter binds a new parameter. This method adds a parameter placeholder to the current query
	// and adds arg as the bound value for the parameter.
	BindParameter(arg interface{})
}

// Clause defines the interface implemented by all clauses used to describe different parts of a query.
// A clause captures optional arguments.
type Clause interface {
	// Write writes this clause to the given writer.
	Write(w Writer)
}

// ColsClause defines a Clause used to select columns.
type ColsClause struct {
	names []string
}

// Names returns c column names.
func (c *ColsClause) Names() []string {
	return c.names
}

func (c *ColsClause) Write(w Writer) {
	for i, n := range c.names {
		if i > 0 {
			w.WriteRune(',')
		}
		w.WriteString(n)
	}
}

// Cols implements a factory for a ColsClause.
func Cols(cols ...string) *ColsClause {
	return &ColsClause{
		names: cols,
	}
}

// --

// TableClause implements a clause used to name a table.
type TableClause struct {
	name string
}

// Name returns c's name.
func (c *TableClause) Name() string {
	return c.name
}

func (c *TableClause) Write(w Writer) {
	w.WriteString(c.name)
}

// Table creates a TableClause from the single table name.
func Table(name string) *TableClause {
	return &TableClause{
		name: name,
	}
}

// From is an alias for Table supporting a more DSL style interface.
func From(name string) *TableClause {
	return Table(name)
}

// Into is an alias for Table supporting a more DSL style interface.
func Into(name string) *TableClause {
	return Table(name)
}
