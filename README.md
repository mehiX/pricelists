# Price list

A coding challenge

## Requirements

- Docker - not needed if you have a local database
- Go >= 1.19 - not needed if you choose to run the project with `docker compose`

The scripts and Makefile are tested on MacOS and Linux. Windows is not supported at the moment.

Create an environment file by running the following:

```shell
cat > .env <<EOF
PRICELIST_DB_URL="host=localhost port=5435 user=sa dbname=/tmp/pricelist sslmode=disable"
BRANDS_DB_URL="host=localhost port=5435 user=sa dbname=/tmp/brands sslmode=disable"
EOF
```

## Run the tests

```shell
make test
```

This will start the database container and run the tests against it. It will then tear down all the running services.

It also generates the coverage profile: [cover.html](cover.html)

## Run the project with Docker

```shell
make up
```

This will build 2 Docker images (database and application), run the 2 services and expose the necessary ports. The database ports are exposed so that we can run tests with `go test`.

Validate that the server is reachable:

```shell
curl -I -X GET http://localhost:8080/health
```

Make a request for a price details:

```shell
curl -s GET http://localhost:8080/prices/prod/35455/brand/ZARA/date/2020-06-14/time/10:00:00
```

When done, shut down the services:

```shell
make down
```

## Build a binary and run it locally

```shell
make binary db
./dist/prices 127.0.0.1:8080
```

## Considerations and decisions taken

- H2 is a good choice for a Java project. For a Go project it only gives headaches. Luckily I enjoy a challenge, however I would prefer MySQL or Postgres which have better support. In fact, I ended up using the Postgres driver.

- With other databases I would use SQL scripts to initialize the database. Since it would take me longer to investigate how to do that with H2, I chose to initialize the databases inside the code ([domain/brands/h2/repo.go](domain/brands/h2/repo.go) and [domain/brands/h2/repo.go](domain/brands/h2/repo.go))

- Connecting to the database supports **retry** and **exponential backoff** ([internal/h2/connect.go](internal/h2/connect.go)). This is especially useful in a microservices environment where the order in which the services start should not matter. You can see how this works by starting the application without the database and watching the output. You can then start the database and see how the server will connect. As further improvements, the retry and backoff parameters should be customizable and the server should also try reconnection when the database goes down unexpectedly.

- **Brands** and **PriceLists** are separate entities, each with their own database, in the idea that a bigger project would also treat these as separate services with endpoints to also manage the data.

- Price representation - I use prices as `int64` and divide them by 100 only for use by the presentation layer.

- Tests - given the limited time I chose to cover the main functionality with tests ([server_test.go](server_test.go)). The required tests are in the function `TestPrices`. They are more integration tests than just unit tests. Normally they should be separated and ran separately using Go tags.

- Git - I tried to add clear messages to my commits and they should be just as clear as the suggested tags.