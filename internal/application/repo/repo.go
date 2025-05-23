package repo

import (
	"context"
	"errors"
)

var NotExistsKey = errors.New("key does not exist")

type Repo interface {
	Save(ctx context.Context, key string, val []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
}
