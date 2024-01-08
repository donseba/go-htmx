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
	equal(t, string(j), resp.Header.Get(htmx.HXLocation.String()))
	equal(t, pushURL, resp.Header.Get(htmx.HXPushUrl.String()))
	equal(t, redirect, resp.Header.Get(htmx.HXRedirect.String()))
	equal(t, htmx.HxBoolToStr(refresh), resp.Header.Get(htmx.HXRefresh.String()))
	equal(t, replaceURL, resp.Header.Get(htmx.HXReplaceUrl.String()))
	equal(t, reSwap, resp.Header.Get(htmx.HXReswap.String()))
	equal(t, reTarget, resp.Header.Get(htmx.HXRetarget.String()))
	equal(t, trigger, resp.Header.Get(htmx.HXTrigger.String()))
	equal(t, triggerAfterSwap, resp.Header.Get(htmx.HXTriggerAfterSwap.String()))
	equal(t, triggerAfterSettle, resp.Header.Get(htmx.HXTriggerAfterSettle.String()))
	equalInt(t, http.StatusAccepted, resp.StatusCode)
}

func TestHxStrToBool(t *testing.T) {
	equalBool(t, true, htmx.HxStrToBool("true"))
	equalBool(t, false, htmx.HxStrToBool("false"))
	equalBool(t, false, htmx.HxStrToBool("not a bool"))
}

func TestHxBoolToStr(t *testing.T) {
	equal(t, "true", htmx.HxBoolToStr(true))
	equal(t, "false", htmx.HxBoolToStr(false))
}

func TestHxResponseKey_String(t *testing.T) {
	equal(t, "HX-Location", htmx.HXLocation.String())
	equal(t, "HX-Push-Url", htmx.HXPushUrl.String())
	equal(t, "HX-Redirect", htmx.HXRedirect.String())
	equal(t, "HX-Refresh", htmx.HXRefresh.String())
	equal(t, "HX-Replace-Url", htmx.HXReplaceUrl.String())
	equal(t, "HX-Reswap", htmx.HXReswap.String())
	equal(t, "HX-Retarget", htmx.HXRetarget.String())
	equal(t, "HX-Reselect", htmx.HXReselect.String())
	equal(t, "HX-Trigger", htmx.HXTrigger.String())
	equal(t, "HX-Trigger-After-Settle", htmx.HXTriggerAfterSettle.String())
	equal(t, "HX-Trigger-After-Swap", htmx.HXTriggerAfterSwap.String())
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

	equalInt(t, htmx.StatusStopPolling, resp.StatusCode)
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

	equal(t, "innerHTML scroll:top settle:1s", resp.Header.Get(htmx.HXReswap.String()))
}

func equal(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func equalInt(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected %d, got %d", expected, actual)
	}
}

func equalBool(t *testing.T, expected, actual bool) {
	if expected != actual {
		t.Errorf("expected %t, got %t", expected, actual)
	}
}
