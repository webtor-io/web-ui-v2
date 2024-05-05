package template

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

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
	if s.LayoutBody != "" {
		t, err = template.New(s.Name).Parse(s.LayoutBody)
		if err != nil {
			return nil, err
		}
	} else if s.Layout != "" {
		templates = append(templates, s.LayoutPath)
		t = template.New(filepath.Base(s.LayoutPath))
	} else {
		t = template.New(filepath.Base(s.Path))
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

	} else if s.Layout != "" {
		name += "_" + s.Layout
	}
	return name
}

func (s *View) Render() (string, error) {
	s.once.Do(func() {
		s.templateName = s.makeTemplateName()
		var t *template.Template
		t, s.err = s.makeTemplate()
		if s.err != nil {
			return
		}
		s.re.Add(s.templateName, t)
	})
	return s.templateName, s.err
}

type Context struct {
	Data any
	Err  error
}

func NewContext(c *gin.Context, obj any, err error) any {
	return &Context{
		Data: obj,
		Err:  err,
	}
}

type Manager struct {
	re             multitemplate.Renderer
	funcs          FuncMap
	contextWrapper func(c *gin.Context, obj any, err error) any
	layouts        []string
	partials       []string
	views          []*View
	mux            sync.Mutex
}

func NewManager(re multitemplate.Renderer) *Manager {
	return &Manager{
		re:             re,
		funcs:          FuncMap{},
		contextWrapper: NewContext,
	}
}

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func (s *Manager) MustRegisterViews(pattern string) *Manager {
	return Must(s.RegisterViews(pattern))
}

func (s *Manager) RegisterViews(pattern string) (m *Manager, err error) {
	m = s
	views, err := filepath.Glob("templates/views/" + pattern)
	if err != nil {
		return
	}
	layouts, err := s.GetLayouts()
	if err != nil {
		return
	}
	partials, err := s.GetPartials()
	if err != nil {
		return
	}
	for _, v := range views {
		f, _ := os.Stat(v)
		if f.IsDir() {
			continue
		}
		s.views = append(s.views, s.makeView(v, "", partials, s.funcs))
		for _, l := range layouts {
			s.views = append(s.views, s.makeView(v, l, partials, s.funcs))
		}
	}

	return
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

func (s *Manager) makeViewWithLayoutBody(mv *View, layout string) *View {
	cv := s.makeView(mv.Path, mv.LayoutPath, mv.Partials, mv.Funcs)
	cv.LayoutBody = layout
	return cv
}

func (s *Manager) WithFuncs(f FuncMap) *Manager {
	for k, v := range f {
		s.funcs[k] = v
	}
	return s
}

func (s *Manager) firstToLower(in string) string {
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError && size <= 1 {
		return in
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return in
	}
	return string(lc) + in[size:]
}

func (s *Manager) WithContextWrapper(cw func(c *gin.Context, obj any, err error) any) *Manager {
	s.contextWrapper = cw
	return s
}

func (s *Manager) WithHelper(h any) *Manager {
	fooType := reflect.TypeOf(h)
	for i := 0; i < fooType.NumMethod(); i++ {
		method := fooType.Method(i)
		s.funcs[s.firstToLower(method.Name)] = func(args ...any) any {
			args = append([]any{h}, args...)
			inputs := make([]reflect.Value, len(args))
			for i := range args {
				inputs[i] = reflect.ValueOf(args[i])
			}
			return method.Func.Call(inputs)[0].Interface()
		}
	}
	return s
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

func (s *Manager) RenderViewByNameAndLayout(name string, layout string) (string, error) {
	for _, v := range s.views {
		if v.Name == name && v.Layout == layout {
			return s.renderView(v)
		}
	}
	return "", errors.New("view not found")
}

func (s *Manager) RenderViewByNameAndLayoutBody(name string, layout string) (string, error) {
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
		cv = s.makeViewWithLayoutBody(mv, layout)
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
	layout     string
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
			name, rerr = s.tm.RenderViewByNameAndLayoutBody(s.name, context.GetHeader("X-Layout"))
			if rerr != nil {
				panic(rerr)
			}
		}
	} else {
		name, rerr = s.tm.RenderViewByNameAndLayout(s.name, s.layout)
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
		v, err = s.tm.RenderViewByNameAndLayout(s.name, s.layout)
		if err != nil {
			return
		}
	} else {
		v, err = s.tm.RenderViewByNameAndLayoutBody(s.name, s.layoutBody)
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

func (s *Manager) Build(name string) *Template {
	return &Template{
		name: name,
		tm:   s,
	}
}

func (s *Template) WithLayout(name string) *Template {
	s.layout = name
	return s
}

func (s *Template) WithLayoutBody(body string) *Template {
	s.layoutBody = body
	return s
}

func Must(m *Manager, err error) *Manager {
	if err != nil {
		panic(err)
	}
	return m
}

type BuilderWithLayout struct {
	tm     *Manager
	layout string
}

func (s *BuilderWithLayout) Build(name string) *Template {
	return s.tm.Build(name).WithLayout(s.layout)
}

type Builder interface {
	Build(name string) *Template
}

func (m *Manager) WithLayout(name string) *BuilderWithLayout {
	return &BuilderWithLayout{
		tm:     m,
		layout: name,
	}

}
