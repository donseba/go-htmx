package htmx

import (
	"net/http"
	"testing"
)

func TestNewTriggerMixed(t *testing.T) {
	trigger := NewTrigger()

	if trigger == nil {
		t.Error("expected trigger to not be nil")
	}

	trigger.AddEvent("foo").
		AddEventDetailed("bar", "baz").
		AddEventDetailed("qux", "quux").
		AddEventObject("corge", map[string]any{"grault": "garply", "waldo": "fred", "plugh": "xyzzy", "thud": "wibble"})

	expected := `{"bar":"baz","corge":{"grault":"garply","plugh":"xyzzy","thud":"wibble","waldo":"fred"},"foo":"","qux":"quux"}`

	if trigger.String() != expected {
		t.Errorf("expected trigger to be %v, got %v", expected, trigger.String())
	}
}

func TestNewTriggerSingle(t *testing.T) {
	trigger := NewTrigger()

	if trigger == nil {
		t.Error("expected trigger to not be nil")
	}

	trigger.AddEvent("foo").
		AddEvent("bar").
		AddEvent("baz")

	expected := "foo, bar, baz"

	if trigger.String() != expected {
		t.Errorf("expected trigger to be %v, got %v", expected, trigger.String())
	}
}

func TestNewTriggerMixedNested(t *testing.T) {
	trigger := NewTrigger()

	if trigger == nil {
		t.Error("expected trigger to not be nil")
	}

	trigger.AddEvent("foo").
		AddEventDetailed("bar", "baz").
		AddEventDetailed("qux", "quux").
		AddEventObject("corge", map[string]any{"grault": "garply", "waldo": "fred", "plugh": "xyzzy", "thud": map[string]any{"foo": "bar", "baz": "qux"}}).AddSuccess("successfully tested", map[string]any{"foo": "bar", "baz": "qux"})
	expected := `{"bar":"baz","corge":{"grault":"garply","plugh":"xyzzy","thud":{"baz":"qux","foo":"bar"},"waldo":"fred"},"foo":"","qux":"quux","showMessage":{"baz":"qux","foo":"bar","level":"success","message":"successfully tested"}}`

	if trigger.String() != expected {
		t.Errorf("expected trigger to be %v, got %v", expected, trigger.String())
	}
}

func TestTriggerSuccess(t *testing.T) {
	req := &http.Request{}
	handler := New().NewHandler(dummyWriter{}, req)
	handler.TriggerSuccess("successfully tested", map[string]any{"foo": "bar", "baz": "qux"})

	expected := `{"showMessage":{"baz":"qux","foo":"bar","level":"success","message":"successfully tested"}}`

	equal(t, expected, handler.response.Get(HXTrigger))
}

func TestTriggerError(t *testing.T) {
	req := &http.Request{}
	handler := New().NewHandler(dummyWriter{}, req)
	handler.TriggerError("successfully tested a fail", map[string]any{"foo": "bar", "baz": "qux"})

	expected := `{"showMessage":{"baz":"qux","foo":"bar","level":"error","message":"successfully tested a fail"}}`

	equal(t, expected, handler.response.Get(HXTrigger))
}

func TestTriggerInfo(t *testing.T) {
	req := &http.Request{}
	handler := New().NewHandler(dummyWriter{}, req)
	handler.TriggerInfo("successfully tested some info", map[string]any{"foo": "bar", "baz": "qux"})

	expected := `{"showMessage":{"baz":"qux","foo":"bar","level":"info","message":"successfully tested some info"}}`

	equal(t, expected, handler.response.Get(HXTrigger))
}

func TestTriggerWarning(t *testing.T) {
	req := &http.Request{}
	handler := New().NewHandler(dummyWriter{}, req)
	handler.TriggerWarning("successfully tested a warning", map[string]any{"foo": "bar", "baz": "qux"})

	expected := `{"showMessage":{"baz":"qux","foo":"bar","level":"warning","message":"successfully tested a warning"}}`

	equal(t, expected, handler.response.Get(HXTrigger))
}
