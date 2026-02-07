# grooster-backend

A GEWIS rooster (roster/schedule) maker. This is the backend repository.

## Setup

Make sure you can run Go and have this repository open in your favourite IDE.

Install the required packages:

```
go mod tidy
```

## Configuration

Copy the example environment file:

```
cp .env-example .env
```

Ensure `DEV_TYPE` is set to "local" in your `.env` file if you want to run the project without Keycloak authentication.

Make sure to set the `JWT_SECRET`

## Seeding the Database

If you want to seed the database with initial data (users, organs, rosters), run:

```
go run ./cmd/seeder/main.go
```

## Running the Application

To start the application, run:

```
go run ./cmd/src/main.go
```

## Generating Docs and Client

To generate the documentation (Swagger/OpenAPI) from the root directory, run:

```
swag init -d cmd/src -g main.go -o cmd/src/docs --pd --parseInternal
```

To generate the client (which can be uploaded to npm) from the root directory, run:

```
openapi-generator-cli generate -i cmd/src/docs/swagger.yaml -g typescript-axios -o client/src
```

## Running Locally

1. Ensure your `.env` is configured (see Configuration section above).
2. Seed the database if necessary.
3. Start the server:

```
go run ./cmd/src/main.go
```