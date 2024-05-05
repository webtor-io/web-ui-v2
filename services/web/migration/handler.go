package migration

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(r *gin.Engine) {
	r.GET("/show", func(c *gin.Context) {
		params := c.Request.URL.Query()
		c.Redirect(http.StatusMovedPermanently, "/embed?"+params.Encode())
	})
}
