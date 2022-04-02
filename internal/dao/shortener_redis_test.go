package dao

import (
	"context"

	"github.com/go-redis/cache/v8"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisShortenDAO", func() {
	var redisShortenDAO *redisShortenerDAO
	var pgShortenDAO *pgShortenerDAO
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		pgShortenDAO = NewPGShortenerDAO(pgClient)
		redisShortenDAO = NewRedisShortenerDAO(redisClient, pgShortenDAO)
	})

	Describe("Get", func() {
		var (
			req       *Shortener
			shortener *Shortener
			err       error
		)

		BeforeEach(func() {
			req = NewFakeShortener("fake_url")
		})

		JustBeforeEach(func() {
			shortener, err = redisShortenDAO.Get(ctx, req)
		})

		Context("cache hit", func() {
			AfterEach(func() {
				deleteShortenerInRedis(ctx, redisShortenDAO, req)
			})

			When("success", func() {
				BeforeEach(func() {
					insertShortenerInRedis(ctx, redisShortenDAO, req)
				})

				It("returns the shortener with no error", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(shortener).To(Equal(req))
				})
			})

			When("expiredAt expires", func() {
				BeforeEach(func() {
					req.ExpiredAt = expiredAt
					insertShortenerInRedis(ctx, redisShortenDAO, req)
				})

				It("returns ErrExpiredAt", func() {
					Expect(err).To(MatchError(ErrExpiredat))
					Expect(shortener).To(BeNil())
				})
			})
		})

		Context("cache miss", func() {
			When("success", func() {
				BeforeEach(func() {
					insertShortener(req)
				})

				AfterEach(func() {
					deleteShortener(req.ID)
					deleteShortenerInRedis(ctx, redisShortenDAO, req)
				})

				It("returns the shortener with no error", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(shortener).To(Equal(req))
				})
			})

			When("expiredAt expires", func() {
				BeforeEach(func() {
					req.ExpiredAt = expiredAt
					insertShortener(req)
				})

				AfterEach(func() {
					deleteShortener(req.ID)
					deleteShortenerInRedis(ctx, redisShortenDAO, req)
				})

				It("returns ErrExpiredAt", func() {
					Expect(err).To(MatchError(ErrExpiredat))
					Expect(shortener).To(BeNil())
				})
			})

			When("shorten url not found", func() {
				It("returns ErrShortenURLNotFound error", func() {
					Expect(err).To(MatchError(ErrShortenURLNotFound))
					Expect(shortener).To(BeNil())
				})
			})
		})
	})
})

func insertShortenerInRedis(ctx context.Context, shortenerDAO *redisShortenerDAO, shortener *Shortener) {
	Expect(shortenerDAO.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   getShortenerURL(shortener.ShortenURL),
		Value: shortener,
		TTL:   shortenerDAORedisCacheDuration,
	})).NotTo(HaveOccurred())
}

func deleteShortenerInRedis(ctx context.Context, shortenerDAO *redisShortenerDAO, shortener *Shortener) {
	Expect(shortenerDAO.cache.Delete(ctx, getShortenerURL(shortener.ShortenURL))).NotTo(HaveOccurred())
}
