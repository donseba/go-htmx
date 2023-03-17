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
		handler.WriteHeader(http.StatusOK)

		i, err := handler.Write(nil)
		if err != nil {
			t.Error(err)
		}

		t.Logf("data written: %d", i)
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
	assert.Equal(t, resp.Header.Get(htmx.HXLocation.String()), string(j))
	assert.Equal(t, resp.Header.Get(htmx.HXPushUrl.String()), pushURL)
	assert.Equal(t, resp.Header.Get(htmx.HXRedirect.String()), redirect)
	assert.Equal(t, resp.Header.Get(htmx.HXRefresh.String()), htmx.HxBoolToStr(refresh))
	assert.Equal(t, resp.Header.Get(htmx.HXReplaceUrl.String()), replaceURL)
	assert.Equal(t, resp.Header.Get(htmx.HXReswap.String()), reSwap)
	assert.Equal(t, resp.Header.Get(htmx.HXRetarget.String()), reTarget)
	assert.Equal(t, resp.Header.Get(htmx.HXTrigger.String()), trigger)
	assert.Equal(t, resp.Header.Get(htmx.HXTriggerAfterSwap.String()), triggerAfterSwap)
	assert.Equal(t, resp.Header.Get(htmx.HXTriggerAfterSettle.String()), triggerAfterSettle)
}
