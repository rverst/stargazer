package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type TemplateType string

const (
	ListTemplate  = TemplateType("list")
	TableTemplate = TemplateType("table")

	creditText = "A list of awesome repositories I've starred. Wan't your own? Try: "
	creditUrl  = "https://github.com/rverst/stargazer"
)

var availableFormats = []string{string(ListTemplate), string(TableTemplate)}

//go:embed list_template.md
var list string

//go:embed table_template.md
var table string

var temp *template.Template

type T struct {
	Total       int
	WithToc     bool
	WithLicense bool
	WithStars   bool
	Toc         string
	Stars       map[string][]Star
	Credits     C
}

type C struct {
	Text string
	Url  string
	Link string
}

func initTemplate(tType string) (err error) {
	var t string
	switch tType {
	case "table":
		t = table
	case "list":
		t = list
	default:
		if exists(tType) {
			if b, err := os.ReadFile(tType); err != nil {
				fmt.Printf("cannot read custom template: %s\n%v\n", tType, err)
			} else {
				t = string(b)
			}
		} else {
			t = list
		}
	}

	temp, err = template.New("readme").
		Funcs(template.FuncMap{"anchor": Anchor}).Parse(t)

	return
}

func writeList(path string, stars map[string][]Star, total int, withToc, withLicense, withStars bool) error {
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

	c := C{
		Text: creditText,
		Url:  creditUrl,
		Link: "[stargazer](" + creditUrl + ")",
	}

	return temp.Execute(f, T{
		Stars:       stars,
		Toc:         toc(stars),
		Total:       total,
		Credits:     c,
		WithToc:     withToc,
		WithLicense: withLicense,
		WithStars:   withStars,
	})
}

// todo: check rules for markdown anchors...
func Anchor(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "-", -1)
}

// toc returns a table of contents
func toc(stars map[string][]Star) string {
	rx := regexp.MustCompile(`[^\w\- ]`) // regexp to remove all punctuation

	header := make(map[string]string, 0)

	for k := range stars {
		x := strings.TrimSpace(k)
		x = rx.ReplaceAllString(x, "")
		x = strings.ReplaceAll(x, " ", "-")

		c := 0
		add := true
		for {
			y := x
			if c > 0 {
				y += fmt.Sprintf("-%d", c)
			}
			for h := range header {
				if h == y {
					c++
					add = false
					break
				}
			}

			if !add {
				continue
			}
			header[y] = k
			break
		}
	}
	sb := strings.Builder{}
	for k, v := range header {
		sb.WriteString(fmt.Sprintf("  * [%s](#%s)", v, k))
	}
	return sb.String()
}
