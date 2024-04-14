package resource

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	sv "github.com/webtor-io/web-ui-v2/services"
	sw "github.com/webtor-io/web-ui-v2/services/web"

	ra "github.com/webtor-io/rest-api/services"
)

type PostArgs struct {
	File   *multipart.FileHeader
	Query  string
	ID     string
	Claims *sv.Claims
}

func (s *Handler) bindPostArgs(c *gin.Context) (*PostArgs, error) {
	file, _ := c.FormFile("resource")
	r, ok := c.GetPostFormArray("resource")
	query := ""
	id := ""
	if ok {
		query = r[0]
		sha1 := sv.SHA1R.Find([]byte(query))
		if sha1 == nil {
			return &PostArgs{Query: query}, errors.Errorf("wrong resource provided query=%v", query)
		}
		id = strings.ToLower(string(sha1))
	}

	if file == nil && query == "" {
		return nil, errors.Errorf("no resource provided")
	}
	claims, err := s.MakeClaims(c)
	if err != nil {
		return nil, err
	}
	return &PostArgs{
		File:   file,
		Query:  query,
		Claims: claims,
		ID:     id,
	}, nil
}

func (s *Handler) postFile(ctx context.Context, args *PostArgs) (res *ra.ResourceResponse, err error) {
	f, err := args.File.Open()
	if err != nil {
		return
	}
	defer f.Close()
	in, err := io.ReadAll(f)
	if err != nil {
		return
	}
	res, err = s.api.StoreResource(ctx, args.Claims, []byte(in))
	return
}

func (s *Handler) postQuery(ctx context.Context, args *PostArgs) (res *ra.ResourceResponse, err error) {
	res, err = s.api.GetResource(ctx, args.Claims, args.ID)
	if err != nil {
		return
	}
	if res != nil {
		return
	}
	return
}

type PostData struct {
	sw.JobData
	Args *PostArgs
}

func (s *Handler) post(c *gin.Context) {
	indexTpl := s.tm.MakeTemplate("index")
	var (
		d    PostData
		err  error
		args *PostArgs
		job  *sv.Job
		res  *ra.ResourceResponse
	)
	args, err = s.bindPostArgs(c)
	d.Args = args
	if err != nil {
		indexTpl.HTMLWithErr(errors.Wrap(err, "wrong args provided"), http.StatusBadRequest, c, d)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	if args.File != nil {
		res, err = s.postFile(ctx, args)
		if err != nil {
			indexTpl.HTMLWithErr(errors.Wrap(err, "failed to upload file"), http.StatusInternalServerError, c, d)
			return
		}
	}
	if res == nil && args.Query != "" {
		res, err = s.postQuery(ctx, args)
		if err != nil {
			indexTpl.HTMLWithErr(errors.Wrap(err, "failed to process query"), http.StatusInternalServerError, c, d)
			return
		}
	}
	if res != nil {
		c.Redirect(http.StatusFound, "/"+res.ID)
		return
	}
	if res == nil && args.Query != "" {
		job, err = s.jobs.Magnetize(args.Claims, args.Query)
		if err != nil {
			indexTpl.HTMLWithErr(errors.Wrap(err, "failed to start magnetizing"), http.StatusInternalServerError, c, d)
			return
		}
		d.Job = job
		indexTpl.HTML(http.StatusAccepted, c, d)
		return
	}
}
