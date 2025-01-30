package job

import (
	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui/services/api"
	"github.com/webtor-io/web-ui/services/job"
	"github.com/webtor-io/web-ui/services/template"
)

type Handler struct {
	q   *job.Queues
	tb  template.Builder
	api *api.Api
}

func New(q *job.Queues, tm *template.Manager, api *api.Api) *Handler {
	return &Handler{
		q:   q,
		tb:  tm,
		api: api,
	}
}

func (s *Handler) RegisterHandler(r *gin.Engine) *Handler {
	r.GET("/queue/:queue_id/job/:job_id/log", s.log)
	return s
}
