package main

import (
	"context"

	"github.com/SevereCloud/gitlabvk/pkg/gitlab"
)

type contextKey int

const (
	contextUserID contextKey = iota
	contextEventType
)

// getUserID return userID
func getUserID(ctx context.Context) int {
	if ctx != nil {
		if hc, ok := ctx.Value(contextUserID).(int); ok {
			return hc
		}
	}

	return 0
}

// getEventType return gitlab event type
func getEventType(ctx context.Context) gitlab.EventType {
	if ctx != nil {
		if hc, ok := ctx.Value(contextEventType).(gitlab.EventType); ok {
			return hc
		}
	}

	return ""
}
