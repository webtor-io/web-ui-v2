package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/template"
	j "github.com/webtor-io/web-ui-v2/services/web/job"
)

type Handler struct {
	api  *api.Api
	jobs *j.Handler
	tm   *template.Manager
}

func RegisterHandler(c *cli.Context, r *gin.Engine, tm *template.Manager, api *api.Api, jobs *j.Handler) {
	h := &Handler{
		api:  api,
		jobs: jobs,
		tm:   tm,
	}
	r.POST("/", h.post)
	r.GET("/:resource_id", h.get)

	tm.RegisterViewsWithFuncs("resource/*", template.FuncMap{
		"makeBreadcrumbs":  MakeBreadcrumbs,
		"hasBreadcrumbs":   HasBreadcrumbs,
		"hasPagination":    HasPagination,
		"makePagination":   MakePagination,
		"makeFileDownload": MakeFileDownload,
		"makeDirDownload":  MakeDirDownload,
		"makeImage":        MakeImage,
		"makeAudio":        MakeAudio,
		"makeVideo":        MakeVideo,
	})
}
