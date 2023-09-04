package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

// TemplateRenderer is a custom renderer for templates.
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders the HTML templates.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// Use a custom renderer for HTMX templates.
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	// Serve static files (e.g., CSS, JavaScript, and HTMX).
	e.Static("/static", "static")

	// Define a route for the main page.
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	// Toggle the content.
	var showContent = false
	e.POST("/toggle", func(c echo.Context) error {
		showContent = !showContent
		var content string
		if showContent {
			content = "This is now visible"
		} else {
			content = "This is hidden"
		}
		return c.HTML(http.StatusOK, content)
	})

	// Start the server.
	e.Start(":8080")
}
