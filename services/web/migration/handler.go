package migration

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(r *gin.Engine) {
	r.GET("/show", func(c *gin.Context) {
		params := c.Request.URL.Query()
		if params.Get("downloadId") != "" {
			c.Redirect(http.StatusMovedPermanently, "/ext/download?id="+params.Get("downloadId"))
			return
		}
		if params.Get("magnet") != "" {
			p := url.Values{}
			p.Add("url", params.Get("magnet"))
			c.Redirect(http.StatusMovedPermanently, "/ext/magnet?url="+p.Encode())
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/embed?"+params.Encode())
	})
}
