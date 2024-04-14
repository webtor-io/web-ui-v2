package action

import (
	"context"
	"html/template"
	"net/http"
	"time"

	m "github.com/webtor-io/web-ui-v2/services/models"

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
	Job  *sv.Job
	Args *PostArgs
}

type Handler struct {
	*w.ClaimsHandler
	jobs *j.Handler
	tm   *w.TemplateManager
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *w.TemplateManager, jobs *j.Handler, uc *sv.UserClaims) {
	h := &Handler{
		tm:            tm,
		ClaimsHandler: w.NewClaimsHandler(uc),
		jobs:          jobs,
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

	h.tm.RegisterViewsWithFuncs("action/*", template.FuncMap{
		"getDurationSec": GetDurationSec,
		"getAudioTracks": GetAudioTracks,
		"getSubtitles":   GetSubtitles,
	})
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
	claims, err := s.MakeClaims(c)
	if err != nil {
		return nil, err
	}
	return &PostArgs{
		ResourceID: rID[0],
		ItemID:     iID[0],
		Claims:     claims,
	}, nil
}

func (s *Handler) post(c *gin.Context, action string) {
	var (
		d    PostData
		err  error
		args *PostArgs
		job  *sv.Job
	)
	postTpl := s.tm.MakeTemplate("action/post")
	args, err = s.bindPostArgs(c)
	if err != nil {
		postTpl.HTMLWithErr(errors.Wrap(err, "wrong args provided"), http.StatusBadRequest, c, d)
		return
	}
	d.Args = args
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	job, err = s.jobs.Action(ctx, c, args.Claims, args.ResourceID, args.ItemID, action)
	if err != nil {
		postTpl.HTMLWithErr(errors.Wrap(err, "failed to start downloading"), http.StatusBadRequest, c, d)
		return
	}
	d.Job = job
	postTpl.HTML(http.StatusOK, c, d)
}
