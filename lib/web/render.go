package web

import (
	"html/template"

	"github.com/unrolled/render"
)

var r *render.Render

func init() {
	// setup template rendering
	r = render.New(render.Options{
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		IndentJSON: true,
		Funcs: []template.FuncMap{
			template.FuncMap{"divide": func(input, divisor int64) float64 {
				return float64(input) / float64(divisor)
			},
			},
		},
	})
}

func Render() *render.Render {
	return r
}
