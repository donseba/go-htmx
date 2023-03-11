package htmx

import (
	"context"
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

func (s *HTMX) HxHeader(ctx context.Context) HxHeaderRequest {
	header := ctx.Value(ContextRequestHeader)

	if val, ok := header.(HxHeaderRequest); ok {
		return val
	}

	return HxHeaderRequest{}
}
