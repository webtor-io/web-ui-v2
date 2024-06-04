package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	cl redis.UniversalClient
}

func NewRedis(cl redis.UniversalClient) *Redis {
	return &Redis{
		cl: cl,
	}
}

func (s *Redis) Pub(ctx context.Context, id string, l *LogItem) (err error) {
	j, err := json.Marshal(l)
	if err != nil {
		return err
	}

	cmd := s.cl.RPush(ctx, id, j)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	cmd = s.cl.Publish(ctx, id, string(j))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return
}

func (s *Redis) GetState(ctx context.Context, id string) (state *JobState, err error) {
	ttlCmd := s.cl.TTL(ctx, id)
	if ttlCmd.Err() != nil {
		return nil, ttlCmd.Err()
	}
	var dur time.Duration
	val := ttlCmd.Val()
	if val == -2 {
		return
	}
	if val > 0 {
		dur = val
	}
	return &JobState{
		ID:  id,
		TTL: dur,
	}, nil
}

func (s *Redis) Sub(ctx context.Context, id string) (res chan LogItem, err error) {
	ch, err := s.subRaw(ctx, id)
	if err != nil || ch == nil {
		return
	}
	res = make(chan LogItem)
	go func() {
		for i := range ch {
			var li LogItem
			err = json.Unmarshal([]byte(i), &li)
			if err != nil {
				res <- LogItem{
					Level:   Error,
					Message: err.Error(),
				}
				close(res)
				return
			}
			res <- li
		}
		close(res)
	}()
	return
}
func (s *Redis) subRaw(ctx context.Context, id string) (res chan string, err error) {
	exCmd := s.cl.Exists(ctx, id)
	ex, err := exCmd.Result()
	if err != nil {
		return
	}
	if ex == 0 {
		if deadline, ok := ctx.Deadline(); ok {
			pCmd := s.cl.LPush(ctx, id, "")
			if err = pCmd.Err(); err != nil {
				return
			}
			eCmd := s.cl.ExpireAt(ctx, id, deadline)
			if err = eCmd.Err(); err != nil {
				return
			}
		}
		return
	}
	cmd := s.cl.LRange(ctx, id, 0, -1)
	items, err := cmd.Result()
	if err != nil {
		return
	}
	res = make(chan string)
	go func() {
		for _, i := range items {
			if i == "" {
				continue
			}
			res <- i
		}
		ps := s.cl.Subscribe(ctx, id)
		for m := range ps.Channel() {
			res <- m.Payload
		}
		close(res)
	}()
	return
}

var _ Storage = (*Redis)(nil)
