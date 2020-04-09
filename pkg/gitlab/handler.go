// Package gitlab for gitlab events
package gitlab

import (
	"context"
	"encoding/json"
)

// FuncBuild function for handler
type FuncBuild func(context.Context, EventBuild)

// FuncIssue function for handler
type FuncIssue func(context.Context, EventIssue)

// FuncJob function for handler
type FuncJob func(context.Context, EventJob)

// FuncMergeRequest function for handler
type FuncMergeRequest func(context.Context, EventMergeRequest)

// FuncNote function for handler
type FuncNote func(context.Context, EventNote)

// FuncPipeline function for handler
type FuncPipeline func(context.Context, EventPipeline)

// FuncPush function for handler
type FuncPush func(context.Context, EventPush)

// FuncSystemHook function for handler
// type FuncSystemHook func(context.Context, EventSystemHook)

// FuncTagPush function for handler
type FuncTagPush func(context.Context, EventTagPush)

// FuncWikiPage function for handler
type FuncWikiPage func(context.Context, EventWikiPage)

// FuncUnknown function for handler unknown event
type FuncUnknown func(context.Context, interface{})

// FuncList struct
type FuncList struct {
	build        []FuncBuild
	issue        []FuncIssue
	job          []FuncJob
	mergeRequest []FuncMergeRequest
	note         []FuncNote
	pipeline     []FuncPipeline
	push         []FuncPush
	tagPush      []FuncTagPush
	wikiPage     []FuncWikiPage
	unknown      []FuncUnknown
	// systemHook   []FuncSystemHook
}

// NewFuncList return FuncList
func NewFuncList() *FuncList {
	return &FuncList{}
}

// Handler gitlab events
func (fl FuncList) Handler(ctx context.Context, event EventType, data []byte) error { // nolint:gocyclo
	switch event {
	case EventTypeBuild:
		var obj EventBuild
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.build {
			f(ctx, obj)
		}
	case EventTypeIssue:
		var obj EventIssue
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.issue {
			f(ctx, obj)
		}
	case EventTypeJob:
		var obj EventJob
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.job {
			f(ctx, obj)
		}
	case EventTypeMergeRequest:
		var obj EventMergeRequest
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.mergeRequest {
			f(ctx, obj)
		}
	case EventTypeNote:
		var obj EventNote
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.note {
			f(ctx, obj)
		}
	case EventTypePipeline:
		var obj EventPipeline
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.pipeline {
			f(ctx, obj)
		}
	case EventTypePush:
		var obj EventPush
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.push {
			f(ctx, obj)
		}
	// case EventTypeSystemHook:
	// 	var obj EventSystemHook
	// 	if err := json.Unmarshal(data, &obj); err != nil {
	// 		return err
	// 	}

	// 	for _, f := range fl.systemHook {
	// 		f(ctx, obj)
	// 	}
	case EventTypeTagPush:
		var obj EventTagPush
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.tagPush {
			f(ctx, obj)
		}
	case EventTypeWikiPage:
		var obj EventWikiPage
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.wikiPage {
			f(ctx, obj)
		}
	default:
		var obj interface{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}

		for _, f := range fl.unknown {
			f(ctx, obj)
		}
	}

	return nil
}

// OnBuild event handler
func (fl *FuncList) OnBuild(f FuncBuild) {
	fl.build = append(fl.build, f)
}

// OnIssue event handler
func (fl *FuncList) OnIssue(f FuncIssue) {
	fl.issue = append(fl.issue, f)
}

// OnJob event handler
func (fl *FuncList) OnJob(f FuncJob) {
	fl.job = append(fl.job, f)
}

// OnMergeRequest event handler
func (fl *FuncList) OnMergeRequest(f FuncMergeRequest) {
	fl.mergeRequest = append(fl.mergeRequest, f)
}

// OnNote event handler
func (fl *FuncList) OnNote(f FuncNote) {
	fl.note = append(fl.note, f)
}

// OnPipeline event handler
func (fl *FuncList) OnPipeline(f FuncPipeline) {
	fl.pipeline = append(fl.pipeline, f)
}

// OnPush event handler
func (fl *FuncList) OnPush(f FuncPush) {
	fl.push = append(fl.push, f)
}

// OnSystemHook event handler
// func (fl *FuncList) OnSystemHook(f FuncSystemHook) {
// 	fl.systemHook = append(fl.systemHook, f)
// }

// OnTagPush event handler
func (fl *FuncList) OnTagPush(f FuncTagPush) {
	fl.tagPush = append(fl.tagPush, f)
}

// OnWikiPage event handler
func (fl *FuncList) OnWikiPage(f FuncWikiPage) {
	fl.wikiPage = append(fl.wikiPage, f)
}

// OnUnknown event handler
func (fl *FuncList) OnUnknown(f FuncUnknown) {
	fl.unknown = append(fl.unknown, f)
}
