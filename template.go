package main

import (
	_ "embed"
	"os"
	"strings"
	"text/template"
)

//go:embed template.md
var md string

var temp = template.Must(template.New("readme").
	Funcs(template.FuncMap{"anchor": Anchor}).Parse(md))

type T struct {
	Total int
	Stars map[string][]Star
}

func writeReadme(stars map[string][]Star, total int) error {

	f, err := os.OpenFile("test.md", os.O_RDWR|os.O_CREATE, 0665)
	if err != nil {
		return err
	}
	defer f.Close()

	return temp.Execute(f, T{Stars: stars, Total: total})
}

// todo: check rules for markdown anchors...
func Anchor(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "-", -1)
}
