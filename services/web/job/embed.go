package job

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
	"github.com/webtor-io/web-ui-v2/services/web/job/script"
)

func (s *Handler) Embed(c *gin.Context, claims *api.Claims, args *script.LoadArgs) (j *job.Job, err error) {
	es, hash, err := script.Embed(s.tb, c, s.api, claims, args, "")
	if err != nil {
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	j = s.q.GetOrCreate("embded").Enqueue(ctx, hash, es)
	return
}
