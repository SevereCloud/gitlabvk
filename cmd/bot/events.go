package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"

	"github.com/SevereCloud/gitlabvk/pkg/gitlab"
)

const maxAttemptSendMessage = 3

func baseRef(ref string) string {
	a := strings.Split(ref, "/")
	if len(a) > 2 {
		return strings.Join(a[2:], "/")
	}

	return ref
}

func (s *Service) onPush(ctx context.Context, e gitlab.EventPush) {
	branch := baseRef(e.Ref)
	message := fmt.Sprintf("🛠 %s pushed to %s#%s\n\n", e.UserName, e.Project.Name, branch)

	for _, commit := range e.Commits {
		message += commit.Message + "\n"
	}

	keyboard := object.NewMessagesKeyboardInline()
	link := ""

	if e.After != gitlab.NullSHA {
		switch {
		case e.Before == gitlab.NullSHA && len(e.Commits) > 0:
			link = e.Commits[0].URL
		default:
			link = fmt.Sprintf("%s/-/compare/%s...%s", e.Repository.Homepage, e.Before[:8], e.After[:8])
		}
	}

	if link != "" {
		keyboard.AddRow()
		keyboard.AddOpenLinkButton(link, "Changes", "")
	}

	userID := getUserID(ctx)
	s.sendMessage(userID, message, keyboard)
}

func (s *Service) onTagPush(ctx context.Context, e gitlab.EventTagPush) {
	message := ""
	keyboard := object.NewMessagesKeyboardInline()
	tag := baseRef(e.Ref)

	if e.CheckoutSHA == "" {
		message = fmt.Sprintf("🏷️ remove tag %s#%s\n\n%s", e.Project.Name, tag, e.Message)
	} else {
		message = fmt.Sprintf("🏷️ new tag %s#%s\n\n%s", e.Project.Name, tag, e.Message)
		link := fmt.Sprintf("%s/-/tags/%s", e.Repository.Homepage, tag)
		keyboard.AddRow()
		keyboard.AddOpenLinkButton(link, "Changes", "")
	}

	userID := getUserID(ctx)
	s.sendMessage(userID, message, keyboard)
}

func (s *Service) onIssue(ctx context.Context, e gitlab.EventIssue) {
	message := fmt.Sprintf(
		"🐛 %s %s issue: %s#%d\n%s\n\n%s",
		e.User.Name,
		e.ObjectAttributes.Action,
		e.Project.Name, e.ObjectAttributes.IID,
		e.ObjectAttributes.Title,
		e.ObjectAttributes.Description,
	)

	link := e.ObjectAttributes.URL
	keyboard := object.NewMessagesKeyboardInline()
	keyboard.AddRow()
	keyboard.AddOpenLinkButton(link, "Open issue", "")

	userID := getUserID(ctx)
	s.sendMessage(userID, message, keyboard)
}

func (s *Service) onNote(ctx context.Context, e gitlab.EventNote) {
	message := "💬 " + e.User.Name + " "

	switch e.ObjectAttributes.NoteableType {
	case gitlab.NoteableTypeIssue:
		message += fmt.Sprintf(
			"write comment to issue %s#%d\n\n%s",
			e.Project.Name, e.Issue.IID,
			e.ObjectAttributes.Note,
		)
	case gitlab.NoteableTypeCommit:
		message += fmt.Sprintf(
			"write comment to commit %s#%s\n\n%s",
			e.Project.Name, e.Commit.ID[:8],
			e.ObjectAttributes.Note,
		)
	case gitlab.NoteableTypeMergeRequest:
		message += fmt.Sprintf(
			"write comment to MR %s#%d\n\n%s",
			e.Project.Name, e.MergeRequest.IID,
			e.ObjectAttributes.Note,
		)
	case gitlab.NoteableTypeSnippet:
		message += fmt.Sprintf(
			"write comment to snippet %s $%d\n\n%s",
			e.Project.Name, e.Snippet.ID,
			e.ObjectAttributes.Note,
		)
	default:
		log.WithField("noteable", e.ObjectAttributes.NoteableType).Warn("Not found noteable type")
		message = fmt.Sprintf(
			"💬 %s write comment to %s\n\n%s",
			e.User.Name,
			e.Project.Name,
			e.ObjectAttributes.Note,
		)
	}

	link := e.ObjectAttributes.URL
	keyboard := object.NewMessagesKeyboardInline()
	keyboard.AddRow()
	keyboard.AddOpenLinkButton(link, "Open comment", "")

	userID := getUserID(ctx)
	s.sendMessage(userID, message, keyboard)
}

func (s *Service) onMergeRequest(ctx context.Context, e gitlab.EventMergeRequest) {
	message := fmt.Sprintf(
		"🔀 %s %s MR: %s#%d\n%s\n\n%s",
		e.User.Name,
		e.ObjectAttributes.Action,
		e.Project.Name, e.ObjectAttributes.IID,
		e.ObjectAttributes.Title,
		e.ObjectAttributes.Description,
	)

	link := e.ObjectAttributes.URL
	keyboard := object.NewMessagesKeyboardInline()
	keyboard.AddRow()
	keyboard.AddOpenLinkButton(link, "Open", "")

	userID := getUserID(ctx)
	s.sendMessage(userID, message, keyboard)
}

func (s *Service) onJob(ctx context.Context, e gitlab.EventJob) {
	var message string

	switch e.BuildStatus {
	case gitlab.StatusCreated:
		// Ignore it
		return
	case gitlab.StatusRunning:
		message += "⌚"
	case gitlab.StatusCanceled:
		message += "🚫"
	case gitlab.StatusFailed:
		message += "🗙"
	case gitlab.StatusSuccess:
		message += "✅"
	default:
		message += "💼"
	}

	message += fmt.Sprintf(
		" %s %s %s\n",
		e.BuildStage,
		e.BuildName,
		e.BuildStatus,
	)

	// link := fmt.Sprintf("%s/pipelines/%d", e.Repository.Homepage, e.PipelineID)

	// keyboard := object.NewMessagesKeyboardInline()
	// keyboard.AddRow()
	// keyboard.AddOpenLinkButton(link, "Open pipeline", "")

	userID := getUserID(ctx)
	if s.getKey(userID, pipelineLastID) == strconv.Itoa(e.PipelineID) {
		log.Debug("job in pipe")
		s.sendPipelineMessage(userID, message, nil)
	} else {
		log.Debug("job new")
		s.setKey(userID, pipelineLastID, strconv.Itoa(e.PipelineID))
		id := s.sendMessage(userID, message, nil)
		s.setKey(userID, pipelineMessageID, strconv.Itoa(id))
	}
}

func (s *Service) onPipeline(ctx context.Context, e gitlab.EventPipeline) {
	var message string

	switch e.ObjectAttributes.Status {
	case gitlab.StatusPending:
		message += "⏸️"
	case gitlab.StatusRunning:
		message += "▶️"
	case gitlab.StatusCanceled:
		message += "🚫"
	case gitlab.StatusFailed:
		message += "🗙"
	case gitlab.StatusSuccess:
		message += "✅"
	default:
		message += "💼"
	}

	message += fmt.Sprintf(
		" pipeline #%d %s\n",
		e.ObjectAttributes.ID,
		e.ObjectAttributes.Status,
	)

	userID := getUserID(ctx)

	link := fmt.Sprintf("%s/pipelines/%d", e.Project.WebURL, e.ObjectAttributes.ID)

	keyboard := object.NewMessagesKeyboardInline()
	keyboard.AddRow()
	keyboard.AddOpenLinkButton(link, "Open pipeline", "")

	if s.getKey(userID, pipelineLastID) == strconv.Itoa(e.ObjectAttributes.ID) {
		if e.ObjectAttributes.Status == gitlab.StatusFailed {
			s.sendMessage(userID, message, keyboard)
		} else {
			s.sendPipelineMessage(userID, message, nil)
		}
	} else {
		s.setKey(userID, pipelineLastID, strconv.Itoa(e.ObjectAttributes.ID))
		id := s.sendMessage(userID, message, nil)
		s.setKey(userID, pipelineMessageID, strconv.Itoa(id))
	}
}

func (s *Service) onWikiPage(ctx context.Context, e gitlab.EventWikiPage) {
	message := fmt.Sprintf(
		"📙 %s %s page %s\n%s\n\n%s",
		e.User.Name,
		e.ObjectAttributes.Action,
		e.Project.Name,
		e.ObjectAttributes.Title,
		e.ObjectAttributes.Message,
	)

	link := e.ObjectAttributes.URL
	keyboard := object.NewMessagesKeyboardInline()
	keyboard.AddRow()
	keyboard.AddOpenLinkButton(link, "Open page", "")

	userID := getUserID(ctx)
	s.sendMessage(userID, message, keyboard)
}

func (s *Service) onUnknow(ctx context.Context, e interface{}) {
	userID := getUserID(ctx)
	event := getEventType(ctx)
	message := fmt.Sprintf("❓ Unknown event %s", event)

	log.WithFields(log.Fields{
		"userID": userID,
		"event":  event,
	}).Warn("Unknown event!")

	s.sendMessage(userID, message, nil)
}

func (s *Service) sendMessage(peerID int, message string, keyboard *object.MessagesKeyboard) int {
	b := params.NewMessagesSendBuilder()
	b.PeerID(peerID)
	b.RandomID(0)
	b.Message(message)
	b.DisableMentions(true)
	b.DontParseLinks(true)

	if keyboard != nil {
		b.Keyboard(keyboard)
	}

	attempt := 0
	for attempt < maxAttemptSendMessage {
		retry := false

		id, err := s.vk.MessagesSend(b.Params)

		var errCode api.ErrorType

		errors.As(err, &errCode)

		switch errCode {
		case api.ErrNoType:
			if err != nil {
				log.WithError(err).WithFields(log.Fields(b.Params)).Error("Messages send error")
				return 0
			}

			return id
		case api.ErrTooMany:
			log.WithError(err).WithFields(log.Fields(b.Params)).Warn("Retry send message")

			retry = true
		case api.ErrServer:
			log.WithError(err).WithFields(log.Fields(b.Params)).Warn("Retry send message")

			retry = true
		case api.ErrMessagesDenySend:
			log.WithError(err).WithFields(log.Fields(b.Params)).Info("Messages deny send")
		default:
			log.WithError(err).WithFields(log.Fields(b.Params)).Error("Messages send error")
		}

		if !retry {
			break
		}

		time.Sleep(time.Second)
		attempt++
	}

	return 0
}

func (s *Service) addMessage(peerID, messageID int, message string) error {
	if messageID == 0 {
		return fmt.Errorf("message id=0")
	}

	b0 := params.NewMessagesGetByIDBuilder()
	b0.MessageIDs([]int{messageID})

	resp, err := s.vk.MessagesGetByID(b0.Params)
	if err != nil {
		return err
	}

	if len(resp.Items) == 0 {
		return fmt.Errorf("message not found")
	}

	message = resp.Items[0].Text + "\n" + message

	b := params.NewMessagesEditBuilder()
	b.PeerID(peerID)
	b.MessageID(messageID)
	b.Message(message)
	b.DontParseLinks(true)

	_, err = s.vk.MessagesEdit(b.Params)

	return err
}

func (s *Service) sendPipelineMessage(peerID int, message string, keyboard *object.MessagesKeyboard) {
	id, _ := strconv.Atoi(s.getKey(peerID, pipelineMessageID))

	if id != 0 {
		log.WithField("id", id).Debug("addMessage")

		err := s.addMessage(peerID, id, message)
		if err == nil {
			return
		}
	}

	id = s.sendMessage(peerID, message, keyboard)
	if id != 0 {
		s.setKey(peerID, pipelineMessageID, strconv.Itoa(id))
	}
}
