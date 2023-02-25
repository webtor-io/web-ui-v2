package job

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
	ra "github.com/webtor-io/rest-api/services"
	sv "github.com/webtor-io/web-ui-v2/services"
	"io"
	"strings"
	"time"
)

type StreamContent struct {
	ExportTag  *ra.ExportTag
	MediaProbe *sv.MediaProbe
}

func (s *Handler) streamContent(j *sv.Job, claims *sv.Claims, resourceID string, itemID string, template string) {
	sc := &StreamContent{}
	j.InProgress("retrieving stream url", "retrieving stream")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := s.api.ExportResourceContent(ctx, claims, resourceID, itemID)
	if err != nil {
		log.WithError(err).Error("failed to export resource")
		j.Error("failed to retrieve for 30 seconds", "retrieving stream")
	}
	j.Done("retrieving stream")
	sc.ExportTag = resp.ExportItems["stream"].Tag
	se := resp.ExportItems["stream"]
	if se.ExportMetaItem.Meta.Cache {
		err = s.renderActionTemplate(j, sc, template)
		if err != nil {
			log.WithError(err).Error("failed to render resource")
			j.Error("failed to render resource", "retrieving stream")
		}
		j.InProgress("waiting player initialization", "player init")
		return
	}
	j.Info("sadly, we don't have this file in cache, we have to warm up before proceed")
	if err := s.warmUp(j, resp.ExportItems["download"].URL, resp.ExportItems["torrent_client_stat"].URL, int(resp.Source.Size), 1_000_000, 500_000, "file"); err != nil {
		return
	}
	if se.Meta.Transcode {
		j.Info("content must be transcoded, we have to warm up transcoder")
		if err := s.warmUp(j, resp.ExportItems["stream"].URL, resp.ExportItems["torrent_client_stat"].URL, 0, -1, -1, "stream"); err != nil {
			return
		}
		j.InProgress("probing content media info", "probe media")
		mp, err := s.api.GetMediaProbe(ctx, resp.ExportItems["media_probe"].URL)
		if err != nil {
			j.Error("failed to get probe data", "probe media")
			return
		}
		sc.MediaProbe = mp
		log.Infof("got media probe %+v", mp)
		j.Done("probe media")
	}
	err = s.renderActionTemplate(j, sc, template)
	if err != nil {
		log.WithError(err).Error("failed to render resource")
		j.Error("failed to render resource", "retrieving stream")
	}
	j.InProgress("waiting player initialization", "player init")
}

func (s *Handler) previewImage(j *sv.Job, claims *sv.Claims, resourceID string, itemID string) {
	s.streamContent(j, claims, resourceID, itemID, "preview_image")
}

func (s *Handler) streamAudio(j *sv.Job, claims *sv.Claims, resourceID string, itemID string) {
	s.streamContent(j, claims, resourceID, itemID, "stream_audio")
}

func (s *Handler) streamVideo(j *sv.Job, claims *sv.Claims, resourceID string, itemID string) {
	s.streamContent(j, claims, resourceID, itemID, "stream_video")
}

func (s *Handler) renderActionTemplate(j *sv.Job, sc *StreamContent, name string) error {
	var b bytes.Buffer
	template := "action/" + name + "_async"
	re, _ := s.re.Instance(template, sc).(render.HTML)
	err := re.Template.Execute(&b, re.Data)
	if err != nil {
		return err
	}
	j.RenderTemplate(template, strings.TrimSpace(b.String()))
	return nil
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
	j.Info("sadly, we don't have this file in cache, we have to warm up before proceed")
	if err := s.warmUp(j, resp.ExportItems["download"].URL, resp.ExportItems["torrent_client_stat"].URL, int(resp.Source.Size), 1_000_000, 0, ""); err != nil {
		return
	}
	j.Download(url)
}

func (s *Handler) warmUp(j *sv.Job, u string, su string, size int, limitStart int, limitEnd int, tagSuff string) error {
	tag := "download"
	if tagSuff != "" {
		tag += "-" + tagSuff
	}
	if limitStart > size {
		limitStart = size
	}
	if limitEnd > size-limitStart {
		limitEnd = size - limitStart
	}
	if size > 0 {
		j.InProgress(fmt.Sprintf("try to download %v", humanize.Bytes(uint64(limitStart+limitEnd))), tag)
	} else {
		j.InProgress("try to retrieve file", tag)
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel2()
	j.StatusUpdate("waiting for peers", tag)
	b, err := s.api.DownloadWithRange(ctx2, u, 0, limitStart)
	if err != nil {
		log.WithError(err).Error("failed start download")
		j.Error("failed to start download", tag)
		return err
	}
	defer b.Close()

	go func() {
		ch, err := s.api.Stats(ctx2, su)
		if err != nil {
			log.WithError(err).Error("failed to get stats")
			return
		}
		for ev := range ch {
			j.StatusUpdate(fmt.Sprintf("%v peers", ev.Peers), tag)
		}
	}()

	_, err = io.Copy(io.Discard, b)

	if limitEnd > 0 {
		b2, err := s.api.DownloadWithRange(ctx2, u, size-limitEnd, -1)
		if err != nil {
			log.WithError(err).Error("failed start download")
			j.Error("failed to start download", tag)
			return err
		}
		defer b2.Close()
		_, err = io.Copy(io.Discard, b2)
	}
	if err != nil {
		j.Error("failed to download within 5 minutes", tag)
		log.WithError(err).Error("failed to download bytes within 5 minutes")
		return err
	}

	j.Done(tag)
	return nil
}

func (s *Handler) Action(claims *sv.Claims, resourceID string, itemID string, action string) (job *sv.Job, err error) {
	id := fmt.Sprintf("%x", sha1.Sum([]byte(resourceID+"/"+itemID+"/"+action+"/"+claims.Role+"/"+claims.SessionID)))
	job = s.q.GetOrCreate(action).Enqueue(id, func(j *sv.Job) {
		switch action {
		case "download":
			s.download(j, claims, resourceID, itemID)
			break
		case "preview-image":
			s.previewImage(j, claims, resourceID, itemID)
			break
		case "stream-audio":
			s.streamAudio(j, claims, resourceID, itemID)
			break
		case "stream-video":
			s.streamVideo(j, claims, resourceID, itemID)
			break
		}
	})
	return
}
