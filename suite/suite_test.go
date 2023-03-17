package suite

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/middleware"
	"github.com/stretchr/testify/assert"
)

var (
	location = &htmx.LocationInput{
		Source:  "source",
		Event:   "",
		Handler: "",
		Target:  "http://new-url.com",
		Swap:    "",
		Values:  nil,
		Header:  nil,
	}
	pushURL            = "http://push-url.com"
	redirect           = "http://redirect.com"
	refresh            = true
	replaceURL         = "http://replace-url.com"
	reSwap             = "#reSwap"
	reTarget           = "#reTarget"
	trigger            = "#trigger"
	triggerAfterSettle = "#triggerAfterSettle"
	triggerAfterSwap   = "#triggerAfterSwap"
)

func TestNew(t *testing.T) {
	h := htmx.New()

	svr := httptest.NewServer(middleware.MiddleWare(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := h.NewHandler(w, r)

		_ = handler.Location(location)
		handler.PushURL(pushURL)
		handler.Redirect(redirect)
		handler.Refresh(true)
		handler.ReplaceURL(replaceURL)
		handler.ReSwap(reSwap)
		handler.ReTarget(reTarget)
		handler.Trigger(trigger)
		handler.TriggerAfterSettle(triggerAfterSettle)
		handler.TriggerAfterSwap(triggerAfterSwap)
		handler.WriteHeader(http.StatusAccepted)

		_, err := handler.Write([]byte("hi"))
		if err != nil {
			t.Error(err)
		}
	})))
	defer svr.Close()

	resp, err := http.Get(svr.URL)
	if err != nil {
		t.Error("an error occurred while making the request")
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error("an error occurred when reading the response")
	}

	j, _ := json.Marshal(location)
	assert.Equal(t, string(j), resp.Header.Get(htmx.HXLocation.String()))
	assert.Equal(t, pushURL, resp.Header.Get(htmx.HXPushUrl.String()))
	assert.Equal(t, redirect, resp.Header.Get(htmx.HXRedirect.String()))
	assert.Equal(t, htmx.HxBoolToStr(refresh), resp.Header.Get(htmx.HXRefresh.String()))
	assert.Equal(t, replaceURL, resp.Header.Get(htmx.HXReplaceUrl.String()))
	assert.Equal(t, reSwap, resp.Header.Get(htmx.HXReswap.String()))
	assert.Equal(t, reTarget, resp.Header.Get(htmx.HXRetarget.String()))
	assert.Equal(t, trigger, resp.Header.Get(htmx.HXTrigger.String()))
	assert.Equal(t, triggerAfterSwap, resp.Header.Get(htmx.HXTriggerAfterSwap.String()))
	assert.Equal(t, triggerAfterSettle, resp.Header.Get(htmx.HXTriggerAfterSettle.String()))
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}
