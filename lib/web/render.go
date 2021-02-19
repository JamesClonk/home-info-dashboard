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
					case "plant room":
						return "🌻"
					case "living room":
						return "📺"
					case "home office":
						return "🖥️"
					case "bedroom":
						return "🛏️"
					case "bedroom lamp":
						return "💡"
					case "plant room lamp":
						return "💡"
					case "air quality lamp":
						return "💡"
					case "food plants lamp":
						return "💡"
					case "weather forecast":
						return "🌧️"
					}

					if strings.Contains(strings.ToLower(name), "chili") ||
						strings.Contains(strings.ToLower(name), "pepper") ||
						strings.Contains(strings.ToLower(name), "jalapeno") ||
						strings.Contains(strings.ToLower(name), "capsicum") {
						return "🌶️"
					}
					if strings.Contains(strings.ToLower(name), "food plant") ||
						strings.Contains(strings.ToLower(name), "lycopersicum") ||
						strings.Contains(strings.ToLower(name), "tomato") {
						return "🍅"
					}
					if strings.Contains(strings.ToLower(name), "salad") ||
						strings.Contains(strings.ToLower(name), "lactuca") {
						return "🥬"
					}

					if strings.Contains(strings.ToLower(name), "raised bed") ||
						strings.Contains(strings.ToLower(name), "balcony") ||
						strings.Contains(strings.ToLower(name), "garden") ||
						strings.Contains(strings.ToLower(name), "outdoor") {
						return "🏝"
					}

					if strings.Contains(strings.ToLower(name), "air quality") ||
						strings.Contains(strings.ToLower(name), "epipremnum") ||
						strings.Contains(strings.ToLower(name), "aureum") ||
						strings.Contains(strings.ToLower(name), "sansevieria") ||
						strings.Contains(strings.ToLower(name), "chlorophytum") {
						return "🌿"
					}

					return "❓"
				},
				"moisture": func(value int64) string {
					if value >= 86 {
						return "🤢🌊"
					}
					if value >= 70 {
						return "😄🌊"
					}
					if value >= 64 {
						return "😅💦"
					}
					if value >= 52 {
						return "😥💧"
					}
					return "😫🔥"
				},
			},
		},
	})
}

func Render() *render.Render {
	return r
}
