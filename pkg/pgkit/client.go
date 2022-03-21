package pgkit

import (
	"context"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/logkit"
	"go.uber.org/zap"
)

type PGConfig struct {
	URL string `long:"url" env:"URL" description:"the URL of PostgresSQL" required:"true"`
}

type PGClient struct {
	*pg.DB
	closeFunc func()
}

func (c *PGClient) Close() error {
	if c.closeFunc != nil {
		c.closeFunc()
	}
	return c.DB.Close()
}

func NewPGClient(ctx context.Context, conf *PGConfig) *PGClient {
	if url := os.ExpandEnv(conf.URL); url != "" {
		conf.URL = url
	}

	logger := logkit.FromContext(ctx).With(zap.String("url", conf.URL))

	opts, err := pg.ParseURL(conf.URL)
	if err != nil {
		logger.Fatal("failed to parse PostgresSQL url", zap.Error(err))
	}

	db := pg.Connect(opts).WithContext(ctx)
	if err := db.Ping(ctx); err != nil {
		logger.Fatal("failed to ping PostgresSQL", zap.Error(err))
	}

	logger.Info("create PostgresSQL client suceessfully")

	return &PGClient{
		DB: db,
	}
}
