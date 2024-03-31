package job

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ra "github.com/webtor-io/rest-api/services"
	sv "github.com/webtor-io/web-ui-v2/services"
	m "github.com/webtor-io/web-ui-v2/services/models"
)

type StreamContent struct {
	ExportTag           *ra.ExportTag
	MediaProbe          *sv.MediaProbe
	OpenSubtitles       []sv.OpenSubtitleTrack
	VideoStreamUserData *m.VideoStreamUserData
}

func (s *Handler) streamContent(j *sv.Job, c *gin.Context, claims *sv.Claims, resourceID string, itemID string, template string) {
	sc := &StreamContent{}
	j.InProgress("retrieving stream url", "retrieving stream")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	resp, err := s.api.ExportResourceContent(ctx, claims, resourceID, itemID)
	if err != nil {
		j.Error(err, "failed to retrieve for 5 minutes", "retrieving stream").Close()
		return
	}
	j.Done("retrieving stream")
	sc.ExportTag = resp.ExportItems["stream"].Tag
	se := resp.ExportItems["stream"]
	if !se.ExportMetaItem.Meta.Cache {
		if err := s.warmUp(j, "warming up torrent client", resp.ExportItems["download"].URL, resp.ExportItems["torrent_client_stat"].URL, int(resp.Source.Size), 1_000_000, 500_000, "file", true); err != nil {
			return
		}
		if se.Meta.Transcode {
			if err := s.warmUp(j, "warming up transcoder", resp.ExportItems["stream"].URL, resp.ExportItems["torrent_client_stat"].URL, 0, -1, -1, "stream", false); err != nil {
				return
			}
			j.InProgress("probing content media info", "probe media")
			mp, err := s.api.GetMediaProbe(ctx, resp.ExportItems["media_probe"].URL)
			if err != nil {
				j.Error(err, "failed to get probe data", "probe media").Close()
				return
			}
			sc.MediaProbe = mp
			log.Infof("got media probe %+v", mp)
			j.Done("probe media")
		}
	}
	if resp.Source.MediaFormat == ra.Video {
		vsud := m.NewVideoStreamUserData(resourceID, itemID)
		vsud.FetchSessionData(c)
		sc.VideoStreamUserData = vsud
		j.InProgress("loading OpenSubtitles", "opensubtitles")
		subs, err := s.api.GetOpenSubtitles(ctx, resp.ExportItems["subtitles"].URL)
		if err != nil {
			j.Warn(err, "failed to get OpenSubtitles", "opensubtitles")
		} else {
			sc.OpenSubtitles = subs
			j.Done("opensubtitles")
		}
	}
	err = s.renderActionTemplate(j, c, sc, template)
	if err != nil {
		j.Error(err, "failed to render resource", "retrieving stream").Close()
	}
	j.InProgress("waiting player initialization", "player init").Close()
}

func (s *Handler) previewImage(j *sv.Job, c *gin.Context, claims *sv.Claims, resourceID string, itemID string) {
	s.streamContent(j, c, claims, resourceID, itemID, "preview_image")
}

func (s *Handler) streamAudio(j *sv.Job, c *gin.Context, claims *sv.Claims, resourceID string, itemID string) {
	s.streamContent(j, c, claims, resourceID, itemID, "stream_audio")
}

func (s *Handler) streamVideo(j *sv.Job, c *gin.Context, claims *sv.Claims, resourceID string, itemID string) {
	s.streamContent(j, c, claims, resourceID, itemID, "stream_video")
}

func (s *Handler) renderActionTemplate(j *sv.Job, c *gin.Context, sc *StreamContent, name string) error {
	template := "action/" + name
	tpl := s.tm.MakeTemplateWithLayout(template, `{{ template "main" . }}`)
	str, err := tpl.ToString(c, sc)
	if err != nil {
		return err
	}
	j.RenderTemplate(template, strings.TrimSpace(str))
	return nil
}

func (s *Handler) download(j *sv.Job, claims *sv.Claims, resourceID string, itemID string) {
	j.InProgress("retrieving download link", "retrieving download link")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := s.api.ExportResourceContent(ctx, claims, resourceID, itemID)
	if err != nil {
		j.Error(err, "failed to retrieve for 30 seconds", "retrieving download link").Close()
	}
	j.Done("retrieving download link")
	de := resp.ExportItems["download"]
	url := de.URL
	if !de.ExportMetaItem.Meta.Cache {
		if err := s.warmUp(j, "warming up torrent client", resp.ExportItems["download"].URL, resp.ExportItems["torrent_client_stat"].URL, int(resp.Source.Size), 1_000_000, 0, "", true); err != nil {
			return
		}
	}
	j.Download(url).Close()
}

func (s *Handler) warmUp(j *sv.Job, m string, u string, su string, size int, limitStart int, limitEnd int, tagSuff string, useStatus bool) error {
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
		j.InProgress(fmt.Sprintf("%v, downloading %v", m, humanize.Bytes(uint64(limitStart+limitEnd))), tag)
	} else {
		j.InProgress(m, tag)
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel2()

	if useStatus {
		j.StatusUpdate("waiting for peers", tag)
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
	}

	b, err := s.api.DownloadWithRange(ctx2, u, 0, limitStart)
	if err != nil {
		j.Error(err, "failed to start download", tag).Close()
		return err
	}
	defer b.Close()

	_, err = io.Copy(io.Discard, b)

	if limitEnd > 0 {
		b2, err := s.api.DownloadWithRange(ctx2, u, size-limitEnd, -1)
		if err != nil {
			j.Error(err, "failed to start download", tag).Close()
			return err
		}
		defer b2.Close()
		_, err = io.Copy(io.Discard, b2)
	}
	if err != nil {
		j.Error(err, "failed to download within 5 minutes", tag).Close()
		return err
	}

	j.Done(tag)
	return nil
}

func (s *Handler) Action(ctx context.Context, c *gin.Context, claims *sv.Claims, resourceID string, itemID string, action string) (job *sv.Job, err error) {
	id := fmt.Sprintf("%x", sha1.Sum([]byte(resourceID+"/"+itemID+"/"+action+"/"+claims.Role+"/"+claims.SessionID)))
	job = s.q.GetOrCreate(action).Enqueue(ctx, id, func(j *sv.Job) {
		switch action {
		case "download":
			s.download(j, claims, resourceID, itemID)
		case "preview-image":
			s.previewImage(j, c, claims, resourceID, itemID)
		case "stream-audio":
			s.streamAudio(j, c, claims, resourceID, itemID)
		case "stream-video":
			s.streamVideo(j, c, claims, resourceID, itemID)
		}
	})
	return
}
