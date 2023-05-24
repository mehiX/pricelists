package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	server "github.com/mehix/pricelists"
	"github.com/mehix/pricelists/service/pricelist"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	fmt.Println("Prices API")

	addr := os.Args[1]

	svc, err := pricelist.NewServiceForH2(os.Getenv("PRICELIST_DB_URL"), os.Getenv("BRANDS_DB_URL"))
	if err != nil {
		panic(err)
	}
	defer svc.Close()

	app := server.New(server.WithPricelist(svc))

	srvr := http.Server{
		Addr:    addr,
		Handler: app.Handlers(),
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch

		tc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		fmt.Print("\033[K\rClosing down... ")
		if err := srvr.Shutdown(tc); err != nil {
			if err == context.Canceled {
				fmt.Println("Context expired waiting for server to shutdown")
			} else {
				fmt.Printf("shutting down server: %v\n", err)
			}
		} else {
			fmt.Println("done!")
		}
	}()

	fmt.Printf("Listen on %s\n", srvr.Addr)
	fmt.Printf("Press CTRL+C to end\n")
	if err := srvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-done
	fmt.Println("Exit")
}

func printUsage() {
	fmt.Printf(`Prices API
	
Usage: %s <http address>`, os.Args[0])
}
