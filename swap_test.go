package htmx

import (
	"testing"
	"time"
)

// TestSwapper_DefaultValues tests the default values set by Swapper
func TestSwapper_DefaultValues(t *testing.T) {
	swap := NewSwap()

	if swap.style != SwapInnerHTML {
		t.Errorf("expected default style to be SwapInnerHTML, got %v", swap.style)
	}

	// Add more tests for other default values if necessary
}

// TestStyle tests the Style method
func TestStyle(t *testing.T) {
	swap := NewSwap().Style(SwapOuterHTML)

	if swap.style != SwapOuterHTML {
		t.Errorf("expected style to be SwapOuterHTML, got %v", swap.style)
	}
}

// TestString tests the String method
func TestString(t *testing.T) {
	swap := NewSwap()
	expected := "innerHTML"

	if swap.String() != expected {
		t.Errorf("expected string output to be %v, got %v", expected, swap.String())
	}

	// Additional scenarios for String method can be added here
}

// TestTimingSwap tests the TimingSwap method
func TestTimingSwap(t *testing.T) {
	duration := 100 * time.Millisecond
	swap := NewSwap().Swap(duration)

	if swap.timing == nil || swap.timing.duration != duration {
		t.Errorf("expected timing swap to be %v, got %v", duration, swap.timing.duration)
	}
}

// TestTimingSettle tests the TimingSettle method
func TestTimingSettle(t *testing.T) {
	duration := 200 * time.Millisecond
	swap := NewSwap().Settle(duration)

	if swap.timing == nil || swap.timing.duration != duration {
		t.Errorf("expected timing settle to be %v, got %v", duration, swap.timing.duration)
	}
}

// TestTransition tests the Transition method
func TestTransition(t *testing.T) {
	swap := NewSwap().Transition(true)

	if swap.transition == nil || *swap.transition != true {
		t.Errorf("expected transition to be true, got %v", swap.transition)
	}
	expected := `innerHTML transition:true`
	if swap.String() != expected {
		t.Errorf("expected string output to be %s, got %s", expected, swap.String())
	}

}

// TestIgnoreTitle tests the IgnoreTitle method
func TestIgnoreTitle(t *testing.T) {
	swap := NewSwap().IgnoreTitle(true)

	if swap.ignoreTitle == nil || *swap.ignoreTitle != true {
		t.Errorf("expected ignoreTitle to be true, got %v", swap.ignoreTitle)
	}

	expected := `innerHTML ignoreTitle:true`
	if swap.String() != "innerHTML ignoreTitle:true" {
		t.Errorf("expected string output to be %s, got %s", expected, swap.String())
	}
}

// TestFocusScroll tests the FocusScroll method
func TestFocusScroll(t *testing.T) {
	swap := NewSwap().FocusScroll(true)

	if swap.focusScroll == nil || *swap.focusScroll != true {
		t.Errorf("expected focusScroll to be true, got %v", swap.focusScroll)
	}
}

// TestScrollingScroll tests the ScrollingScroll method
func TestScrollingScroll(t *testing.T) {
	swap := NewSwap().Scroll(SwapDirectionTop)

	if swap.scrolling == nil || swap.scrolling.direction != SwapDirectionTop || swap.scrolling.mode != ScrollingScroll {
		t.Errorf("expected scrolling mode to be ScrollingScroll and direction to be SwapDirectionTop, got mode: %v, direction: %v", swap.scrolling.mode, swap.scrolling.direction)
	}
}

// TestScrollingShow tests the ScrollingShow method
func TestScrollingShow(t *testing.T) {
	swap := NewSwap().Show(SwapDirectionBottom)

	if swap.scrolling == nil || swap.scrolling.direction != SwapDirectionBottom || swap.scrolling.mode != ScrollingShow {
		t.Errorf("expected scrolling mode to be ScrollingShow and direction to be SwapDirectionBottom, got mode: %v, direction: %v", swap.scrolling.mode, swap.scrolling.direction)
	}
}

// TestScrollingScrollTop tests the ScrollingScrollTop method
func TestScrollingScrollTop(t *testing.T) {
	target := "#element"
	swap := NewSwap().ScrollTop(target)

	if swap.scrolling == nil || swap.scrolling.direction != SwapDirectionTop || swap.scrolling.mode != ScrollingScroll || swap.scrolling.target != target {
		t.Errorf("expected scrolling mode to be ScrollingScroll, direction to be SwapDirectionTop, and target to be %v, got mode: %v, direction: %v, target: %v", target, swap.scrolling.mode, swap.scrolling.direction, swap.scrolling.target)
	}
}

// TestScrollingScrollBottom tests the ScrollingScrollBottom method
func TestScrollingScrollBottom(t *testing.T) {
	target := "#element"
	swap := NewSwap().ScrollBottom(target)

	if swap.scrolling == nil || swap.scrolling.direction != SwapDirectionBottom || swap.scrolling.mode != ScrollingScroll || swap.scrolling.target != target {
		t.Errorf("expected scrolling mode to be ScrollingScroll, direction to be SwapDirectionBottom, and target to be %v, got mode: %v, direction: %v, target: %v", target, swap.scrolling.mode, swap.scrolling.direction, swap.scrolling.target)
	}
}

// TestScrollingShowTop tests the ScrollingShowTop method
func TestScrollingShowTop(t *testing.T) {
	target := "#element"
	swap := NewSwap().ShowTop(target)

	if swap.scrolling == nil || swap.scrolling.direction != SwapDirectionTop || swap.scrolling.mode != ScrollingShow || swap.scrolling.target != target {
		t.Errorf("expected scrolling mode to be ScrollingShow, direction to be SwapDirectionTop, and target to be %v, got mode: %v, direction: %v, target: %v", target, swap.scrolling.mode, swap.scrolling.direction, swap.scrolling.target)
	}
}

// TestScrollingShowBottom tests the ScrollingShowBottom method
func TestScrollingShowBottom(t *testing.T) {
	target := "#element"
	swap := NewSwap().ShowBottom(target)

	if swap.scrolling == nil || swap.scrolling.direction != SwapDirectionBottom || swap.scrolling.mode != ScrollingShow || swap.scrolling.target != target {
		t.Errorf("expected scrolling mode to be ScrollingShow, direction to be SwapDirectionBottom, and target to be %v, got mode: %v, direction: %v, target: %v", target, swap.scrolling.mode, swap.scrolling.direction, swap.scrolling.target)
	}
}

// TestStringForVersion4 tests the htmx 4 syntax of the String method
func TestStringForVersion4(t *testing.T) {
	tests := []struct {
		swap     *Swap
		expected string
	}{
		{NewSwap(), "innerHTML"},
		{NewSwap().ScrollBottom(), "innerHTML scroll:bottom"},
		{NewSwap().ScrollBottom("#list"), `innerHTML scroll:bottom scrollTarget:"#list"`},
		{NewSwap().ShowTop("window"), `innerHTML show:top showTarget:"window"`},
		{NewSwap().ShowTop(`[name="q"]`), `innerHTML show:top showTarget:'[name="q"]'`},
		{NewSwap().Swap(), "innerHTML swap:0s"},
		{NewSwap().Settle(), "innerHTML settle:20ms"},
		{NewSwap().Swap(90 * time.Second), "innerHTML swap:90s"},
		{NewSwap().Settle(1500 * time.Millisecond), "innerHTML settle:1500ms"},
		{NewSwap().Swap(500 * time.Microsecond), "innerHTML swap:0.5ms"},
		{NewSwap().Transition(true).IgnoreTitle(false), "innerHTML transition:true ignoreTitle:false"},
		{NewSwap().FocusScroll(false), "innerHTML focusScroll:false"},
		{NewSwap().Style(SwapInnerMorph), "innerMorph"},
		{NewSwap().Style(SwapOuterMorph), "outerMorph"},
		{NewSwap().Style(SwapOuterSync), "outerSync"},
		{NewSwap().Style(SwapTextContent), "textContent"},
	}

	for _, tt := range tests {
		if got := tt.swap.StringFor(Version4); got != tt.expected {
			t.Errorf("expected string output to be %s, got %s", tt.expected, got)
		}
	}
}

// TestStringForVersion2 tests that the htmx 2 syntax of the String method stays as it is
func TestStringForVersion2(t *testing.T) {
	swap := NewSwap().Style(SwapOuterHTML).ShowTop("#main").Transition(true).IgnoreTitle(false).FocusScroll(true).Settle(1500 * time.Millisecond)
	expected := "outerHTML show:#main:top transition:true ignoreTitle:false focus-scroll:true settle:1.5s"

	if swap.String() != expected {
		t.Errorf("expected string output to be %s, got %s", expected, swap.String())
	}

	if swap.StringFor(Version2) != expected {
		t.Errorf("expected Version2 output to be %s, got %s", expected, swap.StringFor(Version2))
	}

	if swap.StringFor(Version(3)) != expected {
		t.Errorf("expected an unknown version to write htmx 2 syntax %s, got %s", expected, swap.StringFor(Version(3)))
	}

	bare := NewSwap().Swap()
	if bare.String() != "innerHTML swap" {
		t.Errorf("expected string output to be innerHTML swap, got %s", bare.String())
	}
}
