package action

import (
	"context"
	m "github.com/webtor-io/web-ui-v2/services/models"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	sv "github.com/webtor-io/web-ui-v2/services"
	j "github.com/webtor-io/web-ui-v2/services/job"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type PostArgs struct {
	ResourceID string
	ItemID     string
	Claims     *sv.Claims
}

type TrackPutArgs struct {
	ID         string `json:"id"`
	ResourceID string `json:"resourceID"`
	ItemID     string `json:"itemID"`
}

type PostData struct {
	w.ErrorData
	Job  *sv.Job
	Args *PostArgs
}

type Handler struct {
	*w.TemplateHandler
	*w.ClaimsHandler
	jobs *j.Handler
}

func RegisterHandler(c *cli.Context, r *gin.Engine, re multitemplate.Renderer, jobs *j.Handler) {
	h := &Handler{
		TemplateHandler: w.NewTemplateHandler(c, re),
		ClaimsHandler:   w.NewClaimsHandler(),
		jobs:            jobs,
	}
	r.POST("/download-file", func(c *gin.Context) {
		h.post(c, "download")
	})
	r.POST("/download-dir", func(c *gin.Context) {
		h.post(c, "download")
	})
	r.POST("/preview-image", func(c *gin.Context) {
		h.post(c, "preview-image")
	})
	r.POST("/stream-audio", func(c *gin.Context) {
		h.post(c, "stream-audio")
	})
	r.POST("/stream-video", func(c *gin.Context) {
		h.post(c, "stream-video")
	})
	r.PUT("/stream-video/subtitle", func(c *gin.Context) {
		a := TrackPutArgs{}
		if err := c.BindJSON(&a); err != nil {
			c.Error(err)
			return
		}
		vsud := m.NewVideoStreamUserData(a.ResourceID, a.ItemID)
		vsud.SubtitleID = a.ID
		if err := vsud.UpdateSessionData(c); err != nil {
			c.Error(err)
		}
	})
	r.PUT("/stream-video/audio", func(c *gin.Context) {
		a := TrackPutArgs{}
		if err := c.BindJSON(&a); err != nil {
			c.Error(err)
			return
		}
		vsud := m.NewVideoStreamUserData(a.ResourceID, a.ItemID)
		vsud.AudioID = a.ID
		if err := vsud.UpdateSessionData(c); err != nil {
			c.Error(err)
		}
	})
	h.RegisterTemplate(
		"action/post",
		[]string{"async"},
		[]string{},
		template.FuncMap{},
	)
	h.RegisterTemplate(
		"action/preview_image",
		[]string{"async"},
		[]string{},
		template.FuncMap{},
	)
	h.RegisterTemplate(
		"action/stream_video",
		[]string{"async"},
		[]string{},
		template.FuncMap{
			"getDurationSec": GetDurationSec,
			"getAudioTracks": GetAudioTracks,
			"getSubtitles":   GetSubtitles,
		},
	)
	h.RegisterTemplate(
		"action/stream_audio",
		[]string{"async"},
		[]string{},
		template.FuncMap{
			"getDurationSec": GetDurationSec,
		},
	)
}

func (s *Handler) bindPostArgs(c *gin.Context) (*PostArgs, error) {
	rID, ok := c.GetPostFormArray("resource-id")
	if !ok {
		return nil, errors.Errorf("no resource id provided")
	}
	iID, ok := c.GetPostFormArray("item-id")
	if !ok {
		return nil, errors.Errorf("no item id provided")
	}
	return &PostArgs{
		ResourceID: rID[0],
		ItemID:     iID[0],
		Claims:     s.MakeClaims(c),
	}, nil
}

func (s *Handler) post(c *gin.Context, action string) {
	var (
		d    PostData
		err  error
		args *PostArgs
		job  *sv.Job
	)
	post := s.MakeTemplate(c, "action/post", &d)
	args, err = s.bindPostArgs(c)
	if err != nil {
		d.Err = errors.Wrap(err, "wrong args provided")
		post.R(http.StatusBadRequest)
		return
	}
	d.Args = args
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	job, err = s.jobs.Action(ctx, c, args.Claims, args.ResourceID, args.ItemID, action)
	if err != nil {
		d.Err = errors.Wrap(err, "failed to start downloading")
		post.R(http.StatusBadRequest)
		return
	}
	d.Job = job
	post.R(http.StatusOK)
}
