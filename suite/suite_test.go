package suite

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestHxStrToBool(t *testing.T) {
	assert.True(t, htmx.HxStrToBool("true"))
	assert.False(t, htmx.HxStrToBool("false"))
	assert.False(t, htmx.HxStrToBool("not a bool"))
}

func TestHxBoolToStr(t *testing.T) {
	assert.Equal(t, "true", htmx.HxBoolToStr(true))
	assert.Equal(t, "false", htmx.HxBoolToStr(false))
}

func TestHxResponseKey_String(t *testing.T) {
	assert.Equal(t, "HX-Location", htmx.HXLocation.String())
	assert.Equal(t, "HX-Push-Url", htmx.HXPushUrl.String())
	assert.Equal(t, "HX-Redirect", htmx.HXRedirect.String())
	assert.Equal(t, "HX-Refresh", htmx.HXRefresh.String())
	assert.Equal(t, "HX-Replace-Url", htmx.HXReplaceUrl.String())
	assert.Equal(t, "HX-Reswap", htmx.HXReswap.String())
	assert.Equal(t, "HX-Retarget", htmx.HXRetarget.String())
	assert.Equal(t, "HX-Reselect", htmx.HXReselect.String())
	assert.Equal(t, "HX-Trigger", htmx.HXTrigger.String())
	assert.Equal(t, "HX-Trigger-After-Settle", htmx.HXTriggerAfterSettle.String())
	assert.Equal(t, "HX-Trigger-After-Swap", htmx.HXTriggerAfterSwap.String())
}

func TestStopPolling(t *testing.T) {
	h := htmx.New()

	svr := httptest.NewServer(middleware.MiddleWare(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := h.NewHandler(w, r)

		_ = handler.Location(location)
		handler.WriteHeader(htmx.StatusStopPolling)

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

	assert.Equal(t, htmx.StatusStopPolling, resp.StatusCode)
}

func TestSwap(t *testing.T) {
	h := htmx.New()

	svr := httptest.NewServer(middleware.MiddleWare(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := h.NewHandler(w, r)

		_ = handler.Location(location)
		handler.ReSwapWithObject(htmx.NewSwap().ScrollTop().Settle(1 * time.Second))

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

	t.Log(resp.Header.Get(htmx.HXReswap.String()))

	assert.Equal(t, "innerHTML scroll:top settle:1s", resp.Header.Get(htmx.HXReswap.String()))
}
