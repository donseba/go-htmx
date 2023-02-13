package htmx

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type (
	Service struct {
		config  *Config
		appDir  string
		baseURL *url.URL
		tmpl    *template.Template
	}

	Config struct {
		ServerAddress      string
		TemplateDir        string
		TemplateFuncs      map[string]interface{}
		ErrorTemplate      string
		DefaultTemplates   []string
		DefaultTemplatesHx []string
		Logger             Logger
	}

	Logger interface {
		Print(v ...any)
		Printf(format string, v ...any)
		Fatal(v ...any)
	}
)

func NewServiceWithDefaults() (*Service, error) {
	config := &Config{
		ServerAddress:      "localhost:8888",
		TemplateDir:        "templates",
		TemplateFuncs:      nil,
		DefaultTemplates:   []string{filepath.Join("index.gohtml")},
		DefaultTemplatesHx: []string{filepath.Join("hx", "index.gohtml")},
	}

	return NewService(config)
}

func NewService(config *Config) (*Service, error) {
	if config.Logger == nil {
		config.Logger = log.New(os.Stdout, "go-htmx | ", 0)
	}

	s := &Service{
		config: config,
	}

	// base template with base functions.
	s.tmpl = template.New("htmx")
	s.tmpl = s.tmpl.Funcs(config.TemplateFuncs)

	// get executable path
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	s.appDir = filepath.Dir(ex)

	s.baseURL, err = url.Parse(s.config.ServerAddress)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) Redirect(w http.ResponseWriter, r *http.Request, url string, statusCode int) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Push", "true")
		w.Header().Set("HX-Redirect", "true")
	}

	http.Redirect(w, r, url, statusCode)
}

func (s *Service) Logger() Logger {
	return s.config.Logger
}

func (s *Service) Config() *Config {
	return s.config
}

func (s *Service) Root() string {
	return s.appDir
}
