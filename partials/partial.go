package partial

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"sync"
)

var (
	DefaultPartialHeader   = "Hx-Partial"
	DefaultTemplateFuncMap = template.FuncMap{}
	UseTemplateCache       = true
	templateCache          = sync.Map{}
)

type (
	Partial struct {
		id         string
		parent     *Partial
		oob        bool
		templates  []string
		functions  template.FuncMap
		data       map[string]any
		any        any
		globalData *GlobalData
		children   map[string]*Partial

		// internal data
		oobChildren map[string]struct{}
		partials    map[string]template.HTML
		wrapper     *Partial
		url         *url.URL
	}

	Data struct {
		Ctx      context.Context
		URL      *url.URL
		Data     map[string]any
		Global   map[string]any
		Any      any
		Partials map[string]template.HTML
	}

	GlobalData map[string]any
)

// New creates a new root.
func New(templates ...string) *Partial {
	return &Partial{
		id:          "root",
		templates:   templates,
		functions:   make(template.FuncMap),
		data:        make(map[string]any),
		globalData:  &GlobalData{},
		children:    make(map[string]*Partial),
		oobChildren: make(map[string]struct{}),
		partials:    make(map[string]template.HTML),
	}
}

// NewID creates a new instance with the provided ID.
func NewID(id string, templates ...string) *Partial {
	return New(templates...).ID(id)
}

func (p *Partial) ID(id string) *Partial {
	p.id = id
	return p
}

func (p *Partial) Templates(templates ...string) *Partial {
	p.templates = templates
	return p
}

func (p *Partial) Reset() *Partial {
	p.data = make(map[string]any)
	p.globalData = &GlobalData{}
	p.children = make(map[string]*Partial)
	p.oobChildren = make(map[string]struct{})
	p.partials = make(map[string]template.HTML)

	return p
}

// SetData sets the data for the partial.
func (p *Partial) SetData(data map[string]any) *Partial {
	p.data = data
	return p
}

// AddData adds data to the partial.
func (p *Partial) AddData(key string, value any) *Partial {
	p.data[key] = value
	return p
}

// SetGlobalData sets the global data for the partial.
func (p *Partial) SetGlobalData(data map[string]any) *Partial {
	*p.globalData = data
	return p
}

// AddGlobalData adds global data to the partial.
func (p *Partial) AddGlobalData(key string, value any) *Partial {
	(*p.globalData)[key] = value
	return p
}

// SetAny sets the any for the partial.
func (p *Partial) SetAny(any any) *Partial {
	p.any = any
	return p
}

// SetFuncs sets the functions for the partial.
func (p *Partial) SetFuncs(funcs template.FuncMap) *Partial {
	p.functions = funcs
	return p
}

// AddFunc adds a function to the partial.
func (p *Partial) AddFunc(name string, fn interface{}) *Partial {
	p.functions[name] = fn
	return p
}

// AppendFuncs appends functions to the partial if they do not exist.
func (p *Partial) AppendFuncs(funcs template.FuncMap) *Partial {
	for k, v := range funcs {
		if _, ok := p.functions[k]; !ok {
			p.functions[k] = v
		}
	}

	return p
}

// AddTemplate adds a template to the partial.
func (p *Partial) AddTemplate(template string) *Partial {
	p.templates = append(p.templates, template)
	return p
}

// With adds a child partial to the partial.
func (p *Partial) With(child *Partial) *Partial {
	p.children[child.id] = child
	p.children[child.id].globalData = p.globalData
	p.children[child.id].parent = p

	return p
}

// WithOOB adds an out-of-band child partial to the partial.
func (p *Partial) WithOOB(child *Partial) *Partial {
	p.With(child)
	p.oobChildren[child.id] = struct{}{}
	child.oob = true

	return p
}

// Wrap wraps the component with the given renderer
func (p *Partial) Wrap(renderer *Partial) *Partial {
	p.wrapper = renderer
	p.wrapper.With(p)

	return p
}

// RenderWithRequest renders the partial with the request.
func (p *Partial) RenderWithRequest(ctx context.Context, r *http.Request) (template.HTML, error) {
	p.url = r.URL
	var renderTarget, doRenderPartial = renderPartial(r)

	if renderTarget != "" {
		// render the partial with the request
		if c, ok := p.children[renderTarget]; ok {
			var (
				out template.HTML
				err error
			)
			c.AppendFuncs(p.functions)
			out, err = c.renderNamed(ctx, path.Base(c.templates[0]), c.templates)
			if err != nil {
				return "", err
			}

			// find all the oob children and add them to the output
			for id := range c.oobChildren {
				if child, cok := c.children[id]; cok {
					child.AppendFuncs(p.functions)
					if childData, err := child.renderNamed(ctx, path.Base(c.templates[0]), child.templates); err == nil {
						out += childData
					}
				}
			}

			return out, nil
		}

		if p.id != renderTarget {
			return "", fmt.Errorf("partial %s not found, got %s", renderTarget, p.id)
		}

	}

	// gather all children and render them into a map
	for id, child := range p.children {
		if child.oob {
			continue
		}

		if childData, err := child.RenderWithRequest(ctx, r); err == nil {
			p.partials[id] = childData
		} else {
			p.partials[id] = template.HTML(err.Error())
		}
	}

	if !doRenderPartial && p.wrapper != nil {
		parent := p.wrapper
		p.wrapper = nil

		return parent.RenderWithRequest(ctx, r)
	}

	return p.renderNamed(ctx, path.Base(p.templates[0]), p.templates)
}

// Render renders the partial.
func (p *Partial) renderNamed(ctx context.Context, name string, templates []string) (template.HTML, error) {
	if len(templates) == 0 {
		return "", nil
	}

	var err error
	functions := make(template.FuncMap)
	for key, value := range DefaultTemplateFuncMap {
		functions[key] = value
	}

	if p.functions != nil {
		for key, value := range p.functions {
			functions[key] = value
		}
	}

	cacheKey := generateCacheKey(templates, functions)
	tmpl, cached := templateCache.Load(cacheKey)
	if !cached || !UseTemplateCache {
		// Parse and cache template as before
		tmpl, err = template.New(name).Funcs(functions).ParseFiles(templates...)
		if err != nil {
			return "", err
		}
		templateCache.Store(cacheKey, tmpl)
	}

	data := &Data{
		URL:      p.url,
		Ctx:      ctx,
		Any:      p.any,
		Data:     p.data,
		Global:   *p.globalData,
		Partials: p.partials,
	}

	if t, ok := tmpl.(*template.Template); ok {
		var buf bytes.Buffer
		err = t.Execute(&buf, data)
		if err != nil {
			return "", err
		}

		return template.HTML(buf.String()), nil // Return rendered content
	}

	return "", errors.New("template is not a *template.Template")
}

func renderPartial(r *http.Request) (string, bool) {
	hxRequest := r.Header.Get("Hx-Request")
	hxBoosted := r.Header.Get("Hx-Boosted")
	hxHistoryRestoreRequest := r.Header.Get("Hx-History-Restore-Request")

	return r.Header.Get(DefaultPartialHeader), (hxRequest == "true" || hxBoosted == "true") && hxHistoryRestoreRequest != "true"
}

// Generate a hash of the function names to include in the cache key
func generateCacheKey(templates []string, funcMap template.FuncMap) string {
	var funcNames []string
	for name := range funcMap {
		funcNames = append(funcNames, name)
	}
	// Sort function names to ensure consistent ordering
	sort.Strings(funcNames)
	hash := sha256.Sum256([]byte(strings.Join(funcNames, ",")))
	return templates[0] + ":" + hex.EncodeToString(hash[:])
}
