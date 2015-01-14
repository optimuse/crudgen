package main

//go:generate go-bindata -o bindata.go template/...

import (
	"strings"
	"text/template"
)

func unCapitalize(s string) string {
	if len(s) < 1 {
		return s
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

var funcMap = template.FuncMap{
	"UnCapitalize": unCapitalize,
}

func parseTmpl(name string, t *template.Template) *template.Template {
	data, err := Asset("template/" + name)
	if err != nil {
		panic(err)
	}
	if t == nil {
		return template.Must(template.New(name).Funcs(funcMap).Parse(string(data)))
	} else {
		t, err = t.New(name).Funcs(funcMap).Parse(string(data))
		if err != nil {
			panic(err)
		}
	}
	return t
}

var tmpl *template.Template

func init() {
	tmpl = parseTmpl("list.tmpl", nil)
	tmpl = parseTmpl("create.tmpl", tmpl)
	tmpl = parseTmpl("read.tmpl", tmpl)
	tmpl = parseTmpl("update.tmpl", tmpl)
	tmpl = parseTmpl("delete.tmpl", tmpl)
	tmpl = parseTmpl("crud.tmpl", tmpl)
}
