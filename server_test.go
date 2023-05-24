package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthEndpoint(t *testing.T) {

	app := New()
	srvr := httptest.NewTLSServer(app.Handlers())
	defer srvr.Close()

	resp, err := srvr.Client().Get(srvr.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		t.Fatal("/health is not defined")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wrong status code. expected: %d, got: %d", http.StatusOK, resp.StatusCode)
	}
}

func mustParse(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestPrices(t *testing.T) {

	type scenario struct {
		date       string
		time       string
		productID  int64
		brand      string
		finalPrice int64
	}

	scenarios := []scenario{
		{date: "2020-06-14", time: "10:00:00", productID: 35455, brand: "ZARA", finalPrice: 3550},
		{date: "2020-06-14", time: "16:00:00", productID: 35455, brand: "ZARA", finalPrice: 2545},
		{date: "2020-06-14", time: "21:00:00", productID: 35455, brand: "ZARA", finalPrice: 3550},
		{date: "2020-06-15", time: "10:00:00", productID: 35455, brand: "ZARA", finalPrice: 3050},
		{date: "2020-06-16", time: "21:00:00", productID: 35455, brand: "ZARA", finalPrice: 33895},
	}

	app := New()
	srvr := httptest.NewTLSServer(app.Handlers())
	t.Cleanup(srvr.Close)

	for _, s := range scenarios {
		url := fmt.Sprintf("%s/prices/prod/%d/brand/%s/date/%s/time/%s", srvr.URL, s.productID, s.brand, s.date, s.time)
		t.Run(url, func(t *testing.T) {
			//s := s

			t.Parallel()

			resp, err := srvr.Client().Get(url)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("wrong response code. expected: %d, got: %d", http.StatusOK, resp.StatusCode)
			}

		})
	}

}
