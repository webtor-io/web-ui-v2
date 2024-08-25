package job

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Handler) log(c *gin.Context) {
	l, err := s.q.GetOrCreate(c.Param("queue_id")).Log(c.Request.Context(), c.Param("job_id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-l; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
