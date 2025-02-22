package job

import (
	"context"
	"github.com/webtor-io/web-ui/handlers/job/script"
	"github.com/webtor-io/web-ui/services/models"
	"github.com/webtor-io/web-ui/services/web"
	"time"

	"github.com/webtor-io/web-ui/services/job"

	m "github.com/webtor-io/web-ui/services/models"
)

func (s *Handler) Action(c *web.Context, resourceID string, itemID string, action string, settings *models.StreamSettings, purge bool, vsud *m.VideoStreamUserData) (j *job.Job, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	as, id := script.Action(s.tb, s.api, c, resourceID, itemID, action, settings, nil, vsud)
	j = s.q.GetOrCreate(action).Enqueue(ctx, cancel, id, as, purge)
	return
}
