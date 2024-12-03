package script

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/pkg/errors"
	"github.com/webtor-io/web-ui-v2/services/claims"
	"github.com/webtor-io/web-ui-v2/services/models"
	"io"
	"net/http"
	"regexp"
	"strings"
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

type EmbedScript struct {
	api        *api.Api
	apiClaims  *api.Claims
	settings   *models.EmbedSettings
	file       string
	tb         template.Builder
	c          *gin.Context
	hCl        *http.Client
	userClaims *claims.Data
}

func NewEmbedScript(tb template.Builder, hCl *http.Client, c *gin.Context, api *api.Api, apiClaims *api.Claims, userClaims *claims.Data, settings *models.EmbedSettings, file string) *EmbedScript {
	return &EmbedScript{
		api:        api,
		apiClaims:  apiClaims,
		userClaims: userClaims,
		settings:   settings,
		file:       file,
		tb:         tb,
		c:          c,
		hCl:        hCl,
	}
}

func (s *EmbedScript) makeLoadArgs(settings *models.EmbedSettings) (*LoadArgs, error) {
	la := &LoadArgs{}
	if settings.TorrentURL != "" {
		resp, err := s.hCl.Get(settings.TorrentURL)
		if err != nil {
			return nil, err
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		la.File = body
	} else if settings.Magnet != "" {
		la.Query = settings.Magnet
	}
	return la, nil
}

func (s *EmbedScript) Run(j *job.Job) (err error) {
	args, err := s.makeLoadArgs(s.settings)
	if err != nil {
		return
	}
	ls, _, err := Load(s.api, s.apiClaims, args)
	if err != nil {
		return err
	}
	err = ls.Run(j)
	if err != nil {
		return err
	}
	id := j.Context.Value("respID").(string)
	i, err := s.getBestItem(j, id, s.settings)
	if err != nil {
		return err
	}
	var action string
	if i.MediaFormat == ra.Video {
		action = "stream-video"
	} else if i.MediaFormat == ra.Audio {
		action = "stream-audio"
	}
	as, _ := Action(s.tb, s.api, s.apiClaims, s.userClaims, s.c, id, i.ID, action, &s.settings.StreamSettings)
	err = as.Run(j)
	if err != nil {
		return err
	}
	return
}

func (s *EmbedScript) getBestItem(j *job.Job, id string, settings *models.EmbedSettings) (i *ra.ListItem, err error) {
	j.InProgress("searching for stream content")
	ctx, cancel := context.WithTimeout(j.Context, 10*time.Second)
	defer cancel()
	pwd := settings.PWD
	file := settings.File
	if settings.Path != "" {
		parts := strings.Split(settings.Path, "/")
		file = parts[len(parts)-1]
		pwd = strings.Join(parts[:len(parts)-1], "/")
	}
	l, err := s.api.ListResourceContent(ctx, s.apiClaims, id, &api.ListResourceContentArgs{
		Path:   pwd,
		Output: api.OutputTree,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list resource content")
	}
	if len(l.Items) == 1 && l.Items[0].Type == ra.ListTypeDirectory {
		l, err = s.api.ListResourceContent(ctx, s.apiClaims, id, &api.ListResourceContentArgs{
			Path:   l.Items[0].PathStr,
			Output: api.OutputTree,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list resource content")
		}
	}
	if file != "" {
		for _, f := range l.Items {
			if f.Name == file {
				i = &f
				break
			}
		}
	} else {
		i = s.findBestItem(l)
	}
	if i == nil {
		return nil, errors.Wrap(err, "failed to find stream content")
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

func Embed(tb template.Builder, hCl *http.Client, c *gin.Context, api *api.Api, apiClaims *api.Claims, userClaims *claims.Data, settings *models.EmbedSettings, file string) (r job.Runnable, hash string, err error) {
	hash = fmt.Sprintf("%x", sha1.Sum([]byte(apiClaims.Role+"/"+fmt.Sprintf("%+v", settings))))
	r = NewEmbedScript(tb, hCl, c, api, apiClaims, userClaims, settings, file)
	return
}
