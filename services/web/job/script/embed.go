package script

import (
	"context"
	"crypto/sha1"
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
	"github.com/webtor-io/web-ui-v2/services/template"

	ra "github.com/webtor-io/rest-api/services"
)

var (
	sampleReg = regexp.MustCompile("/sample/i")
)

type EmbedSettings struct {
	StreamSettings
	Version string `json:"version"`
	Magnet  string `json:"magnet"`
}

type EmbedScript struct {
	api      *api.Api
	claims   *api.Claims
	settings *EmbedSettings
	file     string
	tb       template.Builder
	c        *gin.Context
}

func NewEmbedScript(tb template.Builder, c *gin.Context, api *api.Api, claims *api.Claims, settings *EmbedSettings, file string) *EmbedScript {
	return &EmbedScript{
		api:      api,
		claims:   claims,
		settings: settings,
		file:     file,
		tb:       tb,
		c:        c,
	}
}

func (s *EmbedScript) Run(j *job.Job) (err error) {
	ls, _, err := Load(s.api, s.claims, &LoadArgs{
		Query: s.settings.Magnet,
	})
	if err != nil {
		return err
	}
	err = ls.Run(j)
	if err != nil {
		return err
	}
	id := j.Context.Value("respID").(string)
	i, err := s.getBestItem(j, id)
	if err != nil {
		return err
	}
	var action string
	if i.MediaFormat == ra.Video {
		action = "stream-video"
	} else if i.MediaFormat == ra.Audio {
		action = "stream-audio"
	}
	as, _ := Action(s.tb, s.api, s.claims, s.c, id, i.ID, action, &s.settings.StreamSettings)
	err = as.Run(j)
	if err != nil {
		return err
	}
	return
}

func (s *EmbedScript) getBestItem(j *job.Job, id string) (i *ra.ListItem, err error) {
	j.InProgress("searching for stream content")
	ctx, cancel := context.WithTimeout(j.Context, 10*time.Second)
	defer cancel()
	l, err := s.api.ListResourceContent(ctx, s.claims, id, &api.ListResourceContentArgs{
		Output: api.OutputTree,
	})
	if err != nil {
		return nil, j.Error(err, "failed to list resource content")
	}
	i = s.findBestItem(l)
	if err != nil {
		return nil, j.Error(err, "failed to find stream content")
	}
	if i == nil {
		return nil, j.Error(err, "failed to find stream content")
	}
	j.Done()
	return
}

func (s *EmbedScript) findBestItem(l *ra.ListResponse) *ra.ListItem {
	for _, v := range l.Items {
		if v.MediaFormat == ra.Video && !sampleReg.MatchString(v.Name) {
			return &v
		}
	}
	for _, v := range l.Items {
		if v.MediaFormat == ra.Audio && !sampleReg.MatchString(v.Name) {
			return &v
		}
	}
	for _, v := range l.Items {
		if v.Type == ra.ListTypeFile {
			return &v
		}
	}
	return nil
}

func Embed(tb template.Builder, c *gin.Context, api *api.Api, claims *api.Claims, settings *EmbedSettings, file string) (r job.Runnable, hash string, err error) {
	hash = fmt.Sprintf("%x", sha1.Sum([]byte(claims.Role+"/"+claims.SessionID+"/"+file+"/"+fmt.Sprintf("%+v", settings))))
	r = NewEmbedScript(tb, c, api, claims, settings, file)
	return
}
