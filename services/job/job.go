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
	storage   Storage
	main      bool
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

func New(ctx context.Context, id string, queue string, runnable Runnable, storage Storage) *Job {
	return &Job{
		ID:        id,
		Queue:     queue,
		runnable:  runnable,
		Context:   ctx,
		l:         []LogItem{},
		observers: []*Observer{},
		storage:   storage,
		main:      true,
	}
}

func (s *Job) Run(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("job panic: %v", r)
		}
	}()
	items, err := s.storage.Sub(ctx, s.ID)
	if err != nil {
		return err
	}
	if items != nil {
		s.main = false
		for i := range items {
			err = s.log(i)
			if err != nil {
				return err
			}
		}
		return
	}
	if s.runnable != nil {
		return s.runnable.Run(s)
	}
	return
}

func (s *Job) ObserveLog() *Observer {
	o := NewObserver()
	s.observers = append(s.observers, o)
	for _, i := range s.l {
		o.Push(i)
		if i.Level == Close {
			o.Close()
		}
	}
	return o
}

func (s *Job) log(l LogItem) error {
	l.Timestamp = time.Now()
	s.l = append(s.l, l)
	for _, o := range s.observers {
		o.Push(l)
		if l.Level == Close {
			o.Close()
		}
	}
	if s.main {
		err := s.storage.Pub(s.Context, s.ID, &l)
		if err != nil {
			return err
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
	return nil
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
	queue   string
	mux     sync.Mutex
	jobs    map[string]*Job
	storage Storage
}

func newJobs(queue string, storage Storage) *Jobs {
	return &Jobs{
		queue:   queue,
		jobs:    map[string]*Job{},
		storage: storage,
	}
}

func (s *Jobs) Enqueue(ctx context.Context, id string, r Runnable) *Job {
	s.mux.Lock()
	defer s.mux.Unlock()
	if _, ok := s.jobs[id]; ok {
		return s.jobs[id]
	}
	j := New(ctx, id, s.queue, r, s.storage)
	s.jobs[id] = j
	go func() {
		err := j.Run(ctx)
		if err != nil {
			log.WithError(err).Error("got job error")
		}
		j.Close()
		<-ctx.Done()
		s.mux.Lock()
		defer s.mux.Unlock()
		delete(s.jobs, id)
	}()
	return j
}

func (s *Jobs) Log(ctx context.Context, id string) (c chan LogItem, err error) {
	c = make(chan LogItem)
	j, ok := s.jobs[id]
	if !ok {
		var state *JobState
		state, err = s.storage.GetState(ctx, id)
		if err != nil || state == nil {
			close(c)
			return
		}
		jCtx := context.Background()
		if (state.TTL) > 0 {
			jCtx, _ = context.WithTimeout(context.Background(), state.TTL)
		}
		j = s.Enqueue(jCtx, id, nil)
	}
	go func() {
		o := j.ObserveLog()
		for i := range o.C {
			c <- i
		}
		close(c)
	}()
	return
}

type Queues struct {
	jobs    map[string]*Jobs
	storage Storage
}

var queueMux sync.Mutex

func NewQueues(storage Storage) *Queues {
	return &Queues{
		jobs:    map[string]*Jobs{},
		storage: storage,
	}
}

func (s Queues) GetOrCreate(name string) *Jobs {
	queueMux.Lock()
	defer queueMux.Unlock()
	_, ok := s.jobs[name]
	if !ok {
		s.jobs[name] = newJobs(name, s.storage)
	}
	return s.jobs[name]
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
