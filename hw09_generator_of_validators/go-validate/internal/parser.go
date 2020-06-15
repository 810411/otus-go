package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type StructV struct {
	Name    string
	VarName string
	Fields  []FieldV
}

type FieldV struct {
	VarName string
	Name    string
	Type    string
	VType   string
	VValue  string
}

func Parse(path string) (string, []StructV, error) {
	structVMap := make(map[string]StructV)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}

	for _, decl := range f.Decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range decl.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			fillStructVMap(currType, structVMap)
		}
	}

	if len(structVMap) == 0 {
		return "", nil, fmt.Errorf("file %s doesn't contain structs fields for generating validators", path)
	}
	res := make([]StructV, 0, len(structVMap))
	for _, v := range structVMap {
		res = append(res, v)
	}

	return f.Name.Name, res, nil
}

func fillStructVMap(currType *ast.TypeSpec, structVMap map[string]StructV) {
	if currStruct, ok := currType.Type.(*ast.StructType); ok {
		structName := currType.Name.String()

		for _, field := range currStruct.Fields.List {
			if field.Tag == nil {
				continue
			}

			vStr := getValidatorStr(field.Tag.Value)
			if vStr == "" {
				continue
			}

			fieldType, ok := getTypeName(field)
			if !ok {
				continue
			}

			vStrArr := strings.Split(vStr, "|")
			for _, v := range vStrArr {
				vArr := strings.SplitN(v, ":", 2)
				varName := strings.ToLower(structName[0:1])
				structVMap[structName] = StructV{
					Name:    structName,
					VarName: varName,
					Fields: append(
						structVMap[structName].Fields,
						FieldV{
							varName,
							field.Names[0].String(),
							fieldType,
							vArr[0], vArr[1],
						},
					),
				}
			}
		}
	}
}

func getValidatorStr(tagValue string) string {
	split := strings.Split(strings.Trim(tagValue, "`"), " ")
	for _, v := range split {
		if strings.HasPrefix(v, "validate") {
			return strings.Trim(strings.TrimPrefix(v, "validate:"), "\"")
		}
	}
	return ""
}

func getTypeName(f *ast.Field) (string, bool) {
	checkStringOrInt := func(expr ast.Expr) (string, bool) {
		typeStr := fmt.Sprint(expr)
		if typeStr == "string" || typeStr == "int" {
			return typeStr, true
		}
		return "", false
	}

	switch t := f.Type.(type) {
	case *ast.Ident:
		if typeStr, ok := checkStringOrInt(t); ok {
			return typeStr, true
		}
	case *ast.ArrayType:
		if typeStr, ok := checkStringOrInt(t.Elt); ok {
			return "[]" + typeStr, true
		}
	}
	return "", false
}
