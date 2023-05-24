package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mehix/pricelists/service/pricelist"
)

func init() {
	godotenv.Load()
}

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

func TestNotImplementedIfNoService(t *testing.T) {

	app := New()
	srvr := httptest.NewTLSServer(app.Handlers())
	defer srvr.Close()

	url := fmt.Sprintf("%s/prices/prod/%d/brand/%s/date/%s/time/%s", srvr.URL, 12345, "ZARA", "2020-06-14", "15:00:00")

	resp, err := srvr.Client().Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotImplemented {
		t.Fatalf("wrong response code. expected: %d, got: %d", http.StatusNotImplemented, resp.StatusCode)
	}
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

	// start the application with supporting services
	svc, err := pricelist.NewServiceForH2(os.Getenv("DB_URL"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(svc.Close)

	app := New(WithPricelist(svc))
	srvr := httptest.NewTLSServer(app.Handlers())
	t.Cleanup(srvr.Close)

	// run tests
	for _, s := range scenarios {
		url := fmt.Sprintf("%s/prices/prod/%d/brand/%s/date/%s/time/%s", srvr.URL, s.productID, s.brand, s.date, s.time)
		t.Run(url, func(t *testing.T) {
			s := s

			t.Parallel()

			resp, err := srvr.Client().Get(url)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("wrong response code. expected: %d, got: %d", http.StatusOK, resp.StatusCode)
			}

			expCT := "application/json"
			gotCT := resp.Header.Get("Content-Type")
			if gotCT != expCT {
				t.Fatalf("wrong Content-Type header. expected: %s, got: %s", expCT, gotCT)
			}

			var price PriceResponse
			if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
				t.Fatal(err)
			}

			if price.Price != s.finalPrice {
				t.Fatalf("wrong price. expected: %d, got: %d", s.finalPrice, price.Price)
			}
		})
	}

}
