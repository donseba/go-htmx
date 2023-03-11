# go-htmx
unobstructive HTMX integration in golang applications.

This package consists of two main parts, 
1) Middleware to catch the HTMX headers from the request
2) Handler with io.Writer interface to serve content.

# Getting started
`go get github.com/donseba/go-htmx`

initialise the htmx service like so : 
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

	// set the headers for the response, see docs for more options
	h.PushURL("http://push.url")
	h.ReTarget("#ReTarged")

	// write the output like you normally do.
	// check inspector tool in browser to see that the headers are set.
	_, _ = h.Write([]byte("OK"))
}
```


echo middleware example: 
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