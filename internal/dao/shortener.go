package dao

import (
	"context"
	"errors"
	"time"
)

type URL struct {
	ID         string
	Url        string
	ShortenUrl string
	ExpiredAt  time.Time
}

type URLDAO interface {
	Shorten(ctx context.Context, req URL) (string, string, error)
	Get(ctx context.Context, req URL) (string, error)
}

var (
	ExpiredURLError = errors.New("id has alread expired")
)

type Test struct {
	a int
}

var _ URLDAO = (*Test)(nil)

func NewTestDAO() *Test {
	return &Test{
		a: 10,
	}
}

func (t *Test) Shorten(ctx context.Context, req URL) (string, string, error) {
	x := "a"
	y := "b"

	return x, y, nil
}

func (t *Test) Get(ctx context.Context, req URL) (string, error) {
	x := req.ID

	return x, nil
}
