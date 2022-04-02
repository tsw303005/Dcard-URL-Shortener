package dao

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/logkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/migrationkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/pgkit"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/rediskit"
)

func TestDAO(t *testing.T) {
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

	migrationConf := &migrationkit.MigrationConfig{
		Source: "file://../../migrations",
		URL:    pgConf.URL,
	}

	redisConf := &rediskit.Redisconfig{
		Addr: "redis:6379",
	}

	ctx := logkit.NewLogger(&logkit.LoggerConfig{
		Development: true,
	}).WithContext(context.Background())

	migration := migrationkit.NewMigration(ctx, migrationConf)
	defer func() {
		Expect(migration.Close()).NotTo(HaveOccurred())
	}()

	Expect(migration.Up()).NotTo(HaveOccurred())

	pgClient = pgkit.NewPGClient(ctx, pgConf)
	redisClient = rediskit.NewRedisClient(ctx, redisConf)
})

var _ = AfterSuite(func() {
	Expect(pgClient.Close()).ToNot(HaveOccurred())
	Expect(redisClient.Close()).ToNot(HaveOccurred())
})

var pgExec = func(query string, params ...interface{}) {
	_, err := pgClient.Exec(query, params...)
	Expect(err).NotTo(HaveOccurred())
}
