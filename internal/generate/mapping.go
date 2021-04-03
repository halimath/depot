package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type FieldMapping struct {
	Field  string
	Column string
	Type   string
	Opts   FieldOptions
}

func (f *FieldMapping) SQLType() string {
	if strings.HasPrefix(f.Type, "int") || strings.HasPrefix(f.Type, "uint") {
		return "int64"
	}
	if strings.HasPrefix(f.Type, "float") {
		return "float64"
	}

	return f.Type
}

type FieldOptions struct {
	ID bool
}

type StructMapping struct {
	Package string
	Name    string
	Fields  []FieldMapping
}

func (s *StructMapping) ID() *FieldMapping {
	for _, f := range s.Fields {
		if f.Opts.ID {
			return &f
		}
	}

	return nil
}

func determineMapping(filename string, typename string) (*StructMapping, error) {
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, filename, nil, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}

	var result *StructMapping
	var resultErr error

	ast.Inspect(fileAst, func(n ast.Node) bool {
		if result != nil || err != nil {
			return false
		}

		switch node := n.(type) {
		case *ast.TypeSpec:
			if node.Name.Name != typename {
				return true
			}

			strct, ok := node.Type.(*ast.StructType)

			if !ok {
				return true
			}

			result = &StructMapping{
				Package: fileAst.Name.Name,
				Name:    typename,
				Fields:  make([]FieldMapping, 0, len(strct.Fields.List)),
			}

			for _, f := range strct.Fields.List {
				if f.Tag == nil {
					continue
				}

				t, err := typeString(f.Type)
				if err != nil {
					resultErr = err
					return false
				}

				fieldMapping := FieldMapping{
					Field: f.Names[0].Name,
					Type:  t,
				}

				if ok := parseTag(f.Tag.Value, &fieldMapping); !ok {
					continue
				}

				result.Fields = append(result.Fields, fieldMapping)
			}
		}
		return true
	})

	if resultErr != nil {
		return nil, resultErr
	}

	if result == nil {
		return nil, fmt.Errorf("type %s not found in %s", typename, filename)
	}

	return result, nil
}

func typeString(t ast.Expr) (string, error) {
	switch typ := t.(type) {
	case *ast.Ident:
		// TODO: Restrict types to those supported by sql package, i.e. int, float, bool, ...
		return typ.Name, nil
	case *ast.ArrayType:
		// TODO: Restrict type to be byte slice
		if typ.Len != nil {
			return "", fmt.Errorf("unsupported fixed length array: %#v", t)
		}
		t, err := typeString(typ.Elt)
		return "[]" + t, err
	case *ast.SelectorExpr:
		// TODO: Restrict type to be time.Time
		// TODO: Alternatively, add support for embedded types
		t, err := typeString(typ.X)
		return t + "." + typ.Sel.Name, err
	default:
		return "", fmt.Errorf("unsupported persistent field type: %#v", t)
	}
}

func parseTag(tag string, f *FieldMapping) (ok bool) {
	val, ok := findDepotTagValue(tag)
	if !ok {
		return false
	}

	parts := strings.Split(val, ",")
	f.Column = parts[0]

	if len(parts) > 1 && strings.ToLower(parts[1]) == "id" {
		f.Opts.ID = true
	}

	return true
}

func findDepotTagValue(tag string) (value string, ok bool) {
	tag, err := strconv.Unquote(tag)
	if err != nil {
		return "", false
	}
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if name == "depot" {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}
