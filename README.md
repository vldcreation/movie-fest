# Backend Engineering Interview Assignment (Golang) Movie Festival APp

## Preview
This is a sample project for backend engineering interview.
If you find these routing isn't best practice, because currently open api doesn't support middleware per endpoint for auto generated endpoint.
That's why i separate the routing into different file. so each routing can use it's own middleware. see more [issue](https://github.com/oapi-codegen/oapi-codegen/issues/518)

It's still not as clean as i want, but we achieve fully testable, automatically generated API, automatically generate querier so it's safe to use CQRS pattern and easy to maintain if so.

All these endpoint is untested propperly due to timeline conflicts

## Requirements

To run this project you need to have the following installed:

1. [Go](https://golang.org/doc/install) version 1.21
2. [GNU Make](https://www.gnu.org/software/make/)
3. [oapi-codegen](https://github.com/deepmap/oapi-codegen)
    Install the latest version with:
    ```
    go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
    ```
4. [golang-migrate](https://github.com/golang-migrate/migrate)
    Install the latest version with:
    ```
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
    [sqlc](github.com/kyleconroy/sqlc)
    Install the latest version with:
    ```
    go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
    ```
5. [golangci-lint](https://github.com/deepmap/oapi-codegen)
    Install the latest version with:
    ```
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    ```
6. [mock](https://github.com/uber-go/mock)

    Install the latest version with:
    ```
    go install go.uber.org/mock/mockgen@latest
    ```

7. [Docker](https://docs.docker.com/get-docker/) version 20

8. [Docker Compose](https://docs.docker.com/compose/install/) version 1.29

9. [Node](https://nodejs.org/en) v20

10. [NPM](https://www.npmjs.com/) v10

## Initiate The Project

1. To start working, execute

```
make init
```

2. cp config/config.yaml.example config/config.yaml


## Running

Run using the script `run.sh`:

```bash
./run.sh
```

You may see some errors since you have not created the API yet.

However for testing, you can use Docker run the project, run the following command:

```
docke -compose up --build
```

You should be able to access the API at http://localhost:8080

If you change `database.sql` file, you need to reinitate the database by running:

```
docker compose down --volumes
```

## Testing

To run test, run the following command:

```
make test
```

```
cp rest_client.http.example rest_client.http
```
