package embed

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
	"github.com/webtor-io/web-ui-v2/services/web/job/script"
)

type EmbedSettings struct {
	BaseURL   string   `json:"baseUrl"`
	Width     string   `json:"width"`
	Height    string   `json:"height"`
	Mode      string   `json:"mode"`
	Subtitles []string `json:"subtitles"`
	Poster    string   `json:"poster"`
	Header    bool     `json:"header"`
	Title     string   `json:"title"`
	ImdbID    string   `json:"imdbId"`
	Version   string   `json:"version"`
	Lang      string   `json:"lang"`
	I18n      struct{} `json:"i18n"`
	Features  struct{} `json:"features"`
	El        struct{} `json:"el"`
	Magnet    string   `json:"magnet"`
	Controls  bool     `json:"controls"`
}
type PostArgs struct {
	ID       string
	Settings *EmbedSettings
	Claims   *api.Claims
}

type PostData struct {
	ID       string
	Settings *EmbedSettings
	Job      *job.Job
}

func (s *Handler) bindPostArgs(c *gin.Context) (*PostArgs, error) {
	rawSettings := c.PostForm("settings")
	var settings EmbedSettings
	err := json.Unmarshal([]byte(rawSettings), &settings)
	if err != nil {
		return nil, err
	}
	id := c.Query("id")

	return &PostArgs{
		ID:       id,
		Settings: &settings,
		Claims:   api.GetClaimsFromContext(c),
	}, nil

}

func (s *Handler) post(c *gin.Context) {
	tpl := s.tb.Build("embed/post")
	pd := PostData{}
	args, err := s.bindPostArgs(c)
	if err != nil {
		tpl.HTMLWithErr(err, http.StatusBadRequest, c, pd)
		return
	}
	pd.ID = args.ID
	pd.Settings = args.Settings
	job, err := s.jobs.Embed(c, args.Claims, &script.LoadArgs{
		Query: args.Settings.Magnet,
	})
	if err != nil {
		tpl.HTMLWithErr(err, http.StatusBadRequest, c, pd)
		return
	}
	pd.Job = job
	tpl.HTML(http.StatusAccepted, c, pd)
}
