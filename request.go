package htmx

import (
	"context"
	"strings"
)

const ContextRequestHeader = "htmx-request-header"

type (
	HxHeaderRequest struct {
		HxBoosted               bool
		HxCurrentURL            string
		HxHistoryRestoreRequest bool
		HxPrompt                string
		HxRequest               bool
		HxTarget                string
		HxTriggerName           string
		HxTrigger               string
	}
)

func (s *Service) HxHeader(ctx context.Context) HxHeaderRequest {
	header := ctx.Value(ContextRequestHeader)

	if val, ok := header.(HxHeaderRequest); ok {
		return val
	}

	return HxHeaderRequest{}
}

func HxStrToBool(str string) bool {
	if strings.EqualFold(str, "true") {
		return true
	}

	return false
}
