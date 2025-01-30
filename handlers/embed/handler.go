package embed

import (
	j "github.com/webtor-io/web-ui/handlers/job"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui/services/embed"
	"github.com/webtor-io/web-ui/services/template"
)

type Handler struct {
	tb   template.Builder
	hCl  *http.Client
	jobs *j.Handler
	ds   *embed.DomainSettings
}

func RegisterHandler(hCl *http.Client, r *gin.Engine, tm *template.Manager, jobs *j.Handler, ds *embed.DomainSettings) {
	h := &Handler{
		tb:   tm.MustRegisterViews("embed/*"),
		jobs: jobs,
		ds:   ds,
		hCl:  hCl,
	}
	r.GET("/embed", h.get)
	r.POST("/embed", h.post)
}
