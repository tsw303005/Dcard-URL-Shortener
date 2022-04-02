package dao

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/logkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/pgkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/rediskit"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test DAO")
}

var (
	pgClient    *pgkit.PGClient
	redisClient *rediskit.RedisClient
)

var _ = BeforeSuite(func() {
	pgConf := &pgkit.PGConfig{
		URL: "postgres://postgres@postgres:5432/postgres?sslmode=disable",
	}
	if url := os.Getenv("POSTGRES_URL"); url != "" {
		pgConf.URL = url
	}

	redisConf := &rediskit.Redisconfig{
		Addr: "redis:6379",
	}

	ctx := logkit.NewLogger(&logkit.LoggerConfig{
		Development: true,
	}).WithContext(context.Background())

	pgClient = pgkit.NewPGClient(ctx, pgConf)
	redisClient = rediskit.NewRedisClient(ctx, redisConf)
})

var _ = AfterSuite(func() {
	Expect(pgClient.Close()).ToNot(HaveOccurred())
	Expect(redisClient.Close()).ToNot(HaveOccurred())
})
