package main

type Field struct {
	Name string
	Type string
}

type StructData struct {
	Package string
	Struct  string
	Parent  string
	Fields  []Field
}

func (d StructData) FieldsNoID() []Field {
	fields := make([]Field, 0)
	for _, f := range d.Fields {
		if f.Name != "ID" {
			fields = append(fields, f)
		}
	}
	return fields
}
