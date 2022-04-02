package rediskit

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/logkit"
)

var _ = Describe("RedisClient", func() {
	Describe("NewRedisClient", func() {
		var (
			ctx         context.Context
			redisClient *RedisClient
			redisConf   *Redisconfig
		)

		BeforeEach(func() {
			ctx = logkit.NewLogger(&logkit.LoggerConfig{
				Development: true,
			}).WithContext(context.Background())

			redisConf = &Redisconfig{
				Addr: "localhost:6379",
			}
			if addr := os.Getenv("REDIS_ADDR"); addr != "" {
				redisConf.Addr = addr
			}
		})

		AfterEach(func() {
			Expect(redisClient.Close()).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			redisClient = NewRedisClient(ctx, redisConf)
		})

		When("success", func() {
			It("returns redis cleint without error", func() {
				Expect(redisClient).NotTo(BeNil())
			})
		})
	})
})
