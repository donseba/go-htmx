package sse

import (
	"testing"
)

func TestMessageString(t *testing.T) {
	tests := []struct {
		message  Envelope
		expected string
	}{
		{NewMessage("<div>hi</div>"), "data: <div>hi</div>\n\n"},
		{NewMessage("<div>hi</div>").WithEvent("time"), "event: time\ndata: <div>hi</div>\n\n"},
	}

	for _, tt := range tests {
		if got := tt.message.String(); got != tt.expected {
			t.Errorf("expected message %q, got %q", tt.expected, got)
		}
	}
}
