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

package generate

import (
	"fmt"
	"strings"
)

// Type represents a mapped field's type. It provides
// methods to generate code to interact with the field.
type Type interface {
	// Expr returns a string containing a Go expression to define
	// a variable or parameter of this type.
	Expr() string

	// Returns a Go statement that obtains the column value from
	// a depot.Values object and assigns it to a variable and a
	// boolean flag indicating success.
	AssignNonNil(assignVar, okVar, valuesExpr, columnExpr string) string
}

// --

// NamedType is an implementation of Type which uses
// a bare name to refer to a type, such as string, int
// or time.Time.
type NamedType struct {
	Name string
}

var _ Type = &NamedType{}

func (n *NamedType) Expr() string {
	return n.Name
}

func (n *NamedType) AssignNonNil(assignVar, okVar, valuesExpr, columnExpr string) string {
	return fmt.Sprintf("%s, %s = %s.%s(%s)", assignVar, okVar, valuesExpr, n.valuesGetterName(), columnExpr)
}

func (n *NamedType) valuesGetterName() string {
	if n.Name == "time.Time" {
		return "GetTime"
	}

	return "Get" + strings.ToUpper(n.Name[0:1]) + n.Name[1:]
}

// --

// ByteSlice implements a Type that describes a slice of
// bytes, i.e. []byte.
type ByteSlice struct{}

var _ Type = &ByteSlice{}

func (b *ByteSlice) Expr() string {
	return "[]byte"
}

func (b *ByteSlice) AssignNonNil(assignVar, okVar, valuesExpr, columnExpr string) string {
	return fmt.Sprintf("%s, %s = %s.GetBytes(%s)", assignVar, okVar, valuesExpr, columnExpr)
}

// --

type PointerType struct {
	NamedType
}

var _ Type = &PointerType{}

func (p *PointerType) Expr() string {
	return "*" + p.NamedType.Expr()
}

func (p *PointerType) AssignNonNil(assignVar, okVar, valuesExpr, columnExpr string) string {
	return fmt.Sprintf(`if %s; k {
		%s = &u
	} else {
		%s = false
	}`, strings.Replace(p.NamedType.AssignNonNil("u", "k", valuesExpr, columnExpr), "=", ":=", 1), assignVar, okVar)
}

// --

// FieldMapping defines the mapping of a single struct field.
type FieldMapping struct {
	Field  string
	Column string
	Type   Type
	Opts   FieldOptions
}

// FieldOptions defines the additional options to be marked on field.s
type FieldOptions struct {
	// Flag indicating that this field is mapped to the primary key column
	ID bool
	// Flag indicating whether values mapped to this field can be null.
	Nullable bool
}

// StructMapping defines how a single struct is mapped.
type StructMapping struct {
	Package string
	Name    string
	Fields  []FieldMapping
}

// ID returns the field mapping defining the primary key or nil
// if no such mapping is defined.
func (s *StructMapping) ID() *FieldMapping {
	for _, f := range s.Fields {
		if f.Opts.ID {
			return &f
		}
	}

	return nil
}
