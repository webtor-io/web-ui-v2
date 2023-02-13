package job

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	ra "github.com/webtor-io/rest-api/services"
	sv "github.com/webtor-io/web-ui-v2/services"
	"io"
	"time"
)

func (s *Handler) previewImage(j *sv.Job, claims *sv.Claims, resourceID string, itemID string) {
	j.InProgress("retrieving preview", "retrieving preview")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := s.api.ExportResourceContent(ctx, claims, resourceID, itemID)
	if err != nil {
		log.WithError(err).Error("failed to export resource")
		j.Error("failed to retrieve for 30 seconds", "retrieving preview")
	}
	j.Done("retrieving preview")
	se := resp.ExportItems["stream"]
	if se.ExportMetaItem.Meta.Cache {
		<-time.After(300 * time.Millisecond)
		j.RenderTag(se.ExportStreamItem.Tag)
		return
	}
	if err := s.warmUp(j, resp); err != nil {
		return
	}
	j.RenderTag(se.ExportStreamItem.Tag)
}

func (s *Handler) download(j *sv.Job, claims *sv.Claims, resourceID string, itemID string) {
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
	if err := s.warmUp(j, resp); err != nil {
		return
	}
	j.Download(url)
}

func (s *Handler) warmUp(j *sv.Job, resp *ra.ExportResponse) error {
	j.Info("sadly, we don't have this file in cache, we have to warm up before proceed")
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
		return err
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
		return err
	}
	j.Done("download")
	return nil
}

func (s *Handler) Action(claims *sv.Claims, resourceID string, itemID string, action string) (job *sv.Job, err error) {
	id := fmt.Sprintf("%x", sha1.Sum([]byte(resourceID+"/"+itemID+"/"+action)))
	job = s.q.GetOrCreate(action).Enqueue(id, func(j *sv.Job) {
		switch action {
		case "download":
			s.download(j, claims, resourceID, itemID)
			break
		case "preview-image":
			s.previewImage(j, claims, resourceID, itemID)
			break
		}
	})
	return
}
