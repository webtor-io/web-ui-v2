package job

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
)

func (s *Handler) Download(claims *api.Claims, resourceID string, itemID string) (j *job.Job, err error) {
	id := fmt.Sprintf("%x", sha1.Sum([]byte(resourceID+"/"+itemID)))
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	j = s.q.GetOrCreate("download").Enqueue(ctx, id, func(j *job.Job) {
		j.InProgress("retriving download link", "retriving download link")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		resp, err := s.api.ExportResourceContent(ctx, claims, resourceID, itemID)
		if err != nil {
			j.Error(err, "failed to retrive for 30 seconds", "retriving download link")
		}
		de := resp.ExportItems["download"]
		url := de.URL
		if de.ExportMetaItem.Meta.Cache {
			j.Done("retriving download link")
			j.Download(url)
			return
		}
		j.Info("sadly, we don't have this file in cache, so we have to proceed with warm up before download")
		j.InProgress("downloading ", "downloading")
	})
	return
}
