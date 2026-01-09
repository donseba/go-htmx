package htmx

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	location = &LocationInput{
		Path:    "http://new-url.com",
		Source:  "source",
		Event:   "",
		Handler: "",
		Target:  "body",
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
	reSelect           = "#reSelect"
	trigger            = "#trigger"
	triggerAfterSettle = "#triggerAfterSettle"
	triggerAfterSwap   = "#triggerAfterSwap"

	reswapWithObject = NewSwap().ScrollTop().Settle(1 * time.Second)
)

func TestNew(t *testing.T) {
	h := New()
	if h == nil {
		t.Errorf("expected htmx to be initialized")
	}

	w := &dummyWriter{
		Writer: &httptest.ResponseRecorder{},
	}
	r := &http.Request{
		Header: http.Header{},
	}
	r.Header.Set("HX-Request", "true")
	r.Header.Set("HX-Boosted", "true")
	r.Header.Set("HX-History-Restore-Request", "false")

	handler := h.NewHandler(w, r)
	if handler == nil {
		t.Errorf("expected handler to be initialized")
	}

	equalBool(t, true, handler.IsHxRequest())
	equalBool(t, true, handler.IsHxBoosted())
	equalBool(t, false, handler.IsHxHistoryRestoreRequest())
	equalBool(t, true, handler.RenderPartial())

	_ = handler.Location(location)
	handler.PushURL(pushURL)
	handler.Redirect(redirect)
	handler.Refresh(refresh)
	handler.ReplaceURL(replaceURL)
	handler.ReSwap(reSwap)
	handler.ReTarget(reTarget)
	handler.ReSelect(reSelect)
	handler.Trigger(trigger)
	handler.TriggerAfterSettle(triggerAfterSettle)
	handler.TriggerAfterSwap(triggerAfterSwap)
	handler.WriteHeader(http.StatusAccepted)

	j, _ := json.Marshal(location)
	equal(t, string(j), handler.ResponseHeader(HXLocation))
	equal(t, pushURL, handler.ResponseHeader(HXPushUrl))
	equal(t, redirect, handler.ResponseHeader(HXRedirect))
	equal(t, HxBoolToStr(refresh), handler.ResponseHeader(HXRefresh))
	equal(t, replaceURL, handler.ResponseHeader(HXReplaceUrl))
	equal(t, reSwap, handler.ResponseHeader(HXReswap))
	equal(t, reTarget, handler.ResponseHeader(HXRetarget))
	equal(t, reSelect, handler.ResponseHeader(HXReselect))
	equal(t, trigger, handler.ResponseHeader(HXTrigger))
	equal(t, triggerAfterSwap, handler.ResponseHeader(HXTriggerAfterSwap))
	equal(t, triggerAfterSettle, handler.ResponseHeader(HXTriggerAfterSettle))

	handler.ReSwapWithObject(reswapWithObject)
	handler.TriggerWithObject(NewTrigger().AddEvent(trigger))
	handler.TriggerAfterSettleWithObject(NewTrigger().AddEvent(triggerAfterSettle))
	handler.TriggerAfterSwapWithObject(NewTrigger().AddEvent(triggerAfterSwap))

	equal(t, reswapWithObject.String(), handler.ResponseHeader(HXReswap))

	head := handler.Header()
	equal(t, "true", head.Get("Hx-Request"))

	handler.StopPolling()

	req := handler.Request()

	equalBool(t, req.HxBoosted, handler.IsHxBoosted())
	equalBool(t, req.HxHistoryRestoreRequest, handler.IsHxHistoryRestoreRequest())
	equalBool(t, req.HxRequest, handler.IsHxRequest())

	equalBool(t, req.HxBoosted, IsHxBoosted(r))
	equalBool(t, req.HxHistoryRestoreRequest, IsHxHistoryRestoreRequest(r))
	equalBool(t, req.HxRequest, IsHxRequest(r))

	i, _ := handler.Write([]byte("hi"))
	equalInt(t, 2, i)
}

func TestNoRouter(t *testing.T) {
	h := New()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
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
	equal(t, string(j), resp.Header.Get(HXLocation.String()))
	equal(t, pushURL, resp.Header.Get(HXPushUrl.String()))
	equal(t, redirect, resp.Header.Get(HXRedirect.String()))
	equal(t, HxBoolToStr(refresh), resp.Header.Get(HXRefresh.String()))
	equal(t, replaceURL, resp.Header.Get(HXReplaceUrl.String()))
	equal(t, reSwap, resp.Header.Get(HXReswap.String()))
	equal(t, reTarget, resp.Header.Get(HXRetarget.String()))
	equal(t, trigger, resp.Header.Get(HXTrigger.String()))
	equal(t, triggerAfterSwap, resp.Header.Get(HXTriggerAfterSwap.String()))
	equal(t, triggerAfterSettle, resp.Header.Get(HXTriggerAfterSettle.String()))
	equalInt(t, http.StatusAccepted, resp.StatusCode)
}

func TestHxResponseKey_String(t *testing.T) {
	equal(t, "HX-Location", HXLocation.String())
	equal(t, "HX-Push-Url", HXPushUrl.String())
	equal(t, "HX-Redirect", HXRedirect.String())
	equal(t, "HX-Refresh", HXRefresh.String())
	equal(t, "HX-Replace-Url", HXReplaceUrl.String())
	equal(t, "HX-Reswap", HXReswap.String())
	equal(t, "HX-Retarget", HXRetarget.String())
	equal(t, "HX-Reselect", HXReselect.String())
	equal(t, "HX-Trigger", HXTrigger.String())
	equal(t, "HX-Trigger-After-Settle", HXTriggerAfterSettle.String())
	equal(t, "HX-Trigger-After-Swap", HXTriggerAfterSwap.String())
}

func TestStopPolling(t *testing.T) {
	h := New()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := h.NewHandler(w, r)

		_ = handler.Location(location)
		handler.WriteHeader(StatusStopPolling)

		_, err := handler.Write([]byte("hi"))
		if err != nil {
			t.Error(err)
		}
	}))
	defer svr.Close()

	resp, err := http.Get(svr.URL)
	if err != nil {
		t.Error("an error occurred while making the request")
		return
	}
	defer resp.Body.Close()

	equalInt(t, StatusStopPolling, resp.StatusCode)
}

func TestSwap(t *testing.T) {
	h := New()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := h.NewHandler(w, r)

		_ = handler.Location(location)
		handler.ReSwapWithObject(NewSwap().ScrollTop().Settle(1 * time.Second))

		_, err := handler.Write([]byte("hi"))
		if err != nil {
			t.Error(err)
		}
	}))
	defer svr.Close()

	resp, err := http.Get(svr.URL)
	if err != nil {
		t.Error("an error occurred while making the request")
		return
	}
	defer resp.Body.Close()

	equal(t, "innerHTML scroll:top settle:1s", resp.Header.Get(HXReswap.String()))
}

func TestHxStrToBool(t *testing.T) {
	equalBool(t, true, HxStrToBool("true"))
	equalBool(t, false, HxStrToBool("false"))
	equalBool(t, false, HxStrToBool("not a bool"))
}

func TestHxBoolToStr(t *testing.T) {
	equal(t, "true", HxBoolToStr(true))
	equal(t, "false", HxBoolToStr(false))
}

type dummyWriter struct {
	io.Writer
}

func (dummyWriter) Header() http.Header {
	h := http.Header{}
	h.Set("HX-Request", "true")
	return h
}
func (dummyWriter) WriteHeader(int) {}

func equalBool(t *testing.T, expected, actual bool) {
	if expected != actual {
		t.Errorf("expected %t, got %t", expected, actual)
	}
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
