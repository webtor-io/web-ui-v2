package job

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	sv "github.com/webtor-io/web-ui-v2/services"
)

func (s *Handler) Download(claims *sv.Claims, resourceID string, itemID string) (job *sv.Job, err error) {
	id := fmt.Sprintf("%x", sha1.Sum([]byte(resourceID+"/"+itemID)))
	job = s.q.GetOrCreate("download").Enqueue(id, func(j *sv.Job) {
		j.InProgress("retrieving download link", "retrieving download link")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		resp, err := s.api.ExportResourceContent(ctx, claims, resourceID, itemID)
		if err != nil {
			log.WithError(err).Error("failed to export resource")
			j.Error("failed to retrieve for 30 seconds", "retrieving download link")
		}
		j.Done("retrieving download link")
		de := resp.ExportItems["download"]
		url := de.URL
		if de.ExportMetaItem.Meta.Cache {
			j.Download(url)
			return
		}
		j.Info("sadly, we don't have this file in cache, we have to warm up before download")
		size := uint64(resp.Source.Size)
		var limit uint64 = 1_000_000
		if limit > size {
			limit = size
		}
		j.InProgress(fmt.Sprintf("try to download %v", humanize.Bytes(limit)), "download")
		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel2()
		j.StatusUpdate("waiting for peers", "download")
		b, err := s.api.Download(ctx2, resp)
		if err != nil {
			log.WithError(err).Error("failed start download")
			j.Error("failed to start download", "download")
			return
		}
		defer b.Close()

		go func() {
			ch, err := s.api.Stats(ctx2, resp)
			if err != nil {
				log.WithError(err).Error("failed to get stats")
				return
			}
			for ev := range ch {
				j.StatusUpdate(fmt.Sprintf("%v peers", ev.Peers), "download")
			}
		}()

		_, err = io.CopyN(io.Discard, b, int64(limit))
		if err != nil {
			j.Error("failed to download within 5 minutes", "download")
			log.WithError(err).Error("failed to download bytes within 5 minutes")
			return
		}
		j.Done("download")
		j.Download(url)
	})
	return
}
