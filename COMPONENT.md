# go-htmx Component Documentation

The `go-htmx` package provides a flexible and efficient way to render partial or full HTML pages in Go applications. It leverages Go's `html/template` package to render templates with dynamic data, supporting features like partial rendering, template wrapping, and data injection.

---

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Creating Components](#creating-components)
4. [Rendering Components](#rendering-components)
5. [Wrapping Components](#wrapping-components)
6. [Adding Partials](#adding-partials)
7. [Attaching Templates](#attaching-templates)
8. [Working with Data](#working-with-data)
9. [Template Functions](#template-functions)
10. [Reusing Components](#reusing-components)
11. [Caveats and Warnings](#caveats-and-warnings)
12. [Example Usage](#example-usage)
13. [Configuration Options](#configuration-options)
14. [Conclusion](#conclusion)
15. [Additional Notes](#additional-notes)
16. [Internal Details](#internal-details)
17. [Caveats and Warnings (Detailed)](#caveats-and-warnings-detailed)
18. [Feedback and Contributions](#feedback-and-contributions)

---

## Introduction

The `go-htmx` package simplifies the process of rendering HTML templates in Go by introducing the concept of components. Components encapsulate templates, data, and rendering logic, making it easier to build complex web pages with reusable parts.

---

## Getting Started

To use the `go-htmx` package, you need to import it into your Go project:

```go
import "github.com/donseba/go-htmx"
```

Ensure you have the package installed:

```bash
go get github.com/donseba/go-htmx
```

---

## Creating Components
A component represents a renderable unit, typically associated with one or more template files. You can create a new component using the `NewComponent` function:

```go
component := htmx.NewComponent("templates/base.html")
```

You can specify multiple template files if needed:

```go
component := htmx.NewComponent("templates/header.html", "templates/body.html", "templates/footer.html")
```

---

## Rendering Components
To render a component, you call its `Render` method, passing in a `context.Context`:

```go
htmlContent, err := component.Render(ctx)
if err != nil {
    // Handle error
}
```

The `Render` method processes the templates and returns the rendered HTML content as a `template.HTML` type.

---

## Wrapping Components
Components can be wrapped inside other components, allowing you to nest templates and create complex layouts. Use the Wrap method to wrap a component:

```go
wrapperComponent := htmx.NewComponent("templates/wrapper.html")
component.Wrap(wrapperComponent, "content")
```

In the wrapper template, you can define a placeholder (e.g., `{{ .Partials.content }})` where the wrapped component's content will be inserted.

--- 

## Adding Partials
Partials are sub-components that can be included within a component. Use the With method to add a partial:
```go
partialComponent := htmx.NewComponent("templates/partial.html")
component.With(partialComponent, "sidebar")
```

In your main component's template, you can reference the partial using `{{ .Partials.sidebar }}`.

--- 

## Attaching Templates
If you have additional templates that you want to include without rendering them as partials, you can use the Attach method:
```go
component.Attach("templates/extra.html")
```
This method appends the template to the component's template list.

## Working with Data
You can pass dynamic data to your templates using the SetData and AddData methods.

### Setting Data

Set multiple data values at once:
```go
data := map[string]interface{}{
    "Title":   "Welcome Page",
    "Message": "Hello, World!",
}
component.SetData(data)
```

### Adding Data
Add individual data values:
```go
component.AddData("Title", "Welcome Page")
component.AddData("Message", "Hello, World!")
``` 

### Global Data
Set data that is accessible to all components and partials using `SetGlobalData` and `AddGlobalData`:
```go
component.SetGlobalData(map[string]interface{}{
    "AppName": "My Go App",
})

component.AddGlobalData("Version", "1.0.0")
```

--- 

## Template Functions
You can enhance your templates with custom functions using `AddTemplateFunction` and `AddTemplateFunctions`.

### Adding a Single Function
```go 
component.AddTemplateFunction("formatDate", func(t time.Time) string {
    return t.Format("Jan 2, 2006")
})
```

Use `{{ formatDate .Data.Timestamp }}` in your template.

### Adding Multiple Functions
```go
funcMap := template.FuncMap{
    "toUpper": strings.ToUpper,
    "safeHTML": func(s string) template.HTML {
        return template.HTML(s)
    },
}

component.AddTemplateFunctions(funcMap)
```

--- 

## Reusing Components
If you need to reuse a Component instance, be aware that internal state (like data and partials) may persist between renders. To reset the component's state, use the Reset method:
```go
component.Reset()
```

This method clears the component's data, global data, partials, and URL, allowing you to reuse it without residual state.

--- 

## Caveats and Warnings

### Thread Safety

Warning: Component instances are not thread-safe and should not be shared across goroutines. If you need to render components concurrently, create separate instances for each goroutine.
Reusing Components

- **State Persistence**: When reusing a component, internal state such as data and partials may persist. Always call Reset() before reusing a component to avoid unintended data leakage.
- **Resetting Components**: The Reset method clears data, global data, partials, and the URL. It does not reset templates or functions.

### URL Propagation

- **Setting the URL**: When you set the URL on a component using SetURL, it is recursively propagated to all partials, including nested ones.
- **Adding Partials After Setting URL**: If you add partials after setting the URL, you may need to call SetURL again to ensure the new partials receive the URL.

### Data Overwriting in injectData

- **Non-Overwriting Behavior**: The injectData method does not overwrite existing keys in templateData. If a key already exists, it will not be replaced.
- **Recommendation**: Be mindful of this behavior when injecting data to avoid unexpected results.

--- 

## Example Usage

Here's a complete example demonstrating how to use the go-htmx package:

```go
package main

import (
	"context"
	"fmt"
	"github.com/donseba/go-htmx"
	"net/http"
)

type (
	App struct {
		htmx *htmx.HTMX
	}
)

func (a *App) handler(w http.ResponseWriter, r *http.Request) {
	// Create main component
	mainComponent := htmx.NewComponent("templates/main.html")

	// Set data
	mainComponent.SetData(map[string]interface{}{
		"Title": "Home Page",
	})

	// Add a partial
	headerComponent := htmx.NewComponent("templates/header.html")
	mainComponent.With(headerComponent, "header")

	// Set URL (propagates to all partials)
	mainComponent.SetURL(r.URL)

	h := a.htmx.NewHandler(w, r)

	_, err = h.Render(r.Context(), mainComponent)
	if err != nil {
		fmt.Printf("error rendering page: %v", err.Error())
	}
}

func main() {
	app := &App{
		htmx: htmx.New(),
	}

	http.HandleFunc("/", app.handler)
	http.ListenAndServe(":8080", nil)
}
```

--- 

## Configuration Options

### Template Functions
The package provides a default function map `DefaultTemplateFuncs` that you can populate with common functions.
```go 
htmx.DefaultTemplateFuncs = template.FuncMap{
    "toUpper": strings.ToUpper,
}
```

### Template Caching
Templates are cached by default for performance. You can control this behavior using the `UseTemplateCache` variable:
```go
htmx.UseTemplateCache = false // Disable template caching
```

--- 

## Conclusion
The Component addition offers a powerful way to manage and render templates in Go applications. By structuring your templates into components and partials, and by leveraging data injection and custom template functions, you can build dynamic and maintainable web pages.


--- 

## Additional Notes 

- Context in Templates: The context passed to Render is available in templates as `{{ .Ctx }}`.
- Data Access in Templates
  - **Accessing Data**: Use `{{ .Data.Key }}` to access data values in templates.
  - **Global Data**: Global data is accessible as `{{ .Global.Key }}` in templates.
  - **Partials**: Partials are available as `{{ .Partials.Key }}` in templates.
  - **URL**: The URL is accessible as `{{ .URL }}` in templates.

--- 

## Internal Details (For Advanced Users)

### Data Structures
- **Component**: Implements the `RenderableComponent` interface and holds all the necessary information to render templates, including data, partials, and template functions.

### Rendering Process
1. **Partial Rendering**: Before rendering the main template, any partial components added via `With` are rendered, and their output is stored.
2. **Template Parsing**: Templates are parsed and cached (if caching is enabled).
3. **Data Preparation**: A data structure containing context, data, global data, partials, and URL is prepared.
4. **Execution**: The template is executed with the prepared data.

### Important Methods
- `Render(ctx context.Context) (template.HTML, error)`: Renders the component.
- `Wrap(renderer RenderableComponent, target string) RenderableComponent`: Wraps the component with another renderer.
- `With(r RenderableComponent, target string) RenderableComponent`: Adds a partial component.
- `SetData(input map[string]interface{}) RenderableComponent`: Sets the template data.
- `AddTemplateFunction(name string, function interface{}) RenderableComponent`: Adds a custom template function.
- `Reset() *Component`: Resets the component's state.

--- 

## Caveats and Warnings (Detailed)

### Thread Safety
Important: `Component` instances are not thread-safe. Do not share a single Component instance across multiple goroutines. Each goroutine should create its own instance of Component to avoid race conditions and undefined behavior.
There is currently no real usage so far for using Components across multiple goroutines, if the need arises, we can discuss and implement a thread-safe version of the Component.

### Reusing Components
- **State Persistence**: The `Component` retains its state between renders. If you reuse a Component, data, partials, and other settings from previous renders may persist.
- **Using Reset**: To reuse a `Component` safely, call `Reset()` to clear its state before setting new data or partials.

### Data Injection Behavior
- **Non-Overwriting in `injectData`**: The `injectData` method does not overwrite existing keys in templateData. This means that if a key exists in both the component's data and the injected data, the component's data takes precedence.
- **Best Practice**: Be explicit with your data keys and manage them carefully to prevent unexpected behavior.

### URL Handling
- **Propagation to Partials**: When you set the URL on a component using `SetURL`, it automatically propagates to all its partials, including nested ones.
- **Adding Partials After URL Is Set**: If you add partials after setting the URL, you need to call `SetURL` again to ensure the new partials receive the URL.

### Template Caching
- **Function Map Consideration**: The template caching mechanism accounts for custom function maps. Templates with different function maps are cached separately.
- **Disabling Cache**: You can disable template caching during development or debugging by setting `UseTemplateCache` to `false`.

--- 

## Feedback and Contributions
We welcome feedback, suggestions, and contributions to the `go-htmx` package. If you have ideas for improvements, new features, or bug fixes, please open an issue or submit a pull request on the GitHub repository.