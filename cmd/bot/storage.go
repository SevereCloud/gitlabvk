package main

import (
	"fmt"

	"github.com/SevereCloud/vksdk/api"
	log "github.com/sirupsen/logrus"
)

const prefixKey = "gitlabvk_"

// keys
const (
	pipelineMessageID = "pipeline_message_id"
	pipelineLastID    = "pipeline_last_id"
)

func (s *Service) getKey(userID int, key string) string {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	// check cache
	if v, ok := s.storageCache[fmt.Sprintf("%d_%s", userID, key)]; ok {
		return v
	}

	// get from VK API
	r, err := s.vk.StorageGet(api.Params{
		"key":     prefixKey + key,
		"user_id": userID,
	})
	if err != nil {
		log.WithError(err).Fatal("VK API storage.get")
	}

	// save cache
	s.storageCache[fmt.Sprintf("%d_%s", userID, key)] = r[0].Value

	return r[0].Value
}

func (s *Service) setKey(userID int, key, value string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.getKey(userID, key) != value {
		// save cache
		s.storageCache[fmt.Sprintf("%d_%s", userID, key)] = value

		_, err := s.vk.StorageSet(api.Params{
			"key":     prefixKey + key,
			"value":   value,
			"user_id": userID,
		})
		if err != nil {
			log.WithError(err).Fatal("VK API storage.set")
		}
	}
}
