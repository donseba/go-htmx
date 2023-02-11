package htmx

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type (
	Service struct {
		logger        Logger
		Router        chi.Router
		config        *Config
		rootDir       string
		tmpl          *template.Template
		TemplateFuncs map[string]interface{}
	}

	Logger interface {
		Print(v ...any)
		Printf(format string, v ...any)
		Fatal(v ...any)
	}

	Server interface {
		ListenAndServe() error
	}

	Pager interface {
		Template(r *http.Request) (*template.Template, error)
		Serve(templateData *TemplateData, templates []string) http.HandlerFunc
		Error(w http.ResponseWriter, r *http.Request, err error, statusCode int)
	}
)

func NewService(logger Logger, templateFuncs map[string]interface{}) *Service {
	if logger == nil {
		logger = log.New(os.Stdout, "go-htmx | ", 0)
	}

	app := &Service{
		logger: logger,
		config: parseConfig(),
	}

	// base template with base functions
	app.tmpl = template.New("htmx")
	app.tmpl = app.tmpl.Funcs(templateFuncs)

	// get executable path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	app.rootDir = filepath.Dir(ex)

	return app
}

func (s *Service) Redirect(w http.ResponseWriter, r *http.Request, url string, statusCode int) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Push", "true")
		w.Header().Set("HX-Redirect", "true")
	}

	http.Redirect(w, r, url, statusCode)
}

func (s *Service) Logger() Logger {
	return s.logger
}

func (s *Service) Config() *Config {
	return s.config
}

func (s *Service) Root() string {
	return s.rootDir
}
