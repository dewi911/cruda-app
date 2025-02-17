# cruda-app

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USERNAME=postgres
export DB_NAME=postgres
export DB_SSLMODE=disable
export DB_PASSWORD=qwerty123

```


```bash
source .env && go build -o app cmd/main.go && ./app

```

```bash
docker run -d --name cruda-db -e POSTGRES_PASSWORD=qwerty123 -v ${HOME}/pgdata/:/var/lib/postgresql/data -p 5432:5432 postgres

```

or

```bash
docker run -d --name=crudd-db -e POSTGRES_PASSWORD=qwerty  -p 5432:5432 -d --rm postgres

```

migrates tables
```bash
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up
```
```bash
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' down
```