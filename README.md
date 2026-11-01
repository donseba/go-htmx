# go-htmx
**Seamless HTMX integration in golang applications.**

[![GoDoc](https://pkg.go.dev/badge/github.com/donseba/go-htmx?status.svg)](https://pkg.go.dev/github.com/donseba/go-htmx?tab=doc)
[![GoMod](https://img.shields.io/github/go-mod/go-version/donseba/go-htmx)](https://github.com/donseba/go-htmx)
[![Size](https://img.shields.io/github/languages/code-size/donseba/go-htmx)](https://github.com/donseba/go-htmx)
[![License](https://img.shields.io/github/license/donseba/go-htmx)](./LICENSE)
[![Stars](https://img.shields.io/github/stars/donseba/go-htmx)](https://github.com/donseba/go-htmx/stargazers)
[![Go Report Card](https://goreportcard.com/badge/github.com/donseba/go-htmx)](https://goreportcard.com/report/github.com/donseba/go-htmx)

## Description

This repository contains the htmx Go package, designed to enhance server-side handling of HTML generated with the [HTMX library](https://htmx.org/). 
It provides a set of tools to easily manage swap behaviors, trigger configurations, and other HTMX-related functionalities in a Go server environment.

## Disclaimer
This package is built around the specific need to be able to work with HTMX in a Go environment.
All functionality found in this repository has a certain use case in various projects that I have worked on.

- Design decisions are documented in the [DESIGN_DECISIONS.md](https://github.com/donseba/go-htmx/blob/main/DESIGN_DECISIONS.md) file.
- Locality of Behavior is documented in the [LOB.md](https://github.com/donseba/go-htmx/blob/main/LOB.md) file.

## Features

- **Component Rendering**: Render (partial) components in response to HTMX requests, enhancing user experience and performance.
- **Swap Configuration**: Configure swap behaviors for HTMX responses, including style, timing, and scrolling.
- **Trigger Management**: Define and manage triggers for HTMX events, supporting both simple and detailed triggers.
- **Middleware Support**: Integrate HTMX seamlessly with Go middleware for easy HTMX header configuration.
- **io.Writer Support**: The HTMX handler implements the io.Writer interface for easy integration with existing Go code.
- **htmx 2 and htmx 4**: htmx 2 by default, htmx 4 with one setting, see [htmx 4](#htmx-4).

---

## Getting Started

### Installation

To install the htmx package, use the following command:

```sh
go get -u github.com/donseba/go-htmx
```

### Usage

initialize the htmx service like so : 
```go
package main

import (
	"log"
	"net/http"

	"github.com/donseba/go-htmx"
)

type App struct {
	htmx *htmx.HTMX
}

func main() {
	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	mux := http.NewServeMux()
	// wrap the htmx example middleware around the http handler
	mux.HandleFunc("/", app.Home)

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	// initiate a new htmx handler
	h := a.htmx.NewHandler(w, r)

	// check if the request is a htmx request
	if h.IsHxRequest() {
		// do something
	}
	
	// check if the request is boosted
	if h.IsHxBoosted() {
		// do something
	}
	
	// check if the request is a history restore request
	if h.IsHxHistoryRestoreRequest() { 
		// do something 
	}
	
	// check if the request is a prompt request
	if h.RenderPartial() { 
		// do something
	}
		
	// set the headers for the response, see docs for more options
	h.PushURL("http://push.url")
	h.ReTarget("#ReTarged")

	// write the output like you normally do.
	// check the inspector tool in the browser to see that the headers are set.
	_, _ = h.Write([]byte("OK"))
}
```

### HTMX Request Checks

The htmx package provides several functions to determine the nature of HTMX requests in your Go application. These checks allow you to tailor the server's response based on specific HTMX-related conditions.

#### IsHxRequest

This function checks if the incoming HTTP request is made by HTMX.

```go
func (h *Handler) IsHxRequest() bool
```
- **Usage**: Use this check to identify requests initiated by HTMX and differentiate them from standard HTTP requests.
- **Example**: Applying special handling or returning partial HTML snippets in response to an HTMX request.

#### IsHxBoosted

Determines if the HTMX request is boosted, which typically indicates an enhancement of the user experience with HTMX's AJAX capabilities.

```go
func (h *Handler) IsHxBoosted() bool
```
- **Usage**: Useful in scenarios where you want to provide an enriched or different response for boosted requests.
- **Example**: Loading additional data or scripts that are specifically meant for AJAX-enhanced browsing.

#### IsHxHistoryRestoreRequest

Checks if the HTMX request is a history restore request. This type of request occurs when HTMX is restoring content from the browser's history.

```go
func (h *Handler) IsHxHistoryRestoreRequest() bool
```
- **Usage**: Helps in handling scenarios where users navigate using browser history, and the application needs to restore previous states or content.
- **Example**: Resetting certain states or re-fetching data that was previously displayed.

#### RenderPartial

This function returns true for HTMX requests that are either standard or boosted, as long as they are not history restore requests. It is a combined check used to determine if a partial render is appropriate.

```go
func (h *Handler) RenderPartial() bool
```
- **Usage**: Ideal for deciding when to render partial HTML content, which is a common pattern in applications using HTMX.
- **Example**: Returning only the necessary HTML fragments to update a part of the webpage, instead of rendering the entire page.

### Swapping
Swapping is a way to replace the content of a dom element with the content of the response.
This is done by setting the `HX-Swap` header to the id of the dom element you want to swap.

```go
func (c *Controller) Route(w http.ResponseWriter, r *http.Request) {
	// initiate a new htmx handler 
	h := a.htmx.NewHandler(w, r)
	
	// Example usage of Swap 
	swap := htmx.NewSwap().Swap(time.Second * 2).ScrollBottom() 
	
	h.ReSwapWithObject(swap)
	
	_, _ = h.Write([]byte("your content"))
}
```

### Trigger Events 
Trigger events are a way to trigger events on the dom element.
This is done by setting the `HX-Trigger` header to the event you want to trigger.

```go
func (c *Controller) Route(w http.ResponseWriter, r *http.Request) {
	// initiate a new htmx handler 
	h := a.htmx.NewHandler(w, r)
	
	// Example usage of Swap 
	trigger := htmx.NewTrigger().AddEvent("event1").AddEventDetailed("event2", "Hello, World!") 
	
	h.TriggerWithObject(trigger)
	// or 
	h.TriggerAfterSettleWithObject(trigger)
	// or
	h.TriggerAfterSwapWithObject(trigger)
	
	_, _ = h.Write([]byte("your content"))
}
```

---

## htmx 4

go-htmx speaks htmx 2 by default. For [htmx 4](https://four.htmx.org/), set the version on the htmx instance:

```go
app := &App{
	htmx: htmx.New(),
}

app.htmx.SetVersion(htmx.Version4)
```

Handlers created after `SetVersion` read the request headers and write the response headers the way htmx 4 does. Without it, nothing changes.

```html
<script src="https://unpkg.com/htmx.org@4.0.0/dist/htmx.min.js"></script>
```

### What changes with htmx 4

| Method | htmx 2 | htmx 4 (`Version4`) |
| --- | --- | --- |
| `h.RenderPartial()` | true for an htmx request that is no history restore | also false for `HX-Request-Type: full`, which htmx 4 sends for body swaps, `hx-select` and history restoration |
| `h.TargetID()` | `HX-Target`, the id of the target | the decoded id of `HX-Target`, which htmx 4 sends as `tagName#id`: `div#content` gives `content` |
| `h.SourceID()` | `HX-Trigger`, the id of the element that triggered the request | the decoded id of `HX-Source`: `button#save` gives `save` |
| `h.Prompt()` | `HX-Prompt` as sent | `HX-Prompt` decoded, as the `hx-prompt` extension of htmx 4 sends it URI encoded |
| `h.ReSwapWithObject(swap)` | `show:#list:top focus-scroll:true` | `show:top showTarget:"#list" focusScroll:true` |
| `h.TriggerAfterSwap*`, `h.TriggerAfterSettle*` | `HX-Trigger-After-Swap`, `HX-Trigger-After-Settle` | nothing, as htmx 4 removed both headers |
| `h.StopPolling()` | status 286, which stops polling | status 286, which htmx 4 swaps like a 200 while it keeps polling |

`h.Request()` holds the request headers as sent, for both versions, now also `HxRequestType` (`HX-Request-Type`) and `HxSource` (`HX-Source`). htmx 4 no longer sends `HX-Trigger` and `HX-Trigger-Name`, and sends `HX-Prompt` only with its `hx-prompt` extension.

The other response headers are the same in htmx 4. `Location` works too; htmx 4 does not use `LocationInput.Handler`.

### Triggers

htmx 4 triggers the events of `HX-Trigger` once the response is swapped and settled, so `h.Trigger` and `h.TriggerWithObject` cover most uses of `TriggerAfterSwap` and `TriggerAfterSettle`. For exact timing, listen to `htmx:after:swap` or `htmx:after:settle` in the browser. The notifications (`TriggerSuccess`, `TriggerError`, …) use `HX-Trigger` and work with both versions.

### Polling

htmx 4 has no status code that stops polling: it stops when the polling element leaves the page. Answer the last poll with markup that replaces the polling element, without its `every` trigger:

```go
func (a *App) Progress(w http.ResponseWriter, r *http.Request) {
	h := a.htmx.NewHandler(w, r)

	if done {
		// <div id="progress" hx-get="/progress" hx-trigger="every 1s"> is replaced, so htmx 4 stops polling
		h.ReTarget("#progress")
		h.ReSwap("outerHTML")
		h.JustWriteString(`<div id="progress">Done</div>`)
		return
	}

	h.JustWriteString("Working…")
}
```

### Swap syntax

`Swap.String()` keeps the htmx 2 syntax. `Swap.StringFor` writes the syntax of a version, for an `hx-swap` attribute in a template or anywhere without a handler:

```go
swap := htmx.NewSwap().ScrollBottom("#messages")

swap.StringFor(htmx.Version4)      // innerHTML scroll:bottom scrollTarget:"#messages"
swap.StringFor(app.htmx.Version()) // the version of the htmx instance
```

```html
<div hx-get="/messages" hx-swap="{{ .Swap.StringFor 4 }}"></div>
```

htmx 4 adds the swap styles `htmx.SwapInnerMorph`, `htmx.SwapOuterMorph`, `htmx.SwapOuterSync` and `htmx.SwapTextContent`. A string passed to `h.ReSwap` is sent as given, so write it in the syntax of htmx 4.

### Error responses

htmx 4 swaps every response except 204 and 304, also 4xx and 5xx responses, which htmx 2 did not swap. Write error responses that work as swap content, or control it per status in the page with `hx-status`, such as `hx-status:5xx="swap:none"`.

### Request checks without a handler

The package functions `htmx.IsHxRequest(r)`, `htmx.IsHxBoosted(r)`, `htmx.IsHxHistoryRestoreRequest(r)` and `htmx.RenderPartial(r)`, and the middleware, have no version: they keep the htmx 2 behaviour. With htmx 4, use the handler of the htmx instance, `app.htmx.NewHandler(w, r).RenderPartial()`.

### Server Sent Events

The `sse` package works with the `hx-sse` extension of htmx 4, which replaces `hx-ext="sse"`, `sse-connect` and `sse-swap`:

```html
<script src="https://unpkg.com/htmx.org@4.0.0/dist/htmx.min.js"></script>
<script src="https://unpkg.com/htmx.org@4.0.0/dist/ext/hx-sse.min.js"></script>

<div hx-sse:connect="/sse"></div>
```

htmx 4 swaps a message without an event name into the element, and triggers a message with an event name as an event on the element. To update several elements from one stream, send messages without an event name, with out-of-band elements:

```go
sseManager.Send(sse.NewMessage(`<div id="time" hx-swap-oob="true">12:00</div>`))
```

### In the browser

htmx 4 changes more in the page than in the headers: attributes are inherited only with `:inherited`, such as `hx-target:inherited="#content"`, and some attributes and events are renamed. See [What's new in htmx 4](https://four.htmx.org/docs/whats-new-in-htmx-4). The [examples](examples) use htmx 4.

---

## utility methods 

### Notification handling 
Comprehensive support for triggering various types of notifications within your Go applications, enhancing user interaction and feedback. The package provides a set of functions to easily manage and trigger different notification types such as success, info, warning, error, and custom notifications.
Available Notification Types

- **Success**: Use for positive confirmation messages.
- **Info**: Ideal for informational messages.
- **Warning**: Suitable for cautionary messages.
- **Error**: Use for error or failure messages.
- **Custom**: Allows for defining your own notification types.

### Usage

Triggering notifications is straightforward. Here are some examples demonstrating how to use each function:

```go
func (h *Handler) MyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Trigger a success notification 
	h.TriggerSuccess("Operation completed successfully")
	
	// Trigger an info notification 
	h.TriggerInfo("This is an informational message")

	// Trigger a warning notification 
	h.TriggerWarning("Warning: Please check your input")
	
	// Trigger an error notification 
	h.TriggerError("Error: Unable to process your request")
	
	// Trigger a custom notification 
	h.TriggerCustom("customType", "This is a custom notification", nil)
}
```

### Notification Levels

The htmx package provides built-in support for four primary notification levels, each representing a different type of message:

- `success`: Indicates successful completion of an operation.
- `info`: Conveys informational messages.
- `warning`: Alerts about potential issues or cautionary information.
- `error`: Signals an error or problem that occurred.

Each notification type is designed to communicate specific kinds of messages clearly and effectively in your application's user interface.
### Triggering Custom Notifications

In addition to these standard notification levels, the htmx package also allows for custom notifications using the TriggerCustom method. This method provides the flexibility to define a custom level and message, catering to unique notification requirements.

```go
func (h *Handler) MyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Trigger standard notifications 
	h.TriggerSuccess("Operation successful")
	h.TriggerInfo("This is for your information")
	h.TriggerWarning("Please be cautious")
	h.TriggerError("An error has occurred")
	
	// Trigger a custom notification 
	h.TriggerCustom("customLevel", "This is a custom notification")
}
```
The TriggerCustom method enables you to specify a custom level (e.g., "customLevel") and an accompanying message. This method is particularly useful when you need to go beyond the predefined notification types and implement a notification system that aligns closely with your application's specific context or branding.

### Advanced Usage with Custom Variables

You can also pass additional data with your notifications. Here's an example:

```go
func (h *Handler) MyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	customData := map[string]string{"key1": "value1", "key2": "value2"}
	h.TriggerInfo("User logged in", customData)
}
```

### the HTMX part 

please refer to the [htmx documentation](https://htmx.org/headers/hx-trigger/) regarding event triggering. and the example [confirmation UI](https://htmx.org/examples/confirm/)

`HX-Trigger: {"showMessage":{"level" : "info", "message" : "Here Is A Message"}}`

And handle this event like so:

```js 
document.body.addEventListener("showMessage", function(evt){
    if(evt.detail.level === "info"){
        alert(evt.detail.message);
    }
})
```
Each property of the JSON object on the right hand side will be copied onto the details object for the event.

### Customizing Notification Event Names

In addition to the standard notification types, the htmx package allows you to customize the event name used for triggering notifications. This is done by modifying the htmx.DefaultNotificationKey. Changing this key will affect the event name in the HTMX trigger, allowing you to tailor it to specific needs or naming conventions of your application.
Setting a Custom Notification Key

Before triggering notifications, you can set a custom event name as follows:

```go
htmx.DefaultNotificationKey = "myCustomEventName"
```

---

## Component Rendering

The components documentation can be found in the [COMPONENTS.md](https://github.com/donseba/go-htmx/blob/main/COMPONENTS.md) file.

---

## Middleware
The htmx package is designed for versatile integration into Go applications, providing support both with and without the use of middleware. Below, we showcase two examples demonstrating the package's usage in scenarios involving middleware.

### standard mux middleware example:

```go
func MiddleWare(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        hxh := htmx.HxRequestHeaderFromRequest(c.Request())

        ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

        next.ServeHTTP(w, r.WithContext(ctx))
    }
    return http.HandlerFunc(fn)
}
```

**NOTE** : The `MiddleWare` function is deprecated but will remain as a reference for users who prefer to use it.
It would be best to create your own middleware function that fits your application's requirements.

### echo middleware example: 

```go
func MiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		hxh := htmx.HxRequestHeaderFromRequest(c.Request())

		ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
```

--- 

## Custom logger 

In case you want to use a custom logger, like zap, you can inject them into the slog package like so:

```go
import (
    "go.uber.org/zap"
    "go.uber.org/zap/exp/zapslog"
)

func main() {
    // create a new htmx instance with the logger
    app := &App{
        htmx: htmx.New(),
    }

    zapLogger := zap.Must(zap.NewProduction())
    defer zapLogger.Sync()
    
    logger := slog.New(zapslog.NewHandler(zapLogger.Core(), nil))
    
    app.htmx.SetLog(logger)
}
```

--- 

## Usage in other frameworks
The htmx package is designed to be versatile and can be used in various Go web frameworks. 
Below are examples of how to use the package in two popular Go web frameworks: Echo and Gin.

### echo

```go
func (c *controller) Hello(c echo.Context) error {
    // initiate a new htmx handler 
    h := c.app.htmx.NewHandler(c.Response(), c.Request())
    
    // Example usage of Swap 
    swap := htmx.NewSwap().Swap(time.Second * 2).ScrollBottom() 
    
    h.ReSwapWithObject(swap)
    
    _, _ = h.Write([]byte("your content"))
}
```

### gin

```go
func (c *controller) Hello(c *gin.Context) {
    // initiate a new htmx handler 
    h := c.app.htmx.NewHandler(c.Writer, c.Request)
    
    // Example usage of Swap 
    swap := htmx.NewSwap().Swap(time.Second * 2).ScrollBottom() 
    
    h.ReSwapWithObject(swap)
    
    _, _ = h.Write([]byte("your content"))
}
```

--- 

## Server Sent Events (SSE)

The htmx package provides support for Server-Sent Events (SSE) in Go applications. This feature allows you to send real-time updates from the server to the client, enabling live updates and notifications in your web application.

You can read about this feature in the [htmx documentation](https://htmx.org/extensions/server-sent-events/) and the [MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events).
For htmx 4, see [Server Sent Events](#server-sent-events) in the htmx 4 section.

### Usage

Create an endpoint in your Go application to handle SSE requests. (see the example for a better understanding)
```go
func (a *App) SSE(w http.ResponseWriter, r *http.Request) {
    cl := &client{
        id: uuid.New().String()
        ch: make(chan *htmx.SSEMessage),
    }
    
    sseManager.Handle(w, cl)
}
```

To send a message to the client, you can use the `Send` method on the `SSEManager` object.

```go
go func() {
    for {
        // Send a message every seconds 
        time.Sleep(1 * time.Second) 
			
        msg := sse.
            NewMessage(fmt.Sprintf("The current time is: %v", time.Now().Format(time.RFC850))).
            WithEvent("Time")
            
        sseManager.Send()
    }
}()
``` 

A message can hold data of several lines, such as rendered HTML: each line is sent in a `data:` field of its own, and the browser joins the lines again.

### HTMX helper methods 

There are two helper methods to simplify the usage of SSE in your HTMX application.
The Manager is created in the background and is not exposed to the user.
You can change the default worker pool size by setting the `htmx.DefaultSSEWorkerPoolSize` variable.

```go

// SSEHandler handles the server-sent events. this is a shortcut and is not the preferred way to handle sse.
func (h *HTMX) SSEHandler(w http.ResponseWriter, cl sse.Client)

// SSESend sends a message to all connected clients.
func (h *HTMX) SSESend(message sse.Envelope)

```
--- 

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

**Remember to give the project a star! Thanks again!**

1. Fork this repo
2. Create a new branch with `main` as the base branch
3. Add your changes
4. Raise a Pull Request

--- 

## License

Distributed under the MIT License. See `LICENSE` for more information.