package htmx

// Version is the major version of htmx the browser runs. It decides how the
// request headers are read and how the response headers are written.
type Version int

const (
	// Version2 is htmx 2, the default.
	Version2 Version = 2

	// Version4 is htmx 4. See https://four.htmx.org/docs/whats-new-in-htmx-4
	Version4 Version = 4
)
