package htmx

import "testing"

func TestNewTriggerMixed(t *testing.T) {
	trigger := NewTrigger()

	if trigger == nil {
		t.Error("expected trigger to not be nil")
	}

	trigger.AddEvent("foo").
		AddEventDetailed("bar", "baz").
		AddEventDetailed("qux", "quux").
		AddEventObject("corge", map[string]string{"grault": "garply", "waldo": "fred", "plugh": "xyzzy", "thud": "wibble"})

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
