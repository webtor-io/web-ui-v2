package resource

import (
	"fmt"
	"strings"

	h "github.com/dustin/go-humanize"
	ra "github.com/webtor-io/rest-api/services"
	w "github.com/webtor-io/web-ui-v2/services/web"
)

type ButtonItem struct {
	ID         string
	CSRF       string
	ItemID     string
	ResourceID string
	Name       string
	Action     string
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

func MakeButton(ctx *w.Context, gd *GetData, name string, icon string, endpoint string) *ButtonItem {
	return &ButtonItem{
		ID:         gd.Item.ID,
		ItemID:     gd.Item.ID,
		ResourceID: gd.Resource.ID,
		Name:       name,
		Action:     icon,
		Endpoint:   endpoint,
		CSRF:       ctx.CSRF,
	}
}

func MakeDirButton(ctx *w.Context, gd *GetData, name string, action string, endpoint string) *ButtonItem {
	return &ButtonItem{
		ID:         gd.List.ID,
		ItemID:     gd.List.ID,
		ResourceID: gd.Resource.ID,
		Name:       name,
		Action:     action,
		Endpoint:   endpoint,
		CSRF:       ctx.CSRF,
	}
}

func MakeFileDownload(ctx *w.Context, gd *GetData) *ButtonItem {
	return MakeButton(ctx, gd,
		fmt.Sprintf("Download [%v]", h.Bytes(uint64(gd.Item.Size))),
		"download",
		"/download-file",
	)
}

func MakeImage(ctx *w.Context, gd *GetData) *ButtonItem {
	return MakeButton(ctx, gd,
		"Preview",
		"preview",
		"/preview-image",
	)
}

func MakeAudio(ctx *w.Context, gd *GetData) *ButtonItem {
	return MakeButton(ctx, gd,
		"Stream",
		"stream",
		"/stream-audio",
	)
}
func MakeVideo(ctx *w.Context, gd *GetData) *ButtonItem {
	return MakeButton(ctx, gd,
		"Stream",
		"stream",
		"/stream-video",
	)
}

func MakeDirDownload(ctx *w.Context, gd *GetData) *ButtonItem {
	return MakeDirButton(ctx, gd,
		fmt.Sprintf("Download Directory as ZIP [%v]", h.Bytes(uint64(gd.List.Size))),
		"download",
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

func (s *Handler) IsDemoMagnet(m string) bool {
	return s.dm == m
}
