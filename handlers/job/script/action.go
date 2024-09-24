package script

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"io"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	ra "github.com/webtor-io/rest-api/services"
	m "github.com/webtor-io/web-ui-v2/services/models"
	"github.com/webtor-io/web-ui-v2/services/template"

	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
)

type StreamContent struct {
	ExportTag           *ra.ExportTag
	MediaProbe          *api.MediaProbe
	OpenSubtitles       []api.OpenSubtitleTrack
	VideoStreamUserData *m.VideoStreamUserData
	Settings            *m.StreamSettings
	ExternalData        *m.ExternalData
}

type TorrentDownload struct {
	Data     []byte
	Infohash string
	Name     string
	Size     int
}

func (s *ActionScript) streamContent(j *job.Job, c *gin.Context, claims *api.Claims, resourceID string, itemID string, template string, settings *m.StreamSettings, vsud *m.VideoStreamUserData) (err error) {
	sc := &StreamContent{
		Settings:     settings,
		ExternalData: &m.ExternalData{},
	}
	j.InProgress("retrieving stream url")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	resp, err := s.api.ExportResourceContent(ctx, claims, resourceID, itemID, settings.ImdbID)
	if err != nil {
		return j.Error(err, "failed to retrieve for 5 minutes")
	}
	j.Done()
	sc.ExportTag = resp.ExportItems["stream"].Tag
	se := resp.ExportItems["stream"]
	if !se.ExportMetaItem.Meta.Cache {
		if err = s.warmUp(j, "warming up torrent client", resp.ExportItems["download"].URL, resp.ExportItems["torrent_client_stat"].URL, int(resp.Source.Size), 1_000_000, 500_000, "file", true); err != nil {
			return
		}
		if se.Meta.Transcode {
			if err = s.warmUp(j, "warming up transcoder", resp.ExportItems["stream"].URL, resp.ExportItems["torrent_client_stat"].URL, 0, -1, -1, "stream", false); err != nil {
				return
			}
			j.InProgress("probing content media info")
			mp, err := s.api.GetMediaProbe(ctx, resp.ExportItems["media_probe"].URL)
			if err != nil {
				return j.Error(err, "failed to get probe data")
			}
			sc.MediaProbe = mp
			log.Infof("got media probe %+v", mp)
			j.Done()
		}
	}
	if resp.Source.MediaFormat == ra.Video {
		sc.VideoStreamUserData = vsud
		j.InProgress("loading OpenSubtitles")
		subs, err := s.api.GetOpenSubtitles(ctx, resp.ExportItems["subtitles"].URL)
		if err != nil {
			j.Warn(err, "failed to get OpenSubtitles")
		} else {
			sc.OpenSubtitles = subs
			j.Done()
		}
	}
	if settings.Poster != "" {
		sc.ExternalData.Poster = s.api.AttachExternalFile(se, settings.Poster)
	}
	for _, v := range settings.Subtitles {
		sc.ExternalData.Tracks = append(sc.ExternalData.Tracks, m.ExternalTrack{
			Src:     s.api.AttachExternalSubtitle(se, v.Src),
			Label:   v.Label,
			SrcLang: v.SrcLang,
			Default: v.Default != nil,
		})
	}
	err = s.renderActionTemplate(j, c, sc, template)
	if err != nil {
		return j.Error(err, "failed to render resource")
	}
	j.InProgress("waiting player initialization")
	return
}

func (s *ActionScript) previewImage(j *job.Job, c *gin.Context, claims *api.Claims, resourceID string, itemID string, settings *m.StreamSettings, vsud *m.VideoStreamUserData) error {
	return s.streamContent(j, c, claims, resourceID, itemID, "preview_image", settings, vsud)
}

func (s *ActionScript) streamAudio(j *job.Job, c *gin.Context, claims *api.Claims, resourceID string, itemID string, settings *m.StreamSettings, vsud *m.VideoStreamUserData) error {
	return s.streamContent(j, c, claims, resourceID, itemID, "stream_audio", settings, vsud)
}

func (s *ActionScript) streamVideo(j *job.Job, c *gin.Context, claims *api.Claims, resourceID string, itemID string, settings *m.StreamSettings, vsud *m.VideoStreamUserData) error {
	return s.streamContent(j, c, claims, resourceID, itemID, "stream_video", settings, vsud)
}

func (s *ActionScript) renderActionTemplate(j *job.Job, c *gin.Context, sc *StreamContent, name string) error {
	actionTemplate := "action/" + name
	tpl := s.tb.Build(actionTemplate).WithLayoutBody(`{{ template "main" . }}`)
	str, err := tpl.ToString(c, sc)
	if err != nil {
		return err
	}
	j.RenderTemplate(actionTemplate, strings.TrimSpace(str))
	return nil
}

func (s *ActionScript) download(j *job.Job, claims *api.Claims, resourceID string, itemID string) (err error) {
	j.InProgress("retrieving download link")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := s.api.ExportResourceContent(ctx, claims, resourceID, itemID, "")
	if err != nil {
		return j.Error(err, "failed to retrieve for 30 seconds")
	}
	j.Done()
	de := resp.ExportItems["download"]
	url := de.URL
	if !de.ExportMetaItem.Meta.Cache {
		if err := s.warmUp(j, "warming up torrent client", resp.ExportItems["download"].URL, resp.ExportItems["torrent_client_stat"].URL, int(resp.Source.Size), 1_000_000, 0, "", true); err != nil {
			return err
		}
	}
	j.Download(url)
	return
}

func (s *ActionScript) downloadTorrent(j *job.Job, c *gin.Context, claims *api.Claims, resourceID string) (err error) {
	j.InProgress("retrieving torrent")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := s.api.GetTorrent(ctx, claims, resourceID)
	if err != nil {
		return j.Error(err, "failed to retrieve for 30 seconds")
	}
	defer func(resp io.ReadCloser) {
		_ = resp.Close()
	}(resp)
	torrent, err := io.ReadAll(resp)
	if err != nil {
		return j.Error(err, "failed to read torrent")
	}
	mi, err := metainfo.Load(bytes.NewBuffer(torrent))
	if err != nil {
		return j.Error(err, "failed to load torrent metainfo")
	}
	info, err := mi.UnmarshalInfo()
	if err != nil {
		return j.Error(err, "failed to unmarshal torrent metainfo")
	}
	j.DoneWithMessage("success! download should start right now!")
	tpl := s.tb.Build("action/download_torrent").WithLayoutBody(`{{ template "main" . }}`)
	name := info.Name
	if name == "" {
		name = resourceID
	}
	str, err := tpl.ToString(c, &TorrentDownload{
		Data:     torrent,
		Infohash: resourceID,
		Name:     name + ".torrent",
		Size:     len(torrent),
	})
	if err != nil {
		return err
	}
	j.RenderTemplate("action/download_torrent", strings.TrimSpace(str))
	return nil
}

func (s *ActionScript) warmUp(j *job.Job, m string, u string, su string, size int, limitStart int, limitEnd int, tagSuff string, useStatus bool) (err error) {
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
		j.InProgress(fmt.Sprintf("%v, downloading %v", m, humanize.Bytes(uint64(limitStart+limitEnd))))
	} else {
		j.InProgress(m)
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel2()

	if useStatus {
		j.StatusUpdate("waiting for peers")
		go func() {
			ch, err := s.api.Stats(ctx2, su)
			if err != nil {
				log.WithError(err).Error("failed to get stats")
				return
			}
			for ev := range ch {
				j.StatusUpdate(fmt.Sprintf("%v peers", ev.Peers))
			}
		}()
	}

	b, err := s.api.DownloadWithRange(ctx2, u, 0, limitStart)
	if err != nil {
		return j.Error(err, "failed to start download")
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(b)

	_, err = io.Copy(io.Discard, b)

	if limitEnd > 0 {
		b2, err := s.api.DownloadWithRange(ctx2, u, size-limitEnd, -1)
		if err != nil {
			return j.Error(err, "failed to start download")
		}
		defer func(b2 io.ReadCloser) {
			_ = b2.Close()
		}(b2)
		_, err = io.Copy(io.Discard, b2)
	}
	if err != nil {
		return j.Error(err, "failed to download within 5 minutes")
	}

	j.Done()
	return
}

type ActionScript struct {
	api        *api.Api
	claims     *api.Claims
	c          *gin.Context
	resourceId string
	itemId     string
	action     string
	tb         template.Builder
	settings   *m.StreamSettings
	vsud       *m.VideoStreamUserData
}

func (s *ActionScript) Run(j *job.Job) (err error) {
	switch s.action {
	case "download":
		return s.download(j, s.claims, s.resourceId, s.itemId)
	case "download-torrent":
		return s.downloadTorrent(j, s.c, s.claims, s.resourceId)
	case "preview-image":
		return s.previewImage(j, s.c, s.claims, s.resourceId, s.itemId, s.settings, s.vsud)
	case "stream-audio":
		return s.streamAudio(j, s.c, s.claims, s.resourceId, s.itemId, s.settings, s.vsud)
	case "stream-video":
		return s.streamVideo(j, s.c, s.claims, s.resourceId, s.itemId, s.settings, s.vsud)
	}
	return
}

func Action(tb template.Builder, api *api.Api, claims *api.Claims, c *gin.Context, resourceID string, itemID string, action string, settings *m.StreamSettings) (r job.Runnable, id string) {
	vsud := m.NewVideoStreamUserData(resourceID, itemID, settings)
	vsud.FetchSessionData(c)
	vsudID := vsud.AudioID + "/" + vsud.SubtitleID + "/" + fmt.Sprintf("%+v", vsud.AcceptLangTags)
	settingsID := fmt.Sprintf("%+v", settings)
	id = fmt.Sprintf("%x", sha1.Sum([]byte(resourceID+"/"+itemID+"/"+action+"/"+claims.Role+"/"+settingsID+"/"+vsudID)))
	return &ActionScript{
		tb:         tb,
		api:        api,
		claims:     claims,
		c:          c,
		resourceId: resourceID,
		itemId:     itemID,
		action:     action,
		settings:   settings,
		vsud:       vsud,
	}, id
}
