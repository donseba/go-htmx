package htmx

import (
	"bytes"
	"context"
	"html/template"
	"net/url"
	"path/filepath"
	"sync"
)

var (
	DefaultTemplateFuncs = template.FuncMap{}
	UseTemplateCache     = true
	templateCache        = sync.Map{} // Cache for parsed templates
)

type (
	RenderableComponent interface {
		Render(ctx context.Context) (template.HTML, error)
		Wrap(renderer RenderableComponent, target string) RenderableComponent
		With(r RenderableComponent, target string) RenderableComponent
		Attach(target string) RenderableComponent
		SetData(input map[string]any) RenderableComponent
		AddData(key string, value any) RenderableComponent
		SetGlobalData(input map[string]any) RenderableComponent
		AddGlobalData(key string, value any) RenderableComponent

		AddTemplateFunction(name string, function interface{}) RenderableComponent
		AddTemplateFunctions(funcs template.FuncMap) RenderableComponent

		SetURL(url *url.URL)

		data() map[string]any
		injectData(input map[string]any)
		addPartial(key string, value any)
		partials() map[string]RenderableComponent
		isWrapped() bool
		wrapper() RenderableComponent
		target() string
	}

	Component struct {
		templateData    map[string]any
		with            map[string]RenderableComponent
		partial         map[string]any
		globalData      map[string]any
		wrappedRenderer RenderableComponent
		wrappedTarget   string
		templates       []string
		url             *url.URL
		functions       template.FuncMap
	}
)

func NewComponent(templates ...string) *Component {
	return &Component{
		templateData: make(map[string]any),
		functions:    make(template.FuncMap),
		partial:      make(map[string]any),
		with:         make(map[string]RenderableComponent),
		templates:    templates,
	}
}

// Render renders the given templates with the given data
// it has all the default template functions and the additional template functions
// that are added with AddTemplateFunction
func (c *Component) Render(ctx context.Context) (template.HTML, error) {
	for key, value := range c.partials() {
		value.SetURL(c.url)
		value.injectData(c.templateData)

		ch, err := value.Render(ctx)
		if err != nil {
			return "", err
		}
		c.addPartial(key, ch)
	}

	//get the name of the first template file
	if len(c.templates) == 0 {
		return "", nil
	}

	return c.renderNamed(ctx, filepath.Base(c.templates[0]), c.templates, c.templateData)
}

// renderNamed renders the given templates with the given data
// it has all the default template functions and the additional template functions
// that are added with AddTemplateFunction
func (c *Component) renderNamed(ctx context.Context, name string, templates []string, input map[string]any) (template.HTML, error) {
	if len(templates) == 0 {
		return "", nil
	}

	var err error
	// Cache template parsing
	tmpl, cached := templateCache.Load(templates[0])
	if !cached || !UseTemplateCache {
		var functions = DefaultTemplateFuncs
		if c.functions != nil {
			for key, value := range c.functions {
				functions[key] = value
			}
		}
		tmpl, err = template.New(name).Funcs(functions).ParseFiles(templates...)
		if err != nil {
			return "", err
		}
		templateCache.Store(templates[0], tmpl)
	}

	data := struct {
		Ctx      context.Context
		Data     map[string]any
		Global   map[string]any
		Partials map[string]any
		URL      *url.URL
	}{
		Ctx:      ctx,
		Data:     input,
		Global:   c.globalData,
		Partials: c.partial,
		URL:      c.url,
	}

	var buf bytes.Buffer
	err = tmpl.(*template.Template).Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

// Wrap wraps the component with the given renderer
func (c *Component) Wrap(renderer RenderableComponent, target string) RenderableComponent {
	c.wrappedRenderer = renderer
	c.wrappedTarget = target

	return c
}

// With adds a partial to the component
func (c *Component) With(r RenderableComponent, target string) RenderableComponent {
	c.with[target] = r

	return c
}

// Attach adds a template to the main component but doesn't pre-render it
func (c *Component) Attach(target string) RenderableComponent {
	c.templates = append(c.templates, target)
	return c
}

func (c *Component) AddTemplateFunction(name string, function interface{}) RenderableComponent {
	c.functions[name] = function

	return c
}

func (c *Component) AddTemplateFunctions(funcs template.FuncMap) RenderableComponent {
	for key, value := range funcs {
		c.functions[key] = value
	}

	return c
}

func (c *Component) SetGlobalData(input map[string]any) RenderableComponent {
	c.globalData = input

	return c
}

func (c *Component) AddGlobalData(key string, value any) RenderableComponent {
	c.globalData[key] = value

	return c
}

// SetData adds data to the component
func (c *Component) SetData(input map[string]any) RenderableComponent {
	c.templateData = input

	return c
}

func (c *Component) AddData(key string, value any) RenderableComponent {
	c.templateData[key] = value

	return c
}

func (c *Component) SetURL(url *url.URL) {
	c.url = url
}

// isWrapped returns true if the component is wrapped
func (c *Component) isWrapped() bool {
	return c.wrappedRenderer != nil
}

// wrapper returns the wrapped renderer
func (c *Component) wrapper() RenderableComponent {
	return c.wrappedRenderer
}

// target returns the target
func (c *Component) target() string {
	return c.wrappedTarget
}

// partials returns the partials
func (c *Component) partials() map[string]RenderableComponent {
	return c.with
}

// injectData injects the input data into the template data
func (c *Component) injectData(input map[string]any) {
	for key, value := range input {
		if _, ok := c.templateData[key]; !ok {
			c.templateData[key] = value
		}
	}
}

func (c *Component) injectGlobalData(input map[string]any) {
	for key, value := range input {
		if _, ok := c.globalData[key]; !ok {
			c.globalData[key] = value
		}
	}
}

// addPartial adds a partial to the component
func (c *Component) addPartial(key string, value any) {
	c.partial[key] = value
}

// data returns the template data
func (c *Component) data() map[string]any {
	return c.templateData
}
