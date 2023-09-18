package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

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

	e.Static("/dist", "dist")

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

	e.GET("/pokemon/random", func(c echo.Context) error {
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(1017)
		url := fmt.Sprintf(`https://pokeapi.co/api/v2/pokemon/%v`, randNum)
		req, _ := http.NewRequest("GET", url, nil)

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		var data map[string]interface{}
		err := json.Unmarshal([]byte(body), &data)
		if err != nil {
			fmt.Printf("could not unmarshal json: %s\n", err)
			return c.HTML(http.StatusBadRequest, "API error")
		}

		// var sprites map[string]interface{}
		rawName := data["name"]
		sprites := data["sprites"].(map[string]interface{})
		rawHeight := data["height"]
		rawWeight := data["weight"]

		castName := rawName.(string)
		name := strings.Title(castName)
		imgUrl := sprites["front_default"].(string)
		imgShinyUrl := sprites["front_shiny"].(string)
		height := rawHeight.(float64) / 10
		weight := rawWeight.(float64) / 10

		return c.HTML(http.StatusOK, fmt.Sprintf(`	
		<div class="max-w-sm bg-white border border-gray-200 rounded-lg shadow dark:bg-gray-800 dark:border-gray-700 flex-1">
		<section class="inline-flex">
		<a href="#">
			<img class="rounded-t-lg w-40" src="%s" alt="normal pokemon" />
		</a>
		<a href="#">
			<img class="rounded-t-lg w-40" src="%s" alt="shiny pokemon" />
		</a>
		</section>
		<div class="p-5">
			<a href="#">
				<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Pokemon: %s</h5>
			</a>
			<p class="mb-3 font-normal text-gray-700 dark:text-gray-400">Weight: %v kg</p>
			<p class="mb-3 font-normal text-gray-700 dark:text-gray-400">Height: %v m</p>
			</div>
		</div>
		`, imgUrl, imgShinyUrl, name, weight, height))

	})

	// Start the server.
	e.Start(":8080")
}
