package resource

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	sv "github.com/webtor-io/web-ui-v2/services"
	j "github.com/webtor-io/web-ui-v2/services/job"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type Handler struct {
	*w.TemplateHandler
	*w.ClaimsHandler
	api  *sv.Api
	jobs *j.Handler
}

func NewHandler(c *cli.Context, api *sv.Api, jobs *j.Handler) *Handler {
	return &Handler{
		TemplateHandler: w.NewTemplateHandler(c),
		ClaimsHandler:   w.NewClaimsHandler(),
		api:             api,
		jobs:            jobs,
	}
}

func (s *Handler) RegisterRoutes(r *gin.Engine) {
	// r.GET("/", s.getIndex)
	r.POST("/", s.post)
	// r.GET("/queue/:queue_id/job/:job_id/log", s.getJobLog)
	r.GET("/:resource_id", s.get)
}

func (s *Handler) RegisterTemplates(r multitemplate.Renderer) {
	s.RegisterTemplate(r, "resource/get", []string{"standard", "async", "async_list", "async_file"})
}
