package pgkit

import (
	"context"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
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

	opts, err := pg.ParseURL(conf.URL)
	if err != nil {
		log.Fatal("failed to parse PostgresSQL url", err)
	}

	db := pg.Connect(opts).WithContext(ctx)
	if err := db.Ping(ctx); err != nil {
		log.Fatal("failed to ping PostgresSQL", err)
	}

	return &PGClient{
		DB: db,
	}
}
