package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/donseba/go-htmx/internal/example"
	"github.com/donseba/go-htmx/pkg/htmx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type (
	App struct {
		logger htmx.Logger
		router chi.Router
		server htmx.Server
		htmx   *htmx.Service
	}
)

func main() {
	app := new(App)

	app.logger = log.New(os.Stdout, "go-htmx | ", 0)

	app.logger.Print("load env file")
	err := godotenv.Load(".env")
	if err != nil {
		app.logger.Fatal(errors.Wrap(err, "error loading .env file"))
	}

	app.logger.Print("start htmx service")
	app.htmx = htmx.NewService(app.logger, nil)

	app.logger.Print("start chi router")
	app.router = chi.NewRouter()
	app.router.Use(middleware.Logger)
	app.router.Use(middleware.Recoverer)
	app.router.Use(middleware.StripSlashes)
	app.router.Use(app.htmx.HxHeaderMiddleWare)

	app.server = &http.Server{
		Addr:    app.htmx.Config().ServerAddress,
		Handler: app.router,
	}

	app.router.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(http.Dir(
		filepath.Join(app.htmx.Root(), "assets"),
	))))

	app.router.Mount("/", example.NewController(app.htmx).Routes())

	app.logger.Printf("start server on : %s", app.htmx.Config().ServerAddress)

	err = app.server.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err)
	}
}
