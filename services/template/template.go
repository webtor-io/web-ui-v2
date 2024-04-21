package template

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	"github.com/gin-contrib/multitemplate"
)

type FuncMap = template.FuncMap

type View struct {
	Name         string
	Path         string
	LayoutPath   string
	Layout       string
	LayoutBody   string
	Partials     []string
	Funcs        FuncMap
	once         sync.Once
	err          error
	re           multitemplate.Renderer
	templateName string
}

func (s *View) makeTemplate() (t *template.Template, err error) {
	templates := []string{}
	if s.LayoutBody == "" {
		templates = append(templates, s.LayoutPath)
		t = template.New(filepath.Base(s.LayoutPath))
	} else {
		t, err = template.New(s.Name).Parse(s.LayoutBody)
		if err != nil {
			return nil, err
		}
	}
	templates = append(templates, s.Path)
	templates = append(templates, s.Partials...)
	return t.Funcs(s.Funcs).ParseFiles(templates...)
}

func (s *View) makeTemplateName() string {
	name := s.Name
	if s.LayoutBody != "" {
		hash := md5.Sum([]byte(s.LayoutBody))
		hashStr := hex.EncodeToString(hash[:])
		name += "_" + hashStr

	} else if s.Layout != "main" {
		name += "_" + s.Layout
	}
	return name
}

func (s *View) Render() (string, error) {
	s.once.Do(func() {
		t, err := s.makeTemplate()
		if err != nil {
			return
		}
		s.templateName = s.makeTemplateName()
		s.re.Add(s.templateName, t)
	})
	return s.templateName, s.err
}

type Manager struct {
	re             multitemplate.Renderer
	funcs          FuncMap
	contextWrapper func(c *gin.Context, obj any, err error) any
	layouts        []string
	partials       []string
	views          []*View
	mux            sync.Mutex
	args           Args
}

type Args interface {
	GetFuncs() FuncMap
	WrapContext(c *gin.Context, obj any, err error) any
}

func NewManager(re multitemplate.Renderer, args Args) *Manager {
	return &Manager{
		re:             re,
		args:           args,
		funcs:          args.GetFuncs(),
		contextWrapper: args.WrapContext,
	}
}

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (s *Manager) RegisterViews(pattern string) error {
	return s.RegisterViewsWithFuncs(pattern, FuncMap{})
}

func (s *Manager) GetLayouts() ([]string, error) {
	if s.layouts != nil {
		return s.layouts, nil
	}
	layouts, err := filepath.Glob("templates/layouts/*")
	if err != nil {
		return nil, err
	}
	s.layouts = layouts
	return layouts, nil
}

func (s *Manager) GetPartials() ([]string, error) {
	if s.partials != nil {
		return s.partials, nil
	}
	partials, err := filepath.Glob("templates/partials/*")
	if err != nil {
		return nil, err
	}
	s.partials = partials
	return partials, nil
}

func (s *Manager) makeView(view string, layout string, partials []string, funcs FuncMap) *View {
	lName := fileNameWithoutExt(filepath.Base(layout))
	vName := fileNameWithoutExt(strings.TrimPrefix(view, "templates/views/"))
	return &View{
		re:         s.re,
		Name:       vName,
		Path:       view,
		LayoutPath: layout,
		Layout:     lName,
		Partials:   partials,
		Funcs:      funcs,
	}
}

func (s *Manager) makeViewWithCustomLayout(mv *View, layout string) *View {
	cv := s.makeView(mv.Path, mv.LayoutPath, mv.Partials, mv.Funcs)
	cv.LayoutBody = layout
	return cv
}

func (s *Manager) RegisterViewsWithFuncs(pattern string, f FuncMap) error {
	for k, v := range f {
		s.funcs[k] = v
	}
	views, err := filepath.Glob("templates/views/" + pattern)
	if err != nil {
		return err
	}
	layouts, err := s.GetLayouts()
	if err != nil {
		return err
	}
	partials, err := s.GetPartials()
	if err != nil {
		return err
	}
	for _, v := range views {
		for _, l := range layouts {
			f, _ := os.Stat(v)
			if !f.IsDir() {
				s.views = append(s.views, s.makeView(v, l, partials, s.funcs))
			}
		}
	}

	return nil
}

func (s *Manager) Init() error {
	for _, v := range s.views {
		_, err := s.renderView(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Manager) RenderViewByName(name string) (string, error) {
	for _, v := range s.views {
		if v.Name == name {
			return s.renderView(v)
		}
	}
	return "", errors.New("view not found")
}

func (s *Manager) RenderViewByNameWithCustomLayout(name string, layout string) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	var cv, mv *View
	for _, v := range s.views {
		if v.Name == name && v.LayoutBody == layout {
			cv = v
		}
		if v.Name == name && v.LayoutBody == "" {
			mv = v
		}
	}
	if cv != nil {
		return s.renderView(cv)
	}
	if mv != nil {
		cv = s.makeViewWithCustomLayout(mv, layout)
		s.views = append(s.views, cv)
		return s.renderView(cv)
	}
	return "", errors.New("view not found")
}

func (s *Manager) renderView(v *View) (string, error) {
	name, err := v.Render()
	if err != nil {
		return "", err
	}
	return name, nil
}

type Template struct {
	name       string
	layoutBody string
	tm         *Manager
}

func (s *Template) HTML(code int, context *gin.Context, obj any) {
	s.HTMLWithErr(nil, code, context, obj)
}

func (s *Template) HTMLWithErr(err error, code int, context *gin.Context, obj any) {
	var name string
	var rerr error
	if context.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		if context.GetHeader("X-Layout") != "" {
			name, rerr = s.tm.RenderViewByNameWithCustomLayout(s.name, context.GetHeader("X-Layout"))
			if rerr != nil {
				panic(rerr)
			}
		}
	} else {
		name, rerr = s.tm.RenderViewByName(s.name)
		if rerr != nil {
			panic(rerr)
		}
	}
	context.Header("X-Template", name)
	context.HTML(code, name, s.tm.contextWrapper(context, obj, err))
}

func (s *Template) ToString(c *gin.Context, obj any) (res string, err error) {
	var b bytes.Buffer
	var v string
	if s.layoutBody == "" {
		v, err = s.tm.RenderViewByName(s.name)
		if err != nil {
			return
		}
	} else {
		v, err = s.tm.RenderViewByNameWithCustomLayout(s.name, s.layoutBody)
		if err != nil {
			return
		}
	}
	re, _ := s.tm.re.Instance(v, s.tm.contextWrapper(c, obj, nil)).(render.HTML)
	err = re.Template.Execute(&b, re.Data)
	if err != nil {
		return
	}
	res = b.String()
	return

}

func (s *Manager) MakeTemplate(name string) *Template {
	return &Template{
		name: name,
		tm:   s,
	}
}

func (s *Manager) MakeTemplateWithLayout(name string, layoutBody string) *Template {
	return &Template{
		name:       name,
		layoutBody: layoutBody,
		tm:         s,
	}
}
