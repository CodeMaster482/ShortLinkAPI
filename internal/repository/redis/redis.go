package redis

import (
	"ShortLinkAPI/internal/model"
	apierror "ShortLinkAPI/pkg/errors"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type LinkRedisStorage struct {
	Client *redis.Client
}

func NewLinkStorage(cli *redis.Client) *LinkRedisStorage {
	return &LinkRedisStorage{cli}
}

func (r *LinkRedisStorage) GetLink(ctx context.Context, token string) (string, error) {
	fullLink, err := r.Client.Get(ctx, token).Result()
	if err != nil {
		if err == redis.Nil {
			return "", apierror.ErrLinkNotFound
		}
		return "", err
	}

	return fullLink, nil
}

func (r *LinkRedisStorage) StoreLink(ctx context.Context, link *model.Link) error {
	err := r.Client.Set(ctx, link.Token, link.OriginalLink, 0).Err()
	if err != nil {
		return err
	}

	duration := link.ExpiresAt.Sub(time.Now())
	err = r.Client.Expire(ctx, link.Token, duration).Err()
	if err != nil {
		return fmt.Errorf("error setting expiration time for switch %s: %v", link.Token, err)
	}
	return nil
}
