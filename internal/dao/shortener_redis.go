package dao

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/google/uuid"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/rediskit"
)

type redisShortenerDAO struct {
	cache   *cache.Cache
	baseDAO ShortenerDAO
}

var _ ShortenerDAO = (*redisShortenerDAO)(nil)

const (
	shortenerDAOLocalCacheSize     = 1024
	shortenerDAOLocalCacheDuration = 1 * time.Minute
	shortenerDAORedisCacheDuration = 3 * time.Minute
)

func NewRedisShortenerDAO(client *rediskit.RedisClient, baseDAO ShortenerDAO) *redisShortenerDAO {
	return &redisShortenerDAO{
		cache: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: cache.NewTinyLFU(shortenerDAOLocalCacheSize, shortenerDAOLocalCacheDuration),
		}),
		baseDAO: baseDAO,
	}
}

func (dao *redisShortenerDAO) Shorten(ctx context.Context, req *Shortener) (uuid.UUID, string, error) {
	return dao.baseDAO.Shorten(ctx, req)
}

func (dao *redisShortenerDAO) Get(ctx context.Context, req *Shortener) (string, error) {
	var shortener Shortener

	if err := dao.cache.Once(&cache.Item{
		Key:   req.ShortenUrl,
		Value: &shortener,
		TTL:   shortenerDAORedisCacheDuration,
		Do: func(*cache.Item) (interface{}, error) {
			return dao.baseDAO.Get(ctx, req)
		},
	}); err != nil {
		return "", err
	}

	now := time.Now()

	if now.After(shortener.ExpiredAt) {
		return "", ErrExpiredat
	}

	return shortener.Url, nil
}
