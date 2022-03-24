# Employee CRUD DB

Employee CRUD database.

## Installing

```
git clone https://github.com/radmirid/grpc-logger.git
```

## Setting environment variables in the .env file  

```
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_NAME=postgres
export DB_SSLMODE=disable
export DB_PASSWORD=password
```

## Building & Running

```
source .env && go build -o crud cmd/main.go && ./crud
```

## Running PostgreSQL in Docker

```
docker run -d --name db -e POSTGRES_PASSWORD=password -v ${HOME}/pgdata/:/var/lib/postgresql/data -p 5432:5432 postgres
```

## LICENSE

[MIT License](LICENSE)
