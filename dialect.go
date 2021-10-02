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
	"strings"

	"github.com/halimath/depot/query"
)

// QueryBuilder defines the interface for types that are used to build clauses.
type QueryBuilder interface {
	query.Writer

	SQL() string
	Args() []interface{}
}

// Dialect abstracts the differences in SQL for different database engines.
type Dialect interface {
	// NewClauseBuilder creates a new QueryBuilder matching the selected database.
	NewClauseBuilder() QueryBuilder
}

// --

type DefaultClauseBuilder struct {
	sql  strings.Builder
	args []interface{}
}

var _ QueryBuilder = &DefaultClauseBuilder{}

func (b *DefaultClauseBuilder) WriteString(s string) { b.sql.WriteString(s) }
func (b *DefaultClauseBuilder) WriteRune(r rune)     { b.sql.WriteRune(r) }
func (b *DefaultClauseBuilder) SQL() string          { return b.sql.String() }
func (b *DefaultClauseBuilder) Args() []interface{}  { return b.args }
func (b *DefaultClauseBuilder) BindParameter(arg interface{}) {
	b.sql.WriteRune('?')
	b.args = append(b.args, arg)
}

// --

type DefaultDialect struct {
}

var _ Dialect = &DefaultDialect{}

func (d *DefaultDialect) NewClauseBuilder() QueryBuilder {
	return &DefaultClauseBuilder{}
}
