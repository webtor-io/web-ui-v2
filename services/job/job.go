package job

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Observer struct {
	C      chan LogItem
	mux    sync.Mutex
	closed bool
}

func (s *Observer) Close() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if !s.closed {
		s.closed = true
		close(s.C)
	}
}

func (s *Observer) Push(v LogItem) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if !s.closed {
		s.C <- v
	}
}

func NewObserver() *Observer {
	return &Observer{
		C: make(chan LogItem, 100),
	}
}

type Job struct {
	ID        string
	Queue     string
	l         []LogItem
	runnable  Runnable
	observers []*Observer
	closed    bool
	mux       sync.Mutex
	cur       string
	Context   context.Context
}

type LogItemLevel string

const (
	Info           LogItemLevel = "info"
	Error          LogItemLevel = "error"
	Warn           LogItemLevel = "warn"
	Done           LogItemLevel = "done"
	InProgress     LogItemLevel = "inprogress"
	Finish         LogItemLevel = "finish"
	Redirect       LogItemLevel = "redirect"
	Download       LogItemLevel = "download"
	RenderTemplate LogItemLevel = "rendertemplate"
	StatusUpdate   LogItemLevel = "statusupdate"
	Close          LogItemLevel = "close"
)

var levelMap = map[LogItemLevel]log.Level{
	Info:           log.InfoLevel,
	Error:          log.ErrorLevel,
	Warn:           log.WarnLevel,
	Done:           log.InfoLevel,
	InProgress:     log.InfoLevel,
	Finish:         log.InfoLevel,
	Download:       log.InfoLevel,
	Redirect:       log.InfoLevel,
	StatusUpdate:   log.InfoLevel,
	RenderTemplate: log.InfoLevel,
	Close:          log.InfoLevel,
}

type LogItem struct {
	Level     LogItemLevel `json:"level,omitempty"`
	Message   string       `json:"message,omitempty"`
	Tag       string       `json:"tag,omitempty"`
	Location  string       `json:"location,omitempty"`
	Template  string       `json:"template,omitempty"`
	Body      string       `json:"body,omitempty"`
	Timestamp time.Time    `json:"timestamp,omitempty"`
}

func New(ctx context.Context, id string, queue string, runnable Runnable) *Job {
	return &Job{
		ID:        id,
		Queue:     queue,
		runnable:  runnable,
		Context:   ctx,
		l:         []LogItem{},
		observers: []*Observer{},
	}
}

func (s *Job) Run(ctx context.Context) {
	s.runnable.Run(s)
}

func (s *Job) ObserveLog() *Observer {
	o := NewObserver()
	s.observers = append(s.observers, o)
	for _, i := range s.l {
		o.Push(i)
	}
	return o
}

func (s *Job) log(l LogItem) {
	l.Timestamp = time.Now()
	s.l = append(s.l, l)
	for _, o := range s.observers {
		o.Push(l)
		if l.Level == Close {
			o.Close()
		}
	}
	message := l.Message
	if l.Level == Done {
		message = "done"
	}
	if l.Level == Finish {
		message = "finish"
	}
	if l.Level == Redirect {
		message = "redirect"
	}
	if l.Level == StatusUpdate {
		message = "statusupdate"
	}
	if l.Level == RenderTemplate {
		message = "rendertemplate"
	}
	if l.Level == Close {
		message = "close"
	}
	log.WithFields(log.Fields{
		"ID":       s.ID,
		"Queue":    s.Queue,
		"Tag":      l.Tag,
		"Location": l.Location,
		"Template": l.Template,
		"Body":     l.Body,
	}).Log(levelMap[l.Level], message)
}

func (s *Job) Info(message string) *Job {
	s.log(LogItem{
		Level:   Info,
		Message: message,
	})
	return s
}

func (s *Job) Warn(err error, message string) *Job {
	log.WithError(err).Error("got job warning")
	s.log(LogItem{
		Level:   Warn,
		Message: message,
		Tag:     s.cur,
	})
	return s
}

func (s *Job) Error(err error, message string) error {
	log.WithError(err).Error("got job error")
	s.log(LogItem{
		Level:   Error,
		Message: message,
		Tag:     s.cur,
	})
	return err
}

func (s *Job) InProgress(message string) *Job {
	s.cur = message
	s.log(LogItem{
		Level:   InProgress,
		Message: message,
		Tag:     s.cur,
	})
	return s
}

func (s *Job) StatusUpdate(message string) *Job {
	s.log(LogItem{
		Level:   StatusUpdate,
		Message: message,
		Tag:     s.cur,
	})
	return s
}

func (s *Job) Done() *Job {
	s.log(LogItem{
		Level: Done,
		Tag:   s.cur,
	})
	return s
}
func (s *Job) Finish() *Job {
	s.log(LogItem{
		Level:   Finish,
		Message: "success!",
	})
	return s
}

func (s *Job) Download(url string) *Job {
	s.log(LogItem{
		Level:    Download,
		Location: url,
	})
	return s
}

func (s *Job) Redirect(url string) *Job {
	s.log(LogItem{
		Level:    Redirect,
		Location: url,
	})
	return s
}

func (s *Job) RenderTemplate(name string, body string) *Job {
	s.log(LogItem{
		Level:    RenderTemplate,
		Template: name,
		Body:     body,
	})
	return s
}

func (s *Job) FinishWithMessage(m string) *Job {
	s.log(LogItem{
		Level:   Finish,
		Message: m,
	})
	return s
}

func (s *Job) Close() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	s.log(LogItem{
		Level: Close,
	})
}

type Jobs struct {
	queue string
	mux   sync.Mutex
	jobs  map[string]*Job
}

func newJobs(queue string) *Jobs {
	return &Jobs{
		queue: queue,
		jobs:  map[string]*Job{},
	}
}

func (s *Jobs) Enqueue(ctx context.Context, id string, r Runnable) *Job {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.jobs[id]; ok {
		return s.jobs[id]
	}
	j := New(ctx, id, s.queue, r)
	s.jobs[id] = j
	go func() {
		j.Run(ctx)
		j.Close()
		<-ctx.Done()
		s.mux.Lock()
		defer s.mux.Unlock()
		delete(s.jobs, id)
	}()
	return j
}

func (s *Jobs) Log(id string) chan LogItem {
	c := make(chan LogItem, 100)
	if _, ok := s.jobs[id]; ok {
		go func() {
			o := s.jobs[id].ObserveLog()
			for i := range o.C {
				c <- i
			}
			close(c)
		}()
	} else {
		close(c)
	}
	return c
}

type Queues map[string]*Jobs

var queueMux sync.Mutex

func NewQueues() *Queues {
	return &Queues{}
}

func (s Queues) Get(name string) *Jobs {
	return s[name]
}

func (s Queues) GetOrCreate(name string) *Jobs {
	queueMux.Lock()
	defer queueMux.Unlock()
	_, ok := s[name]
	if !ok {
		s[name] = newJobs(name)
	}
	return s[name]
}

type Runnable interface {
	Run(j *Job) error
}

type Script struct {
	body func(j *Job) error
}

func (s *Script) Run(j *Job) error {
	return s.body(j)
}

func NewScript(body func(j *Job) error) *Script {
	return &Script{body: body}
}
