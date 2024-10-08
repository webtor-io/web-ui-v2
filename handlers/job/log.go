package job

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Handler) log(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	l, err := s.q.GetOrCreate(c.Param("queue_id")).Log(ctx, c.Param("job_id"))
	if err != nil {
		cancel()
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache,no-store,no-transform")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	clientGone := c.Writer.CloseNotify()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			cancel()
			return false
		case <-c.Request.Context().Done():
			cancel()
			return false
		case msg, ok := <-l:
			if !ok {
				cancel()
				return false
			}
			c.SSEvent("message", msg)
			return true
		}
	})
}
