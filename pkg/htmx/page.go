package htmx

import (
	"context"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
)

type (
	TemplateData struct {
		Context   func() context.Context
		PageTitle string
		BaseURL   string
		Data      map[string]interface{}
	}
)

func (s *Service) TemplateData(ctx context.Context) (*TemplateData, error) {
	app := &TemplateData{
		Context: func() context.Context {
			return ctx
		},
		PageTitle: "go-htmx",
		Data:      make(map[string]interface{}),
		BaseURL:   s.Config().ServerAddress,
	}

	return app, nil
}

func (s *Service) Template(r *http.Request) (*template.Template, error) {
	tmpl, err := s.tmpl.Clone()
	if err != nil {
		return nil, err
	}

	hxh := s.HxHeader(r.Context())
	if hxh.HxRequest {
		var tmp []string

		for i := 0; i < len(s.Config().DefaultTemplatesHx); i++ {
			tmp = append(tmp, filepath.Join(s.rootDir, s.Config().TemplateDir, s.Config().DefaultTemplatesHx[i]))
		}

		return tmpl.ParseFiles(tmp...)
	}

	var tmp []string
	for i := 0; i < len(s.Config().DefaultTemplates); i++ {
		tmp = append(tmp, filepath.Join(s.rootDir, s.Config().TemplateDir, s.Config().DefaultTemplates[i]))
	}

	return tmpl.ParseFiles(tmp...)
}

func (s *Service) Error(w http.ResponseWriter, r *http.Request, inErr error, statusCode int) {
	ctx := r.Context()

	td, err := s.TemplateData(ctx)
	if err != nil {
		s.logger.Print(err)
		return
	}
	td.Data["Error"] = inErr
	td.Data["StatusCode"] = statusCode
	td.Data["Alert"] = false

	hxh := s.HxHeader(ctx)
	if hxh.HxRequest == false {
		w.WriteHeader(statusCode)
	} else {
		td.Data["Alert"] = true
	}

	s.serveError(td).ServeHTTP(w, r)
}

func (s *Service) serveError(td *TemplateData) http.HandlerFunc {
	var (
		tmpl *template.Template
		err  error
	)

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err = s.Template(r)
		if err != nil {
			s.logger.Print("msg", "template status error", "err", err)
			_, _ = w.Write([]byte(errors.Wrap(err, "error serving error").Error()))
			return
		}

		tmpl, err = tmpl.ParseFiles(
			filepath.Join(s.rootDir, s.config.TemplateDir, "error.gohtml"),
		)

		if err != nil {
			s.logger.Print("msg", "template status error", "err", err)
			_, _ = w.Write([]byte(errors.Wrap(err, "error serving error").Error()))
			return
		}

		err = tmpl.ExecuteTemplate(w, "index.gohtml", td)
		if err != nil {
			s.logger.Print("msg", "error serving error", "err", err)
			_, _ = w.Write([]byte(errors.Wrap(err, "error serving error").Error()))
			return
		}
	}
}

func (s *Service) Serve(templateData *TemplateData, templates []string) http.HandlerFunc {
	var (
		tmpl *template.Template
		err  error
	)

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err = s.Template(r)
		if err != nil {
			s.Error(w, r, err, http.StatusInternalServerError)
			return
		}

		var temps = make([]string, len(templates))
		for i := 0; i < len(templates); i++ {
			temps[i] = filepath.Join(s.rootDir, s.Config().TemplateDir, templates[i])
		}

		tmpl, err = tmpl.ParseFiles(temps...)
		if err != nil {
			s.Error(w, r, err, http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "index.gohtml", templateData)
		if err != nil {
			s.Error(w, r, err, http.StatusInternalServerError)
			return
		}
	}
}
