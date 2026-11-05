package htmx

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newVersionHandler(version Version, headers map[string]string) (*Handler, *httptest.ResponseRecorder) {
	h := New()
	h.SetVersion(version)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	for k, v := range headers {
		r.Header.Set(k, v)
	}

	w := httptest.NewRecorder()

	return h.NewHandler(w, r), w
}

func TestVersion(t *testing.T) {
	h := New()
	if h.Version() != Version2 {
		t.Errorf("expected the default version to be Version2, got %d", h.Version())
	}

	h.SetVersion(Version4)
	if h.Version() != Version4 {
		t.Errorf("expected the version to be Version4, got %d", h.Version())
	}
}

func TestHxRequestHeaderFromRequestHTMX4(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("HX-Request", "true")
	r.Header.Set("HX-Request-Type", "partial")
	r.Header.Set("HX-Source", "button#save")
	r.Header.Set("HX-Target", "div#content")

	header := HxRequestHeaderFromRequest(r)

	equal(t, "partial", header.HxRequestType)
	equal(t, "button#save", header.HxSource)
	equal(t, "div#content", header.HxTarget)
}

func TestRenderPartialHTMX4(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		htmx2   bool
		htmx4   bool
	}{
		{"no htmx request", map[string]string{}, false, false},
		{"request without type", map[string]string{"HX-Request": "true"}, true, true},
		{"partial request", map[string]string{"HX-Request": "true", "HX-Request-Type": "partial"}, true, true},
		{"full request", map[string]string{"HX-Request": "true", "HX-Request-Type": "full"}, true, false},
		{"boosted full request", map[string]string{"HX-Boosted": "true", "HX-Request-Type": "full"}, true, false},
		{"history restore", map[string]string{"HX-Request": "true", "HX-History-Restore-Request": "true", "HX-Request-Type": "full"}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h2, _ := newVersionHandler(Version2, tt.headers)
			equalBool(t, tt.htmx2, h2.RenderPartial())

			h4, _ := newVersionHandler(Version4, tt.headers)
			equalBool(t, tt.htmx4, h4.RenderPartial())
		})
	}
}

func TestTargetIDAndSourceID(t *testing.T) {
	h2, _ := newVersionHandler(Version2, map[string]string{
		"HX-Target":  "content",
		"HX-Trigger": "save",
		"HX-Source":  "button#other",
	})

	equal(t, "content", h2.TargetID())
	equal(t, "save", h2.SourceID())

	tests := []struct {
		identifier string
		id         string
	}{
		{"div#content", "content"},
		{"div#my%20list", "my list"},
		{"li#caf%C3%A9", "café"},
		{"div#a#b", "a#b"},
		{"div#row%zz", "row%zz"},
		{"body", ""},
		{"", ""},
	}

	for _, tt := range tests {
		h4, _ := newVersionHandler(Version4, map[string]string{
			"HX-Target":  tt.identifier,
			"HX-Source":  tt.identifier,
			"HX-Trigger": "ignored",
		})

		equal(t, tt.id, h4.TargetID())
		equal(t, tt.id, h4.SourceID())
	}
}

func TestPrompt(t *testing.T) {
	h2, _ := newVersionHandler(Version2, map[string]string{"HX-Prompt": "50%25 off"})
	equal(t, "50%25 off", h2.Prompt())

	h4, _ := newVersionHandler(Version4, map[string]string{"HX-Prompt": "caf%C3%A9%20au%20lait"})
	equal(t, "café au lait", h4.Prompt())

	h4, _ = newVersionHandler(Version4, map[string]string{"HX-Prompt": "100%"})
	equal(t, "100%", h4.Prompt())
}

func TestTriggerAfterHTMX4(t *testing.T) {
	h4, w4 := newVersionHandler(Version4, nil)

	h4.Trigger(trigger)
	h4.TriggerAfterSettle(triggerAfterSettle)
	h4.TriggerAfterSwap(triggerAfterSwap)
	h4.TriggerAfterSettleWithObject(NewTrigger().AddEvent(triggerAfterSettle))
	h4.TriggerAfterSwapWithObject(NewTrigger().AddEvent(triggerAfterSwap))

	equal(t, trigger, w4.Header().Get(HXTrigger.String()))

	for _, key := range []HxResponseKey{HXTriggerAfterSettle, HXTriggerAfterSwap} {
		if _, ok := w4.Header()[key.String()]; ok {
			t.Errorf("expected htmx 4 to leave out %s, got %q", key, w4.Header().Get(key.String()))
		}
	}

	h2, w2 := newVersionHandler(Version2, nil)

	h2.TriggerAfterSettle(triggerAfterSettle)
	h2.TriggerAfterSwap(triggerAfterSwap)

	equal(t, triggerAfterSettle, w2.Header().Get(HXTriggerAfterSettle.String()))
	equal(t, triggerAfterSwap, w2.Header().Get(HXTriggerAfterSwap.String()))
}

func TestReSwapWithObjectHTMX4(t *testing.T) {
	swap := NewSwap().Style(SwapOuterHTML).ShowTop("#main .result").FocusScroll(true).Settle(1500 * time.Millisecond)

	h4, w4 := newVersionHandler(Version4, nil)
	h4.ReSwapWithObject(swap)
	equal(t, `outerHTML show:top showTarget:"#main .result" focusScroll:true settle:1500ms`, w4.Header().Get(HXReswap.String()))

	h2, w2 := newVersionHandler(Version2, nil)
	h2.ReSwapWithObject(swap)
	equal(t, `outerHTML show:#main .result:top focus-scroll:true settle:1.5s`, w2.Header().Get(HXReswap.String()))

	// raw swap strings are sent as given
	h4.ReSwap("innerMorph")
	equal(t, "innerMorph", w4.Header().Get(HXReswap.String()))
}

func TestStopPollingHTMX4(t *testing.T) {
	h4, w4 := newVersionHandler(Version4, nil)
	h4.StopPolling()

	equalInt(t, StatusStopPolling, w4.Code)
}

func TestSwapStringForInTemplate(t *testing.T) {
	tmpl := template.Must(template.New("swap").Parse(`<div hx-swap="{{ .Swap.StringFor 4 }}"></div>`))

	var out strings.Builder
	err := tmpl.Execute(&out, map[string]any{"Swap": NewSwap().ScrollBottom("#list")})
	if err != nil {
		t.Fatal(err)
	}

	equal(t, `<div hx-swap="innerHTML scroll:bottom scrollTarget:&#34;#list&#34;"></div>`, out.String())
}
