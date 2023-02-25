package resource

import (
	"fmt"
	h "github.com/dustin/go-humanize"
	ra "github.com/webtor-io/rest-api/services"
	"strings"
)

type ButtonItem struct {
	ID         string
	CSRF       string
	ItemID     string
	ResourceID string
	Name       string
	Icon       string
	Endpoint   string
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

func MakeButton(gd *GetData, name string, icon string, endpoint string) *ButtonItem {
	return &ButtonItem{
		ID:         gd.Item.ID,
		CSRF:       gd.CSRF,
		ItemID:     gd.Item.ID,
		ResourceID: gd.Resource.ID,
		Name:       name,
		Icon:       icon,
		Endpoint:   endpoint,
	}
}

func MakeFileDownload(gd *GetData) *ButtonItem {
	return MakeButton(gd,
		fmt.Sprintf("Download [%v]", h.Bytes(uint64(gd.Item.Size))),
		"action",
		"/download-file",
	)
}

func MakeImage(gd *GetData) *ButtonItem {
	return MakeButton(gd,
		"Preview",
		"preview",
		"/preview-image",
	)
}

func MakeAudio(gd *GetData) *ButtonItem {
	return MakeButton(gd,
		"Stream",
		"stream",
		"/stream-audio",
	)
}
func MakeVideo(gd *GetData) *ButtonItem {
	return MakeButton(gd,
		"Stream",
		"stream",
		"/stream-video",
	)
}

func MakeDirDownload(gd *GetData) *ButtonItem {
	return MakeButton(gd,
		fmt.Sprintf("Download Directory as ZIP [%v]", h.Bytes(uint64(gd.List.Size))),
		"action",
		"/download-dir",
	)
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
	var res []Breadcrumb
	res = append(res, Breadcrumb{
		Name:    r.Name,
		PathStr: "/",
	})
	if pathStr != "/" {
		var t []string
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
	var res []Pagination
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
