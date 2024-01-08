package htmx

import (
	"fmt"
	"strings"
	"time"
)

type Swap struct {
	style       SwapStyle
	transition  *bool
	timing      *SwapTiming
	scrolling   *SwapScrolling
	ignoreTitle *bool
	focusScroll *bool
}

type SwapTiming struct {
	mode     SwapTimingMode
	duration time.Duration
}

func (s *SwapTiming) String() string {
	var out string

	out = string(s.mode)

	if s.duration != 0 {
		out += ":" + s.duration.String()
	}

	return out
}

type SwapScrolling struct {
	mode      SwapScrollingMode
	target    string
	direction SwapDirection
}

func (s *SwapScrolling) String() string {
	var out string

	out = string(s.mode)

	if s.target != "" {
		out += ":" + s.target
	}

	if s.direction != "" {
		out += ":" + s.direction.String()
	}

	return out
}

// NewSwap returns a new Swap
func NewSwap() *Swap {
	return &Swap{
		style: SwapInnerHTML,
	}
}

// Style sets the style of the swap, default is innerHTML and can be changed in htmx.config.defaultSwapStyle
func (s *Swap) Style(style SwapStyle) *Swap {
	s.style = style
	return s
}

// setScrolling sets the scrolling behavior
func (s *Swap) setScrolling(mode SwapScrollingMode, direction SwapDirection, target ...string) *Swap {
	scrolling := &SwapScrolling{
		mode:      mode,
		direction: direction,
	}

	if len(target) > 0 {
		scrolling.target = target[0]
	}

	s.scrolling = scrolling
	return s
}

// Scroll sets the scrolling behavior to scroll to the top or bottom
func (s *Swap) Scroll(direction SwapDirection, target ...string) *Swap {
	return s.setScrolling(ScrollingScroll, direction, target...)
}

// ScrollTop sets the scrolling behavior to scroll to the top of the target element
func (s *Swap) ScrollTop(target ...string) *Swap {
	return s.Scroll(SwapDirectionTop, target...)
}

// ScrollBottom sets the scrolling behavior to scroll to the bottom of the target element
func (s *Swap) ScrollBottom(target ...string) *Swap {
	return s.Scroll(SwapDirectionBottom, target...)
}

// Show sets the scrolling behavior to scroll to the top or bottom
func (s *Swap) Show(direction SwapDirection, target ...string) *Swap {
	return s.setScrolling(ScrollingShow, direction, target...)
}

// ShowTop sets the scrolling behavior to scroll to the top of the target element
func (s *Swap) ShowTop(target ...string) *Swap {
	return s.Show(SwapDirectionTop, target...)
}

// ShowBottom sets the scrolling behavior to scroll to the bottom of the target element
func (s *Swap) ShowBottom(target ...string) *Swap {
	return s.Show(SwapDirectionBottom, target...)
}

// setTiming modifies the amount of time that htmx will wait after receiving a response to swap or settle the content
func (s *Swap) setTiming(mode SwapTimingMode, swap ...time.Duration) *Swap {
	var duration time.Duration
	if len(swap) > 0 {
		duration = swap[0]
	} else {
		switch mode {
		case TimingSwap:
			duration = DefaultSwapDuration
		case TimingSettle:
			duration = DefaultSettleDelay
		}
	}

	s.timing = &SwapTiming{
		mode:     mode,
		duration: duration,
	}
	return s
}

// Swap modifies the amount of time that htmx will wait after receiving a response to swap the content
func (s *Swap) Swap(swap ...time.Duration) *Swap {
	return s.setTiming(TimingSwap, swap...)
}

// Settle modifies the amount of time that htmx will wait after receiving a response to settle the content
func (s *Swap) Settle(swap ...time.Duration) *Swap {
	return s.setTiming(TimingSettle, swap...)
}

// Transition enables or disables the transition
// see : https://developer.mozilla.org/en-US/docs/Web/API/View_Transitions_API
func (s *Swap) Transition(transition bool) *Swap {
	s.transition = &transition
	return s
}

// IgnoreTitle enables or disables the Title
func (s *Swap) IgnoreTitle(ignoreTitle bool) *Swap {
	s.ignoreTitle = &ignoreTitle
	return s
}

// FocusScroll enables or disables the focus scroll behaviour
func (s *Swap) FocusScroll(focusScroll bool) *Swap {
	s.focusScroll = &focusScroll
	return s
}

// String returns the string representation of the Swap
func (s *Swap) String() string {
	var parts []string

	parts = append(parts, string(s.style))

	if s.scrolling != nil {
		parts = append(parts, s.scrolling.String())
	}

	if s.transition != nil {
		parts = append(parts, fmt.Sprintf("transition:%s", HxBoolToStr(*s.transition)))
	}

	if s.ignoreTitle != nil {
		parts = append(parts, fmt.Sprintf("ignoreTitle:%s", HxBoolToStr(*s.ignoreTitle)))
	}

	if s.focusScroll != nil {
		parts = append(parts, fmt.Sprintf("focus-scroll:%s", HxBoolToStr(*s.focusScroll)))
	}

	if s.timing != nil {
		parts = append(parts, s.timing.String())
	}

	return strings.Join(parts, " ")
}

const (
	// SwapInnerHTML replaces the inner html of the target element
	SwapInnerHTML SwapStyle = "innerHTML"

	// SwapOuterHTML replaces the entire target element with the response
	SwapOuterHTML SwapStyle = "outerHTML"

	// SwapBeforeBegin insert the response before the target element
	SwapBeforeBegin SwapStyle = "beforebegin"

	// SwapAfterBegin insert the response before the first child of the target element
	SwapAfterBegin SwapStyle = "afterbegin"

	// SwapBeforeEnd insert the response after the last child of the target element
	SwapBeforeEnd SwapStyle = "beforeend"

	// SwapAfterEnd insert the response after the target element
	SwapAfterEnd SwapStyle = "afterend"

	// SwapDelete deletes the target element regardless of the response
	SwapDelete SwapStyle = "delete"

	// SwapNone does not append content from response (out of band items will still be processed).
	SwapNone SwapStyle = "none"
)

type SwapStyle string

func (s SwapStyle) String() string {
	return string(s)
}

const (
	// ScrollingScroll You can also change the scrolling behavior of the target element by using the scroll and show modifiers, both of which take the values top and bottom
	ScrollingScroll SwapScrollingMode = "scroll"

	// ScrollingShow You can also change the scrolling behavior of the target element by using the scroll and show modifiers, both of which take the values top and bottom
	ScrollingShow SwapScrollingMode = "show"
)

type SwapScrollingMode string

func (s SwapScrollingMode) String() string {
	return string(s)
}

const (
	// TimingSwap You can modify the amount of time that htmx will wait after receiving a response to swap the content by including a swap modifier
	TimingSwap SwapTimingMode = "swap"

	// TimingSettle you can modify the time between the swap and the settle logic by including a settle modifier:
	TimingSettle SwapTimingMode = "settle"
)

// SwapTimingMode modifies the amount of time that htmx will wait after receiving a response to swap or settle the content
type SwapTimingMode string

func (s SwapTimingMode) String() string {
	return string(s)
}

const (
	SwapDirectionTop    SwapDirection = "top"
	SwapDirectionBottom SwapDirection = "bottom"
)

// SwapDirection modifies the scrolling behavior of the target element
type SwapDirection string

func (s SwapDirection) String() string {
	return string(s)
}
