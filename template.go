package main

import (
	_ "embed"
	"errors"
	"os"
	"strings"
	"text/template"
)

type TemplateType string

const (
	ListTemplate  = TemplateType("list")
	TableTemplate = TemplateType("table")
)

var availableFormats = []string{string(ListTemplate), string(TableTemplate)}

//go:embed list_template.md
var list string

//go:embed table_template.md
var table string

var temp *template.Template

type T struct {
	Total int
	Stars map[string][]Star
}

func initTemplate(t TemplateType) (err error) {
	var md string
	switch t {
	case TableTemplate:
		md = table
	default:
		md = list
	}

	temp, err = template.New("readme").
		Funcs(template.FuncMap{"anchor": Anchor}).Parse(md)

	return
}

func writeList(path string, stars map[string][]Star, total int) error {
	if temp == nil {
		return errors.New("template not initialized")
	}

	if exists(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o665)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	return temp.Execute(f, T{Stars: stars, Total: total})
}

// todo: check rules for markdown anchors...
func Anchor(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "-", -1)
}
