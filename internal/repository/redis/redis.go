package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CodeMaster482/ShortLinkAPI/internal/model"
	apierror "github.com/CodeMaster482/ShortLinkAPI/pkg/errors"

	"github.com/go-redis/redis/v8"
)

type LinkRedisStorage struct {
	Client *redis.Client
}

func NewLinkStorage(cli *redis.Client) *LinkRedisStorage {
	return &LinkRedisStorage{cli}
}

func (r *LinkRedisStorage) GetLink(ctx context.Context, token string) (*model.Link, error) {
	fullLink, err := r.Client.Get(ctx, token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) { /* err == redis.Nil */
			return nil, apierror.ErrLinkNotFound
		}

		return nil, err
	}

	return &model.Link{
		OriginalLink: fullLink,
		Token:        token,
	}, nil
}

func (r *LinkRedisStorage) StoreLink(ctx context.Context, link *model.Link) error {
	err := r.Client.Set(ctx, link.Token, link.OriginalLink, 0).Err()
	if err != nil {
		return err
	}

	duration := time.Until(link.ExpiresAt)

	err = r.Client.Expire(ctx, link.Token, duration).Err()
	if err != nil {
		return fmt.Errorf("error setting expiration time for switch %s: %w", link.Token, err)
	}

	return nil
}

func (r *LinkRedisStorage) StartRecalculation(interval time.Duration, deleted chan []string) {
	return
}
