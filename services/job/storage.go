package job

import (
	"context"
	"time"

	cs "github.com/webtor-io/common-services"
)

type State struct {
	ID  string
	TTL time.Duration
}

type Storage interface {
	Pub(ctx context.Context, id string, l *LogItem) error
	Sub(ctx context.Context, id string) (res chan LogItem, err error)
	GetState(ctx context.Context, id string) (state *State, err error)
	Drop(ctx context.Context, id string) (err error)
}

type NilStorage struct{}

func (s *NilStorage) Pub(_ context.Context, _ string, _ *LogItem) error {
	return nil
}

func (s *NilStorage) Drop(_ context.Context, _ string) (err error) {
	return
}

func (s *NilStorage) Sub(_ context.Context, _ string) (res chan LogItem, err error) {
	return
}

func (s *NilStorage) GetState(_ context.Context, _ string) (state *State, err error) {
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
