package models

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

type VideoStreamUserData struct {
	ResourceID      string
	ItemID          string
	SubtitleID      string
	AudioID         string
	AcceptLangTags  []language.Tag
	FallbackLangTag language.Tag
}

func NewVideoStreamUserData(resourceID string, itemID string) *VideoStreamUserData {
	return &VideoStreamUserData{
		ResourceID:      resourceID,
		ItemID:          itemID,
		FallbackLangTag: language.English,
	}
}

func (s *VideoStreamUserData) FetchSessionData(c *gin.Context) {
	session := sessions.Default(c)
	var subtitleID, audioID string
	audioKey := s.makeKey(s.ResourceID, s.ItemID, "audio")
	subtitleKey := s.makeKey(s.ResourceID, s.ItemID, "subtitle")
	if session.Get(subtitleKey) != nil {
		subtitleID = session.Get(subtitleKey).(string)
	}
	if session.Get(audioKey) != nil {
		audioID = session.Get(audioKey).(string)
	}
	accept := c.GetHeader("Accept-Language")
	tags, _, err := language.ParseAcceptLanguage(accept)
	if err != nil {
		tags = []language.Tag{language.English}
	}
	s.AudioID = audioID
	s.SubtitleID = subtitleID
	s.AcceptLangTags = tags
}

func (s *VideoStreamUserData) makeKey(resourceID string, itemID string, name string) string {
	return fmt.Sprintf("%v_%v_%v_id", resourceID, itemID, name)
}

func (s *VideoStreamUserData) UpdateSessionData(c *gin.Context) error {
	session := sessions.Default(c)
	audioKey := s.makeKey(s.ResourceID, s.ItemID, "audio")
	subtitleKey := s.makeKey(s.ResourceID, s.ItemID, "subtitle")
	if s.SubtitleID == "" {
		session.Delete(subtitleKey)
	} else {
		session.Set(subtitleKey, s.SubtitleID)
	}
	if s.AudioID == "" {
		session.Delete(audioKey)
	} else {
		session.Set(audioKey, s.AudioID)
	}
	return session.Save()
}
