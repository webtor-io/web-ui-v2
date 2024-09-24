package job

import (
	"context"
	"github.com/webtor-io/web-ui-v2/handlers/job/script"
	"github.com/webtor-io/web-ui-v2/services/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
)

func (s *Handler) Embed(c *gin.Context, hCl *http.Client, claims *api.Claims, settings *models.EmbedSettings) (j *job.Job, err error) {
	es, hash, err := script.Embed(s.tb, hCl, c, s.api, claims, settings, "")
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	j = s.q.GetOrCreate("embded").Enqueue(ctx, cancel, hash, es)
	return
}
