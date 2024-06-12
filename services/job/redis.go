package job

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	cl     redis.UniversalClient
	prefix string
}

func NewRedis(cl redis.UniversalClient, prefix string) *Redis {
	return &Redis{
		prefix: prefix,
		cl:     cl,
	}
}

func (s *Redis) Pub(ctx context.Context, id string, l *LogItem) (err error) {
	key := s.makeKey(id)
	j, err := json.Marshal(l)
	if err != nil {
		return err
	}

	cmd := s.cl.RPush(ctx, key, j)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	cmd = s.cl.Publish(ctx, key, string(j))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return
}

func (s *Redis) GetState(ctx context.Context, id string) (state *JobState, err error) {
	key := s.makeKey(id)
	ttlCmd := s.cl.TTL(ctx, key)
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
	key := s.makeKey(id)
	exCmd := s.cl.Exists(ctx, key)
	ex, err := exCmd.Result()
	if err != nil {
		return
	}
	if ex == 0 {
		if deadline, ok := ctx.Deadline(); ok {
			pCmd := s.cl.LPush(ctx, key, "")
			if err = pCmd.Err(); err != nil {
				return
			}
			eCmd := s.cl.ExpireAt(ctx, key, deadline)
			if err = eCmd.Err(); err != nil {
				return
			}
		}
		return
	}
	cmd := s.cl.LRange(ctx, key, 0, -1)
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
		ps := s.cl.Subscribe(ctx, key)
		defer ps.Close()

		if err = ps.Ping(ctx); err != nil {
			return
		}

		for m := range ps.Channel() {
			res <- m.Payload
		}
		close(res)
	}()
	return
}

func (s *Redis) makeKey(id string) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s%s", s.prefix, id)))
	return hex.EncodeToString(hash[:])
}

var _ Storage = (*Redis)(nil)
