package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/SevereCloud/gitlabvk/internal"
	"github.com/SevereCloud/gitlabvk/pkg/gitlab"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/callback"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	getSetting         = "get_setting"
	resetToken         = "reset_token"
	notSupportedButton = "not_supported_button"
)

const (
	permissionsMessages = 1 << 12
	permissionsManage   = 1 << 18
)

// ButtonPayload struct
type ButtonPayload struct {
	ButtonType string `json:"button_type,omitempty"`
	Command    string `json:"command"`
	Payload    string `json:"payload,omitempty"`
}

func (b ButtonPayload) String() string {
	raw, _ := json.Marshal(b)
	return string(raw)
}

const maxContentLength = 1e6 * 5 // 5 Mbyte

// Service struct
type Service struct {
	fl *gitlab.FuncList
	vk *api.VK
	cb *callback.Callback

	verify       *internal.Verification
	mtx          sync.Mutex
	storageCache map[string]string

	domain string
}

// NewService return new *Service
func NewService(domain string) *Service {
	s := &Service{
		fl:           gitlab.NewFuncList(),
		vk:           api.NewVK(os.Getenv("GITLABVK_ACCESS_TOKEN")),
		cb:           callback.NewCallback(),
		storageCache: make(map[string]string),
		domain:       domain,
	}
	s.cb.MessageNew(s.MessageNew)

	s.vk.EnableMessagePack()
	s.vk.EnableZstd()

	tokenPerm, err := s.vk.GroupsGetTokenPermissions(api.Params{})
	if err != nil {
		log.WithError(err).Fatal("VK API groups.getTokenPermissions error")
	}

	// Check permissions
	mask := permissionsMessages + permissionsManage
	if mask&tokenPerm.Mask != mask {
		log.WithField("permissions", tokenPerm.Permissions).Fatal("Token bad permissions. Need messages and manage")
	}

	// Храним ключ у пользователя 2e9
	secret := s.getKey(2e9, "secret")
	if secret == "" {
		secret = GenerateRandomString(32)
		s.setKey(2e9, "secret", secret)
	}

	s.verify = internal.NewVerification(secret)

	// s.fl.OnBuild(s.onBuild)
	s.fl.OnIssue(s.onIssue)
	s.fl.OnJob(s.onJob)
	s.fl.OnMergeRequest(s.onMergeRequest)
	s.fl.OnNote(s.onNote)
	s.fl.OnPipeline(s.onPipeline)
	s.fl.OnPush(s.onPush)
	// s.fl.OnSystemHook(s.onSystemHook)
	s.fl.OnTagPush(s.onTagPush)
	s.fl.OnWikiPage(s.onWikiPage)
	s.fl.OnUnknown(s.onUnknow)

	return s
}

// KeyboardBuild return main keyboard
func (s *Service) KeyboardBuild() *object.MessagesKeyboard {
	keyboard := object.NewMessagesKeyboard(true)
	keyboard.AddRow().AddTextButton(
		"Настройки для webhook",
		ButtonPayload{
			Command: getSetting,
		}.String(),
		"",
	)
	keyboard.AddRow().AddTextButton(
		"Сбросить ключ доступа",
		ButtonPayload{
			Command: resetToken,
		}.String(),
		"negative",
	)

	return keyboard
}

func (s *Service) settingMessageBuild(userID int) (text string) {
	u, err := url.Parse(s.domain)
	if err != nil {
		log.WithError(err).Fatal("Invalid domain")
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	u.Path += "/webhook/" + strconv.Itoa(userID)

	text += "URL: " + u.String() + "\n"
	text += "Secret Token: " + s.generateToken(userID) + "\n"

	return
}

// MessageNew callback handler
func (s *Service) MessageNew(_ context.Context, obj events.MessageNewObject) {
	if obj.Message.PeerID > 2e9 {
		return
	}

	var p ButtonPayload
	_ = json.Unmarshal([]byte(obj.Message.Payload), &p)

	var (
		message    string
		attachment string
		keyboard   *object.MessagesKeyboard
	)

	keyboard = s.KeyboardBuild()

	switch p.Command {
	case notSupportedButton:
		log.WithFields(log.Fields{
			"user_id": obj.Message.FromID,
		}).Info("User not support button")

		message = "Ваш клиент не поддерживает эту кнопку"
	case getSetting:
		log.WithFields(log.Fields{
			"user_id": obj.Message.FromID,
		}).Info("User get setting")

		message = "Ваши настройки для Webhooks\n\n"
		message += s.settingMessageBuild(obj.Message.FromID)
	case resetToken:
		log.WithFields(log.Fields{
			"user_id": obj.Message.FromID,
		}).Info("User reset token")

		_ = s.regenerateToken(obj.Message.FromID)
		message = "Токен сброшен. Новые настройки для Webhooks:\n\n"
		message += s.settingMessageBuild(obj.Message.FromID)
	default:
		message = "Ваши настройки для Webhooks\n\n"
		message += s.settingMessageBuild(obj.Message.FromID)
	}

	params := api.Params{
		"peer_id":          obj.Message.PeerID,
		"random_id":        0,
		"message":          message,
		"attachment":       attachment,
		"keyboard":         keyboard,
		"dont_parse_links": true,
		"disable_mentions": true,
	}

	_, err := s.vk.MessagesSend(params)
	if err != nil {
		log.WithError(err).WithFields(log.Fields(params)).Error("Message send error")
	}
}

// Webhook http handler
func (s *Service) Webhook(w http.ResponseWriter, r *http.Request) {
	// Get var from header
	event := gitlab.EventType(r.Header.Get(gitlab.HeaderEvent))
	token := r.Header.Get(gitlab.HeaderToken)

	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(http.StatusText(http.StatusBadRequest)))

		return
	}

	logField := log.Fields{
		"userID":        userID,
		"ip":            internal.GetIP(r),
		"ContentLength": r.ContentLength,
		"event":         event,
	}

	// Check token
	if !s.checkToken(token, userID) {
		log.WithFields(logField).Info("Forbidden")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(http.StatusText(http.StatusForbidden)))

		return
	}

	// Check Content Length
	if r.ContentLength > maxContentLength {
		log.WithFields(logField).Info("Request Entity Too Large")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, _ = w.Write([]byte(http.StatusText(http.StatusRequestEntityTooLarge)))

		return
	}

	// Get data
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(logField).WithError(err).Info("Read body error")
		// NOTE: 499?
		return
	}

	log.Trace(string(data))

	// Handler event
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextUserID, userID)
	ctx = context.WithValue(ctx, contextEventType, event)

	err = s.fl.Handler(ctx, event, data)
	if err != nil {
		log.WithFields(logField).WithError(err).Error("Handler event error")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(http.StatusText(http.StatusBadRequest)))

		return
	}

	// return ok
	log.WithFields(logField).Info("ok")

	_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
}

// CallbackUpdate update callback setting
func (s *Service) CallbackUpdate() {
	u, err := url.Parse(s.domain)
	if err != nil {
		log.WithError(err).Fatal("Invalid domain")
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	u.Path += "/callback"

	urlCallback := u.String()

	err = s.cb.AutoSetting(s.vk, urlCallback)
	if err != nil {
		log.WithError(err).Fatal("Callback AutoSetting")
	}
}

// Callback func
func (s *Service) Callback(w http.ResponseWriter, r *http.Request) {
	s.cb.HandleFunc(w, r)
}

func init() {
	// Flags
	lvl := flag.String("level", "info", "logger level")
	flag.Parse()

	log.WithField("level", lvl).Debug("logger level")

	// Logrus level
	level, err := log.ParseLevel(*lvl)
	if err != nil {
		log.Warn(err)
		level = log.InfoLevel
	}

	log.SetLevel(level)
}

func main() {
	s := NewService(os.Getenv("GITLABVK_DOMAIN"))

	addr := os.Getenv("GITLABVK_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	// Router setting
	r := mux.NewRouter()
	r.HandleFunc("/webhook/{id}", s.Webhook)
	r.HandleFunc("/callback", s.Callback)
	http.Handle("/", r)
	log.Printf("Start server on %s", addr)

	// паралельно обновляем callback
	go s.CallbackUpdate()

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.WithError(err).Fatal("ListenAndServe error")
	}
}
