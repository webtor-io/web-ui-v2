package ext

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/webtor-io/web-ui/services/template"
	"net/http"
	"strconv"
)

type Handler struct {
	tb template.Builder
}

func RegisterHandler(r *gin.Engine, tm *template.Manager) {
	h := &Handler{
		tb: tm.MustRegisterViews("ext/*"),
	}
	r.GET("/ext/download", h.download)
	r.GET("/ext/magnet", h.magnet)
}

type DownloadData struct {
	DownloadID int
}

func (s *Handler) download(c *gin.Context) {
	i, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.WithError(err).Error("failed to parse download id")
		c.AbortWithStatus(http.StatusBadRequest)
	}

	d := DownloadData{
		DownloadID: i,
	}
	s.tb.Build("ext/download").HTML(http.StatusOK, c, d)
}

type MagnetData struct {
	Magnet string
}

func (s *Handler) magnet(c *gin.Context) {
	magnet := c.Query("url")

	d := MagnetData{
		Magnet: magnet,
	}
	s.tb.Build("ext/magnet").HTML(http.StatusOK, c, d)
}
