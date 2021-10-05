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

// OrderByClause defines an order by clause.
type OrderByClause struct {
	column string
	asc    bool
}

func (c *OrderByClause) Write(w Writer) {
	w.WriteString(c.column)
	w.WriteRune(' ')

	if c.asc {
		w.WriteString("asc")
	} else {
		w.WriteString("desc")
	}
}

var _ Clause = &OrderByClause{}

// OrderBy constructs a new OrderByClause.
func OrderBy(column string, asc bool) *OrderByClause {
	return &OrderByClause{
		column: column,
		asc:    asc,
	}
}

// Asc returns an OrderByClause sorting by the given column in
// ascending order.
func Asc(column string) *OrderByClause {
	return OrderBy(column, true)
}

// Desc returns an OrderByClause sorting by the given column in
// descending order.
func Desc(column string) *OrderByClause {
	return OrderBy(column, false)
}
