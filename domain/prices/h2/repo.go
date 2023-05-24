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

	if conn != nil {
		if _, err := conn.Exec(`CREATE TABLE PRICE_LIST (
			brand_id int not null,
			start_date datetime not null,
			end_date datetime not null,
			price_list int not null,
			product_id int not null,
			priority int not null default 0,
			price int not null,
			currency varchar(3) not null default 'EUR'
		)`); err != nil {
			log.Println("Creating PRICE_LIST", err)
		} else {
			if _, err := conn.Exec(`insert into PRICE_LIST values (
				1, '2020-06-14 00:00:00', '2020-12-31 23:59:59', 1, 35455, 0, 3550, 'EUR'
			)`); err != nil {
				log.Println("Inserting into PRICE_LIST", err)
			} else {
				fmt.Println("PRICE_LIST populated with data")
			}
		}
	}

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
	from PRICE_LIST 
	where ? >= start_date and ? <= end_date and brand_id = ? and product_id = ? 
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
