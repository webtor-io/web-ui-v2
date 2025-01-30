package job

import (
	"context"
	"github.com/webtor-io/web-ui/handlers/job/script"
	"github.com/webtor-io/web-ui/services/claims"
	"github.com/webtor-io/web-ui/services/models"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/webtor-io/web-ui/services/api"
	"github.com/webtor-io/web-ui/services/job"
)

func (s *Handler) Action(c *gin.Context, apiClaims *api.Claims, userClaims *claims.Data, resourceID string, itemID string, action string, settings *models.StreamSettings, purge bool) (j *job.Job, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	as, id := script.Action(s.tb, s.api, apiClaims, userClaims, c, resourceID, itemID, action, settings)
	j = s.q.GetOrCreate(action).Enqueue(ctx, cancel, id, as, purge)
	return
}
