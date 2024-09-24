package resource

import (
	"github.com/pkg/errors"
	"github.com/webtor-io/web-ui-v2/handlers/job/script"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	sv "github.com/webtor-io/web-ui-v2/services"
	"github.com/webtor-io/web-ui-v2/services/api"
	"github.com/webtor-io/web-ui-v2/services/job"
)

type PostArgs struct {
	File        []byte
	Query       string
	Instruction string
	Claims      *api.Claims
}

func (s *Handler) bindPostArgs(c *gin.Context) (*PostArgs, error) {
	file, _ := c.FormFile("resource")
	instruction, _ := c.GetPostForm("instruction")
	r, ok := c.GetPostFormArray("resource")
	query := ""
	if ok {
		query = r[0]
		sha1 := sv.SHA1R.Find([]byte(query))
		if sha1 == nil {
			return &PostArgs{Query: query}, errors.Errorf("wrong resource provided query=%v", query)
		}
	}

	if file == nil && query == "" {
		return nil, errors.Errorf("no resource provided")
	}

	var fd []byte

	if file != nil {
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer func(f multipart.File) {
			_ = f.Close()
		}(f)
		fd, err = io.ReadAll(f)
		if err != nil {
			return nil, err
		}
	}

	return &PostArgs{
		File:        fd,
		Query:       query,
		Claims:      api.GetClaimsFromContext(c),
		Instruction: instruction,
	}, nil
}

type PostData struct {
	Job         *job.Job
	Args        *PostArgs
	Instruction string
}

func (s *Handler) post(c *gin.Context) {
	indexTpl := s.tb.Build("index")
	var (
		d       PostData
		err     error
		args    *PostArgs
		loadJob *job.Job
	)
	args, err = s.bindPostArgs(c)
	if err != nil {
		indexTpl.HTMLWithErr(errors.Wrap(err, "wrong args provided"), http.StatusBadRequest, c, d)
	}
	d.Args = args
	d.Instruction = args.Instruction
	if err != nil {
		indexTpl.HTMLWithErr(errors.Wrap(err, "wrong args provided"), http.StatusBadRequest, c, d)
		return
	}
	loadJob, err = s.jobs.Load(args.Claims, &script.LoadArgs{
		Query: args.Query,
		File:  args.File,
	})
	if err != nil {
		indexTpl.HTMLWithErr(errors.Wrap(err, "failed to load resource"), http.StatusInternalServerError, c, d)
		return
	}
	d.Job = loadJob
	indexTpl.HTML(http.StatusAccepted, c, d)
}
