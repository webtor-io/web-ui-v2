package job

import (
	"context"
	"time"

	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
	"github.com/webtor-io/web-ui-v2/services/web/job/script"
)

func (s *Handler) Load(claims *api.Claims, args *script.LoadArgs) (j *job.Job, err error) {
	ls, hash, err := script.Load(s.api, claims, args)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	j = s.q.GetOrCreate("load").Enqueue(ctx, cancel, hash, job.NewScript(func(j *job.Job) (err error) {
		err = ls.Run(j)
		if err != nil {
			return
		}
		j.Redirect("/" + j.Context.Value("respID").(string))
		return
	}))
	return
}
