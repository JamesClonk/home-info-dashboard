package web

import (
	"html/template"
	"strings"

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
			template.FuncMap{
				"divide": func(input, divisor int64) float64 {
					return float64(input) / float64(divisor)
				},
				"emoji": func(name string) string {
					switch strings.ToLower(name) {
					case "living room":
						return "📺"
					case "home office":
						return "🖥️"
					case "bedroom":
						return "🛏️"
					case "air quality lamp":
						return "💡"
					case "food plants lamp":
						return "💡"
					case "weather forecast":
						return "🌧️"
					}

					if strings.Contains(strings.ToLower(name), "food plant") {
						return "🍅"
					}
					if strings.Contains(strings.ToLower(name), "air quality") {
						return "🌿"
					}

					return ""
				},
			},
		},
	})
}

func Render() *render.Render {
	return r
}
