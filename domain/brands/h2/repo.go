package h2

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mehix/pricelists/domain/brands"
	"github.com/mehix/pricelists/internal/h2"
)

type repo struct {
	conn *sql.DB
}

func NewRepository(url string) (brands.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	f := h2.ConnectWithRetry(h2.ConnectDB, 5, time.Second, time.Minute)
	conn, err := f(ctx, url)

	return &repo{conn: conn}, err
}

func (r *repo) Close() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("closing h2 database: %v\n", err)
		} else {
			fmt.Println("H2 connection closed")
		}
	}
}

func (r *repo) ListAll(ctx context.Context) ([]brands.Brand, error) {
	qry := "select id, name from brands"

	rows, err := r.conn.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []brands.Brand
	for rows.Next() {
		var b brands.Brand
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			log.Printf("scanning brand row: %v\n", err)
			continue
		}
		all = append(all, b)
	}

	return all, nil
}

func (r *repo) FindByName(ctx context.Context, name string) (brands.Brand, error) {

	qry := "select id, name from brands where name = ?"

	row := r.conn.QueryRowContext(ctx, qry, name)
	if row.Err() != nil {
		return brands.Brand{}, row.Err()
	}

	var b brands.Brand
	if err := row.Scan(&b.ID, &b.Name); err != nil {
		return brands.Brand{}, err
	}

	return b, nil
}
