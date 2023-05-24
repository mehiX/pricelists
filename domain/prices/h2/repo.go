package h2

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jmrobles/h2go"

	"github.com/mehix/pricelists/domain/prices"
	"github.com/mehix/pricelists/internal/h2"
)

type repo struct {
	conn *sql.DB
}

func NewRepository(url string) (prices.Repository, error) {
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

func (r *repo) ListAllForTime(ctx context.Context, brandID, productID int64, t time.Time) ([]prices.PriceDetails, error) {
	qry := `select brand_id, start_date, end_date, price_list, product_id, priority, price, currency 
	from pricelist 
	where ? >= startDate and ? <= endDate and brandId = ? and productId = ? 
	order by priority desc;`

	rows, err := r.conn.QueryContext(ctx, qry, t, t, brandID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []prices.PriceDetails
	for rows.Next() {
		var d prices.PriceDetails
		if err := rows.Scan(&d.BrandID, &d.StartDate, &d.EndDate, &d.PriceList, &d.ProductID, &d.Priority, &d.Price, &d.Currency); err != nil {
			log.Printf("scanning sql row: %v\n", err)
			continue
		}
		details = append(details, d)
	}

	return details, nil
}
