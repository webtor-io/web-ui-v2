package job

import (
	"github.com/gin-gonic/gin"
	sv "github.com/webtor-io/web-ui-v2/services"
)

type Handler struct {
	q *sv.JobQueues
}

func RegisterHandler(r *gin.Engine, queues *sv.JobQueues) {
	h := &Handler{
		q: queues,
	}
	r.GET("/queue/:queue_id/job/:job_id/log", h.log)
}
