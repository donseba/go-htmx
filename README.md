# go-htmx
**Seamless HTMX integration in golang applications.**

## Description

This repository contains the htmx Go package, designed to enhance server-side handling of HTML generated with the [HTMX library](https://htmx.org/). 
It provides a set of tools to easily manage swap behaviors, trigger configurations, and other HTMX-related functionalities in a Go server environment.

## Features

- **Swap Configuration**: Configure swap behaviors for HTMX responses, including style, timing, and scrolling.
- **Trigger Management**: Define and manage triggers for HTMX events, supporting both simple and detailed triggers.
- **Middleware Support**: Integrate HTMX seamlessly with Go middleware for easy HTMX header configuration.
- **io.Writer Support**: The HTMX handler implements the io.Writer interface for easy integration with existing Go code.

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
	"github.com/donseba/go-htmx/middleware"
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
	mux.Handle("/", middleware.MiddleWare(http.HandlerFunc(app.Home)))

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
	
	h.TriggerWithObject(swap)
	// or 
	h.TriggerAfterSettleWithObject(swap)
	// or
	h.TriggerAfterSwapWithObject(swap)
	
	_, _ = h.Write([]byte("your content"))
}
```

## Middleware
The htmx package is designed to work with middleware to catch the htmx headers and set them in the request context.
This allows you to use the headers in your handlers.A middleware that works with the standard mux router is included in the package.

### echo middleware example: 
```go
func MiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		hxh := htmx.HxHeaderRequest{
			HxBoosted:               htmx.HxStrToBool(c.Request().Header.Get("HX-Boosted")),
			HxCurrentURL:            c.Request().Header.Get("HX-Current-URL"),
			HxHistoryRestoreRequest: htmx.HxStrToBool(c.Request().Header.Get("HX-History-Restore-Request")),
			HxPrompt:                c.Request().Header.Get("HX-Prompt"),
			HxRequest:               htmx.HxStrToBool(c.Request().Header.Get("HX-Request")),
			HxTarget:                c.Request().Header.Get("HX-Target"),
			HxTriggerName:           c.Request().Header.Get("HX-Trigger-Name"),
			HxTrigger:               c.Request().Header.Get("HX-Trigger"),
		}

		ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
```

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

**Remember to give the project a star! Thanks again!**

1. Fork this repo
2. Create a new branch with `main` as the base branch
3. Add your changes
4. Raise a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.