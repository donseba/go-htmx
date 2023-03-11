package htmx

import "strings"

func HxStrToBool(str string) bool {
	if strings.EqualFold(str, "true") {
		return true
	}
	return false
}

func HxBoolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
