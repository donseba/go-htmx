package htmx

type (
	HxResponseKey string

	HxResponseHeader struct {
		Headers map[HxResponseKey]string
	}
)

var (
	HXLocation           HxResponseKey = "HX-Location"             // Allows you to do a client-side redirect that does not do a full page reload
	HXPushUrl            HxResponseKey = "HX-Push-Url"             // pushes a new url into the history stack
	HXRedirect           HxResponseKey = "HX-Redirect"             // can be used to do a client-side redirect to a new location
	HXRefresh            HxResponseKey = "HX-Refresh"              // if set to "true" the client side will do a full refresh of the page
	HXReplaceUrl         HxResponseKey = "HX-Replace-Url"          // replaces the current URL in the location bar
	HXReswap             HxResponseKey = "HX-Reswap"               // Allows you to specify how the response will be swapped. See hx-swap for possible values
	HXRetarget           HxResponseKey = "HX-Retarget"             // A CSS selector that updates the target of the content update to a different element on the page
	HXTrigger            HxResponseKey = "HX-Trigger"              // allows you to trigger client side events, see the documentation for more info
	HXTriggerAfterSettle HxResponseKey = "HX-Trigger-After-Settle" // allows you to trigger client side events, see the documentation for more info
	HXTriggerAfterSwap   HxResponseKey = "HX-Trigger-After-Swap"   // allows you to trigger client side events, see the documentation for more info
)

func (h HxResponseKey) String() string {
	return string(h)
}

func (h *HxResponseHeader) Set(k HxResponseKey, val string) {
	h.Headers[k] = val
}
