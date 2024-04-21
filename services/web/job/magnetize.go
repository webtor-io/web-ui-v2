package job

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/webtor-io/web-ui-v2/services"

	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
)

func (s *Handler) Magnetize(claims *api.Claims, query string) (j *job.Job, err error) {
	sha1 := services.SHA1R.Find([]byte(query))
	if sha1 == nil {
		err = errors.Errorf("wrong resource provided query=%v", query)
		return
	}
	id := strings.ToLower(string(sha1))
	if !strings.HasPrefix(query, "magnet:") {
		query = "magnet:?xt=urn:btih:" + id
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	j = s.q.GetOrCreate("magnetize").Enqueue(ctx, id, func(j *job.Job) {
		j.Info("sadly, we don't have torrent, so we have to magnetize it from peers")
		j.InProgress("magnetizing", "magnetizing")
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		resp, err := s.api.StoreResource(ctx, claims, []byte(query))
		if err != nil || resp == nil {
			j.Error(err, "failed to magnetize, there were no peers for 30 seconds, try another magnet", "magnetizing")
		} else {
			j.Done("magnetizing")
			j.Redirect("/" + resp.ID).Close()
		}
	})
	return
}
