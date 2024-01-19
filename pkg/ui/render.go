package ui

import (
	"html/template"
	"io"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/egoodhall/go-links/pkg/config"
	"github.com/labstack/echo/v4"
)

var index *template.Template = template.Must(
	template.New("page").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"joinArgUrls": func(urls []config.ArgUrl) string {
			res := new(strings.Builder)
			for i, url := range urls {
				if i > 0 {
					res.WriteRune('\n')
				}
				res.WriteString(string(url))
			}
			return res.String()
		},
	}).ParseFS(templates, "templates/*.go.html"))

var Render = RenderFunc(func(to io.Writer, name string, data interface{}, ctx echo.Context) error {
	return index.ExecuteTemplate(to, name, data)
})

type RenderFunc func(to io.Writer, name string, data interface{}, ctx echo.Context) error

func (rf RenderFunc) Render(to io.Writer, name string, data interface{}, ctx echo.Context) error {
	return rf(to, name, data, ctx)
}
