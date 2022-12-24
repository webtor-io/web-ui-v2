package web

import (
	"fmt"
	"strings"

	ra "github.com/webtor-io/rest-api/services"
	sv "github.com/webtor-io/web-ui-v2/services"

	h "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
)

type Breadcrumb struct {
	Name    string
	PathStr string
	Active  bool
}

type Pagination struct {
	Page   uint
	Active bool
	Prev   bool
	Next   bool
	Number bool
}

// func getAssetHash(name string) (string, error) {
// 	f, err := os.Open(assetsPath + "/" + name)
// 	if err != nil {
// 		return "", err
// 	}
// 	h := md5.New()
// 	if _, err := io.Copy(h, f); err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(h.Sum(nil)), nil
// }

func MakeJobLogURL(j *sv.Job) string {
	return fmt.Sprintf("/queue/%v/job/%v/log", j.Queue, j.ID)
}

func HasBreadcrumbs(lr *ra.ListResponse) bool {
	hasDir := false
	for _, i := range lr.Items {
		if i.Type == ra.ListTypeDirectory {
			hasDir = true
			break
		}
	}
	return hasDir || lr.ListItem.PathStr != "/"
}

func MakeBreadcrumbs(r ra.ResourceResponse, pathStr string) []Breadcrumb {
	res := []Breadcrumb{}
	res = append(res, Breadcrumb{
		Name:    r.Name,
		PathStr: "/",
	})
	if pathStr != "/" {
		t := []string{}
		path := strings.Split(strings.Trim(pathStr, "/"), "/")
		for _, p := range path {
			t = append(t, p)
			res = append(res, Breadcrumb{
				Name:    p,
				PathStr: "/" + strings.Join(t, "/"),
			})
		}
	}
	res[len(res)-1].Active = true
	return res
}

func HasPagination(lr *ra.ListResponse) bool {
	return lr.Count > len(lr.Items)
}

func MakePagination(lr *ra.ListResponse, page uint, pageSize uint) []Pagination {
	res := []Pagination{}
	pages := uint(lr.Count)/pageSize + 1
	prev := page - 1
	if prev < 1 {
		prev = 1
	}
	next := page + 1
	if next > pages {
		next = pages
	}
	res = append(res, Pagination{
		Page:   prev,
		Active: prev != page,
		Prev:   true,
	})
	for i := uint(1); i < pages+1; i++ {
		res = append(res, Pagination{
			Page:   i,
			Active: i != page,
			Number: true,
		})
	}
	res = append(res, Pagination{
		Page:   next,
		Active: next != page,
		Next:   true,
	})
	return res
}

func Log(err error) error {
	log.Error(err)
	return err
}

func ShortErr(err error) string {
	return strings.Split(err.Error(), ":")[0]
}

func BitsForHumans(b int64) string {
	return h.Bytes(uint64(b))
}

type ErrorData struct {
	Err error
}
type CSRFData struct {
	CSRF string
}

type JobData struct {
	Job *sv.Job
}
