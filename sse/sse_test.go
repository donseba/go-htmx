package sse

import (
	"strings"
	"testing"
)

func TestMessageString(t *testing.T) {
	tests := []struct {
		message  Envelope
		expected string
	}{
		{NewMessage("<div>hi</div>"), "data: <div>hi</div>\n\n"},
		{NewMessage("<div>hi</div>").WithEvent("time"), "event: time\ndata: <div>hi</div>\n\n"},
		{NewMessage(""), "data: \n\n"},
		{NewMessage("<ul>\n<li>a</li>\n</ul>"), "data: <ul>\ndata: <li>a</li>\ndata: </ul>\n\n"},
		{NewMessage("a\r\nb\rc"), "data: a\ndata: b\ndata: c\n\n"},
		{NewMessage("a\n"), "data: a\ndata: \n\n"},
		{NewMessage("a\nevent: other").WithEvent("time"), "event: time\ndata: a\ndata: event: other\n\n"},
	}

	for _, tt := range tests {
		if got := tt.message.String(); got != tt.expected {
			t.Errorf("expected message %q, got %q", tt.expected, got)
		}
	}
}

// TestMessageDataArrivesAsSent reads the data of a message the way a browser does,
// which joins the data fields with newlines.
func TestMessageDataArrivesAsSent(t *testing.T) {
	for _, data := range []string{"", "one line", "<table>\n  <tr><td>S</td></tr>\n</table>", "trailing\n", "\n\nleading"} {
		if got := receivedData(NewMessage(data).String()); got != data {
			t.Errorf("expected data %q, got %q", data, got)
		}
	}
}

func receivedData(message string) string {
	var data []string

	for _, line := range strings.Split(strings.TrimSuffix(message, "\n\n"), "\n") {
		if value, ok := strings.CutPrefix(line, "data: "); ok {
			data = append(data, value)
		}
	}

	return strings.Join(data, "\n")
}
