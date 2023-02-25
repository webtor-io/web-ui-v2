package resource

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	sv "github.com/webtor-io/web-ui-v2/services"
	j "github.com/webtor-io/web-ui-v2/services/job"
	w "github.com/webtor-io/web-ui-v2/services/web"
	"html/template"
)

type Handler struct {
	*w.TemplateHandler
	*w.ClaimsHandler
	api  *sv.Api
	jobs *j.Handler
}

func RegisterHandler(c *cli.Context, r *gin.Engine, re multitemplate.Renderer, api *sv.Api, jobs *j.Handler) {
	h := &Handler{
		TemplateHandler: w.NewTemplateHandler(c, re),
		ClaimsHandler:   w.NewClaimsHandler(),
		api:             api,
		jobs:            jobs,
	}
	r.POST("/", h.post)
	r.GET("/:resource_id", h.get)
	h.RegisterTemplate(
		"resource/get",
		[]string{"standard", "async", "async_list", "async_file"},
		[]string{"list", "button", "icons", "file"},
		template.FuncMap{
			"makeBreadcrumbs":  MakeBreadcrumbs,
			"hasBreadcrumbs":   HasBreadcrumbs,
			"hasPagination":    HasPagination,
			"makePagination":   MakePagination,
			"makeFileDownload": MakeFileDownload,
			"makeDirDownload":  MakeDirDownload,
			"makeImage":        MakeImage,
			"makeAudio":        MakeAudio,
			"makeVideo":        MakeVideo,
		},
	)
}
