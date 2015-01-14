package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"

	"go/ast"
)

var idRegexp = regexp.MustCompile("^([a-zA-Z]+)ID$")

type Visitor struct {
	search string
	found  *bool

	Data *StructData
}

func (v Visitor) Visit(node ast.Node) (w ast.Visitor) {
	if d, ok := node.(*ast.GenDecl); ok {
		for _, s := range d.Specs {
			if t, ok := s.(*ast.TypeSpec); ok {
				//all structs
				if st, ok := (t.Type).(*ast.StructType); ok {
					//match given name
					if t.Name.String() == v.search {
						*(v.found) = true
						v.Data.Struct = v.search
						//create fields
						if st.Fields != nil {
							for _, f := range st.Fields.List {
								if len(f.Names) == 1 {
									field := Field{Name: f.Names[0].String()}

									//get string of type
									switch typ := (f.Type).(type) {
									case *ast.Ident:
										field.Type = typ.String()
									case *ast.ArrayType:
										field.Type = "[]" + (typ.Elt).(*ast.Ident).String()
									default:
										continue
									}

									v.Data.Fields = append(v.Data.Fields, field)
								}

								//determine parent
								if f.Tag != nil {
									structTag := reflect.StructTag(f.Tag.Value[1:len(f.Tag.Value)])
									if structTag.Get("crudgen") == "parent" && len(f.Names) == 1 {
										parent := idRegexp.FindStringSubmatch(f.Names[0].String())
										if len(parent) == 2 {
											v.Data.Parent = parent[1]
										}
									}
								}
							}

							return nil
						}
					}
				}
			}
		}
	}
	return v
}

func Parse(dir, structName string) (*StructData, error) {
	fset := token.NewFileSet()
	pkgset, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	v := Visitor{search: structName, found: new(bool), Data: new(StructData)}

	for _, pkg := range pkgset {
		for _, f := range pkg.Files {
			v.Data.Package = f.Name.String()
			ast.Walk(v, f)
			if *(v.found) {
				return v.Data, nil
			}
		}
	}

	return nil, fmt.Errorf("Struct Not Found")
}
