package job

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Handler) log(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()
	l, ok, err := s.q.GetOrCreate(c.Param("queue_id")).Log(ctx, c.Param("job_id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache,no-store,no-transform")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	c.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Done():
			return false
		case msg, ok := <-l:
			if !ok {
				return false
			}
			c.SSEvent("message", msg)
			return true
		}
	})
}
