package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
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
	// check time if it has alread expired
	now := time.Now().Unix()
	expiredAt, _ := time.Parse(time.RFC3339, req.ExpiredAt)

	if now >= expiredAt.Unix() {
		return uuid.Nil, "", ErrShortenURLFail
	}

	id := uuid.New()
	req.ID = id
	req.ShortenURL = "http://" + req.ShortenURL + "/get/" + id.String()

	if _, err := dao.client.ModelContext(ctx, req).Insert(); err != nil {
		return uuid.Nil, "", err
	}

	return req.ID, req.ShortenURL, nil
}

func (dao *pgShortenerDAO) Get(ctx context.Context, req *Shortener) (*Shortener, error) {
	var shortener = &Shortener{}

	if err := dao.client.ModelContext(ctx, shortener).Where("id = ?", req.ID).Select(); err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, ErrShortenURLNotFound
		}
		return nil, err
	}

	now := time.Now().Unix()
	expirtedAt, _ := time.Parse(time.RFC3339, shortener.ExpiredAt)

	if now > expirtedAt.Unix() {
		return nil, ErrExpiredat
	}

	return shortener, nil
}
