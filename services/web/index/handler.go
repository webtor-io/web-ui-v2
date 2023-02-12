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

func NewHandler(c *cli.Context) *Handler {
	return &Handler{
		TemplateHandler: w.NewTemplateHandler(c),
		ClaimsHandler:   w.NewClaimsHandler(),
	}
}

func (s *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", s.index)
}

func (s *Handler) RegisterTemplates(r multitemplate.Renderer) {
	s.RegisterTemplate(
		r,
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
