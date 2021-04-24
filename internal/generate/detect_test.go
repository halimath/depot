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
