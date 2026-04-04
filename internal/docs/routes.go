package docs

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//go:embed openapi.yaml
var openAPISpec []byte

type scalarConfig struct {
	URL string `json:"url"`
}

const scalarHTMLTemplate = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Expense Manager API Reference</title>
    <style>
      html, body, #app {
        height: 100%%;
        margin: 0;
      }
    </style>
  </head>
  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', %s)
    </script>
  </body>
</html>
`

// RegisterRoutes exposes the raw OpenAPI spec and a Scalar API reference page.
func RegisterRoutes(app *fiber.App) {
	app.Get("/openapi.yaml", serveOpenAPISpec)
	app.Get("/docs", serveScalarReference)
	app.Get("/docs/", func(c *fiber.Ctx) error {
		return c.Redirect("/docs", fiber.StatusMovedPermanently)
	})
}

func serveOpenAPISpec(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "application/yaml; charset=utf-8")
	return c.Send(openAPISpec)
}

func serveScalarReference(c *fiber.Ctx) error {
	config := scalarConfig{
		URL: "/openapi.yaml",
	}

	payload, err := json.Marshal(config)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	return c.SendString(sprintfScalarHTML(string(payload)))
}

func sprintfScalarHTML(config string) string {
	return fmt.Sprintf(scalarHTMLTemplate, config)
}
