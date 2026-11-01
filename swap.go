package htmx

import (
	"fmt"
	"strconv"
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

// htmx4String returns the timing in the syntax of htmx 4. It always writes the duration, as htmx 4 reads
// a bare swap or settle as true, and writes it in whole seconds or in milliseconds, which htmx 4 parses.
func (s *SwapTiming) htmx4String() string {
	if s.duration%time.Second == 0 {
		return fmt.Sprintf("%s:%ds", s.mode, s.duration/time.Second)
	}

	ms := strconv.FormatFloat(float64(s.duration)/float64(time.Millisecond), 'f', -1, 64)

	return fmt.Sprintf("%s:%sms", s.mode, ms)
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

// htmx4String returns the scrolling in the syntax of htmx 4, which takes the target in a modifier of its own:
// show:top showTarget:"#element".
func (s *SwapScrolling) htmx4String() string {
	var out string

	out = string(s.mode)

	if s.direction != "" {
		out += ":" + s.direction.String()
	}

	if s.target != "" {
		out += " " + string(s.mode) + "Target:" + hconQuote(s.target)
	}

	return out
}

// hconQuote quotes a value for the HCON syntax of htmx 4, which has no escapes:
// in double quotes, or in single quotes when the value holds a double quote.
func hconQuote(value string) string {
	if strings.Contains(value, `"`) {
		return "'" + value + "'"
	}

	return `"` + value + `"`
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

// String returns the string representation of the Swap, in the syntax of htmx 2
func (s *Swap) String() string {
	return s.StringFor(Version2)
}

// StringFor returns the string representation of the Swap in the syntax of the given htmx version.
// htmx 4 writes focusScroll rather than focus-scroll, and the target of scroll and show in a modifier of its own:
// show:top showTarget:"#element" rather than show:#element:top.
func (s *Swap) StringFor(version Version) string {
	var parts []string

	parts = append(parts, string(s.style))

	if s.scrolling != nil {
		if version == Version4 {
			parts = append(parts, s.scrolling.htmx4String())
		} else {
			parts = append(parts, s.scrolling.String())
		}
	}

	if s.transition != nil {
		parts = append(parts, fmt.Sprintf("transition:%s", HxBoolToStr(*s.transition)))
	}

	if s.ignoreTitle != nil {
		parts = append(parts, fmt.Sprintf("ignoreTitle:%s", HxBoolToStr(*s.ignoreTitle)))
	}

	if s.focusScroll != nil {
		if version == Version4 {
			parts = append(parts, fmt.Sprintf("focusScroll:%s", HxBoolToStr(*s.focusScroll)))
		} else {
			parts = append(parts, fmt.Sprintf("focus-scroll:%s", HxBoolToStr(*s.focusScroll)))
		}
	}

	if s.timing != nil {
		if version == Version4 {
			parts = append(parts, s.timing.htmx4String())
		} else {
			parts = append(parts, s.timing.String())
		}
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

	// SwapInnerMorph morphs the children of the target element into the response (htmx 4)
	SwapInnerMorph SwapStyle = "innerMorph"

	// SwapOuterMorph morphs the target element into the response (htmx 4)
	SwapOuterMorph SwapStyle = "outerMorph"

	// SwapOuterSync replaces the children of the target element and copies the attributes of the response's element (htmx 4)
	SwapOuterSync SwapStyle = "outerSync"

	// SwapTextContent replaces the text content of the target element with the text of the response (htmx 4)
	SwapTextContent SwapStyle = "textContent"
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
