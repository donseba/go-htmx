package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/donseba/go-htmx"
	htmx_chi "github.com/donseba/go-htmx/middleware/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

type (
	App struct {
		server *http.Server
		logger htmx.Logger
		router chi.Router
		htmx   *htmx.Service
	}
)

func main() {
	app := new(App)
	app.logger = log.New(os.Stdout, "go-htmx | ", 0)

	app.logger.Print("start htmx service")
	var err error
	app.htmx, err = htmx.NewService(&htmx.Config{
		ServerAddress:      "localhost:8888",
		TemplateDir:        "templates",
		TemplateFuncs:      nil,
		ErrorTemplate:      filepath.Join("error.gohtml"),
		DefaultTemplates:   []string{filepath.Join("index.gohtml")},
		DefaultTemplatesHx: []string{filepath.Join("hx", "index.gohtml")},
		Logger:             app.logger,
	})
	if err != nil {
		app.logger.Fatal(errors.Wrap(err, "error loading .env file"))
	}

	app.logger.Print("start chi router")
	app.router = chi.NewRouter()
	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)
	app.router.Use(middleware.StripSlashes)
	app.router.Use(htmx_chi.MiddleWare)

	app.server = &http.Server{
		Addr:    app.htmx.Config().ServerAddress,
		Handler: app.router,
	}

	app.router.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir(
		filepath.Join(app.htmx.Root(), "assets"),
	))))

	app.router.Mount("/", NewController(app.htmx).Routes())

	app.logger.Printf("start server on : %s", app.htmx.Config().ServerAddress)

	err = app.server.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err)
	}
}
