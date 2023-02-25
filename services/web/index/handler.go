package index

import (
	csrf "github.com/utrack/gin-csrf"
	"html/template"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type IndexData struct {
	w.ErrorData
	w.CSRFData
	w.JobData
	Args any
}

type Handler struct {
	*w.TemplateHandler
	*w.ClaimsHandler
}

func RegisterHandler(c *cli.Context, r *gin.Engine, re multitemplate.Renderer) {
	h := &Handler{
		TemplateHandler: w.NewTemplateHandler(c, re),
		ClaimsHandler:   w.NewClaimsHandler(),
	}
	r.GET("/", h.index)
	h.RegisterTemplate(
		"index",
		[]string{"standard", "async"},
		[]string{},
		template.FuncMap{},
	)
}

func (s *Handler) index(c *gin.Context) {
	var d IndexData
	d.CSRF = csrf.GetToken(c)
	index := s.MakeTemplate(c, "index", &d)
	index.R(http.StatusOK)
}
