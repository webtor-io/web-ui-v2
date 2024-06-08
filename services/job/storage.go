package job

import (
	"context"
	"time"

	cs "github.com/webtor-io/common-services"
)

type JobState struct {
	ID  string
	TTL time.Duration
}

type Storage interface {
	Pub(ctx context.Context, id string, l *LogItem) error
	Sub(ctx context.Context, id string) (res chan LogItem, err error)
	GetState(ctx context.Context, id string) (state *JobState, err error)
}

type NilStorage struct{}

func (s *NilStorage) Pub(ctx context.Context, id string, l *LogItem) error {
	return nil
}

func (s *NilStorage) Sub(ctx context.Context, id string) (res chan LogItem, err error) {
	return
}

func (s *NilStorage) GetState(ctx context.Context, id string) (state *JobState, err error) {
	return nil, nil
}

var _ Storage = (*NilStorage)(nil)

func NewStorage(rc *cs.RedisClient, prefix string) Storage {
	cl := rc.Get()
	if cl == nil {
		return &NilStorage{}
	}
	return NewRedis(cl, prefix)
}
