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
	"reflect"
	"testing"
)

func Test_detectMapping(t *testing.T) {
	actual, err := detectMapping("test.go", `
		package models

		type (
			// Message demonstrates a persistent struct showing several mapped fields.
			Message struct {
				ID         string     "depot:\"id,id\""
				Text       string     "depot:\"text\""
				OrderIndex int        "depot:\"order_index\""
				Length     float32    "depot:\"len\""
				Attachment []byte     "depot:\"attachment\""
				Created    time.Time  "depot:\"created\""
				Updated	   *time.Time "depot:\"updated,nullable\""
			}
		)`, "Message")

	if err != nil {
		t.Fatalf("expected no error but got %s", err)
	}

	expected := StructMapping{
		Package: "models",
		Name:    "Message",
		Fields: []FieldMapping{
			{
				Field:  "ID",
				Column: "id",
				Type: &NamedType{
					Name: "string",
				},
				Opts: FieldOptions{
					ID: true,
				},
			},
			{
				Field:  "Text",
				Column: "text",
				Type: &NamedType{
					Name: "string",
				},
			},
			{
				Field:  "OrderIndex",
				Column: "order_index",
				Type: &NamedType{
					Name: "int",
				},
			},
			{
				Field:  "Length",
				Column: "len",
				Type: &NamedType{
					Name: "float32",
				},
			},
			{
				Field:  "Attachment",
				Column: "attachment",
				Type:   &ByteSlice{},
			},
			{
				Field:  "Created",
				Column: "created",
				Type: &NamedType{
					Name: "time.Time",
				},
			},
			{
				Field:  "Updated",
				Column: "updated",
				Type: &PointerType{
					NamedType: NamedType{
						Name: "time.Time",
					},
				},
				Opts: FieldOptions{
					Nullable: true,
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, *actual) {
		t.Errorf("expected %#v but got %#v", expected, *actual)
	}

	if !reflect.DeepEqual(expected.Fields[0], *actual.ID()) {
		t.Errorf("expected id mapping %#v but got %#v", expected.Fields[0], *actual.ID())
	}

}
