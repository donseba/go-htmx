package htmx

type HxHeaderResponse string

var (
	HXLocation           HxHeaderResponse = "HX-Location"             // Allows you to do a client-side redirect that does not do a full page reload
	HXPushUrl            HxHeaderResponse = "HX-Push-Url"             // pushes a new url into the history stack
	HXRedirect           HxHeaderResponse = "HX-Redirect"             // can be used to do a client-side redirect to a new location
	HXRefresh            HxHeaderResponse = "HX-Refresh"              // if set to "true" the client side will do a full refresh of the page
	HXReplaceUrl         HxHeaderResponse = "HX-Replace-Url"          // replaces the current URL in the location bar
	HXReswap             HxHeaderResponse = "HX-Reswap"               // Allows you to specify how the response will be swapped. See hx-swap for possible values
	HXRetarget           HxHeaderResponse = "HX-Retarget"             // A CSS selector that updates the target of the content update to a different element on the page
	HXTrigger            HxHeaderResponse = "HX-Trigger"              // allows you to trigger client side events, see the documentation for more info
	HXTriggerAfterSettle HxHeaderResponse = "HX-Trigger-After-Settle" // allows you to trigger client side events, see the documentation for more info
	HXTriggerAfterSwap   HxHeaderResponse = "HX-Trigger-After-Swap"   // allows you to trigger client side events, see the documentation for more info
)

func (h HxHeaderResponse) Values() []HxHeaderResponse {
	return []HxHeaderResponse{
		HXLocation,
		HXPushUrl,
		HXRedirect,
		HXRefresh,
		HXReplaceUrl,
		HXReswap,
		HXRetarget,
		HXTrigger,
		HXTriggerAfterSettle,
		HXTriggerAfterSwap,
	}
}

func (h HxHeaderResponse) String() string {
	return string(h)
}
