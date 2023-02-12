package resource

import (
	ra "github.com/webtor-io/rest-api/services"
	"strings"
)

type DownloadItem struct {
	ID         string
	CSRF       string
	Size       int64
	ItemID     string
	ResourceID string
	Name       string
}

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

func MakeFileDownload(gd *GetData) *DownloadItem {
	return &DownloadItem{
		ID:         gd.Item.ID,
		CSRF:       gd.CSRF,
		Size:       gd.Item.Size,
		ItemID:     gd.Item.ID,
		ResourceID: gd.Resource.ID,
		Name:       "Download",
	}
}

func MakeDirDownload(gd *GetData) *DownloadItem {
	return &DownloadItem{
		ID:         gd.List.ID,
		CSRF:       gd.CSRF,
		Size:       gd.List.Size,
		ItemID:     gd.List.ID,
		ResourceID: gd.Resource.ID,
		Name:       "Download Directory as ZIP",
	}
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
