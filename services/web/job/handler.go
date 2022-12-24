package job

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	sv "github.com/webtor-io/web-ui-v2/services"
)

type Handler struct {
	q *sv.JobQueues
}

func NewHandler(queues *sv.JobQueues) *Handler {
	return &Handler{
		q: queues,
	}
}

func (s *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/queue/:queue_id/job/:job_id/log", s.log)
}

func (s *Handler) RegisterTemplates(r multitemplate.Renderer) {
}
