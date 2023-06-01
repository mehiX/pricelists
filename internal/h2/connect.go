package h2

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

// ConnectDbFn is the signature for a retriable function that would (eventually) return a database connection or an error
type ConnectDbFn func(context.Context, string) (*sql.DB, error)

// ConnectWithRetry tries to connect to a database by repeatedly calling `f`.
// This is especially useful when running in container, since the database may take longer to
// start and be ready to receive connections.
// `base` and `cap` are used for exponential backoff. The function will not wait longer than `cap` to retry
func ConnectWithRetry(f ConnectDbFn, retries int, base time.Duration, cap time.Duration) ConnectDbFn {
	backoff := base

	return func(ctx context.Context, s string) (*sql.DB, error) {
		for r := 0; ; r++ {
			conn, err := f(ctx, s)
			if err == nil || r >= retries {
				return conn, err
			}

			if backoff > cap {
				backoff = cap
			}
			jitter := rand.Int63n(int64(backoff * 3))
			wait := backoff + time.Duration(jitter)

			fmt.Printf("Connect attempt %d failed. Waiting %v before retrying\n", r+1, wait)

			select {
			case <-time.After(wait):
				backoff <<= 1
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}

// ConnectDB tries to connect to an H2 instance. It pings the connection before returning it
func ConnectDB(_ context.Context, url string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		log.Printf("Can't connet to H2 Database: %s\n", err)
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		log.Printf("Can't ping to H2 Database: %s\n", err)
		return nil, err
	}

	fmt.Printf("H2 Database connected\n")

	return conn, nil
}
