package dao

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tsw303005/Dcard-URL-Shortener/pkg/pgkit"
)

type pgShortenerDAO struct {
	client *pgkit.PGClient
}

var _ ShortenerDAO = (*pgShortenerDAO)(nil)

func NewPGShortenerDAO(pgClient *pgkit.PGClient) *pgShortenerDAO {
	return &pgShortenerDAO{
		client: pgClient,
	}
}

func (dao *pgShortenerDAO) Shorten(ctx context.Context, req *Shortener) (uuid.UUID, string, error) {
	id, err := uuid.NewUUID()

	if err != nil {
		log.Fatal("failed to create uuid", err)
	}

	req.ID = id
	req.ShortenUrl = "http://localhost/" + id.String()

	if _, err := dao.client.ModelContext(ctx, req).Insert(); err != nil {
		return uuid.Nil, "", ErrShortenURLFail
	}

	return req.ID, req.ShortenUrl, nil
}

func (dao *pgShortenerDAO) Get(ctx context.Context, req *Shortener) (*Shortener, error) {
	var shortener = &Shortener{}

	if err := dao.client.ModelContext(ctx, shortener).Where("shorten_url = ?", req.ShortenUrl).Select(); err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	expirted_at, _ := time.Parse(time.RFC3339, shortener.ExpiredAt)

	if now > expirted_at.Unix() {
		return nil, ErrExpiredat
	}

	return shortener, nil
}
