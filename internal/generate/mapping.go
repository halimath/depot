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
	"strings"
)

// FieldMapping defines the mapping of a single struct field.
type FieldMapping struct {
	Field  string
	Column string
	Type   string
	Opts   FieldOptions
}

// ValuesGetterName returns the name of the Values.Get* method
// to invoke to obtain a value assignable to the target field.
func (f *FieldMapping) ValuesGetterName() string {
	if f.Type == "[]byte" {
		return "GetBytes"
	}

	if f.Type == "time.Time" {
		return "GetTime"
	}

	return "Get" + strings.ToUpper(f.Type[0:1]) + f.Type[1:]
}

// FieldOptions defines the additional options to be marked on field.s
type FieldOptions struct {
	ID bool
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
