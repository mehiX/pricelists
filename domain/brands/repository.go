package brands

import "context"

type Brand struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Repository interface {
	Close()
	ListAll(ctx context.Context) ([]Brand, error)
	FindByName(ctx context.Context, name string) (Brand, error)
}
