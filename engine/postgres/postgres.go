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

package postgres

import (
	"strconv"
	"strings"

	"github.com/halimath/depot"
)

type Dialect struct{}

var _ depot.Dialect = &Dialect{}

func (d *Dialect) NewClauseBuilder() depot.QueryBuilder { return &clauseBuilder{} }

// --

type clauseBuilder struct {
	sql  strings.Builder
	args []interface{}
}

var _ depot.QueryBuilder = &clauseBuilder{}

func (b *clauseBuilder) WriteString(s string) { b.sql.WriteString(s) }
func (b *clauseBuilder) WriteRune(r rune)     { b.sql.WriteRune(r) }
func (b *clauseBuilder) SQL() string          { return b.sql.String() }
func (b *clauseBuilder) Args() []interface{}  { return b.args }
func (b *clauseBuilder) BindParameter(arg interface{}) {
	b.args = append(b.args, arg)
	b.sql.WriteRune('$')
	b.sql.WriteString(strconv.Itoa(len(b.args)))
}
