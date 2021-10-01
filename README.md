# organization_manager

This repo contains code which stands up a REST API server to manage an organization object. 
The following endpoints are available:
 - GET /api/v1/organizations - Retrieves organizations and can be filtered via query parameters
 - POST /api/v1/organizations - Creates a new organizations from the request body

Organization Object:
```markdown
{
    "id": <internal_id> 
    "name": "CLEAR",
    "creation_date": "2002-09-22T00:00:00Z",
    "employee_count": 10000,
    "is_public": true
}
```

For more detailed endpoint documentation see the swagger docs located in `/documentation/api_docs.yaml`

## Running the server:
This server expects to be run on Golang v 1.14 with the dependencies specified in the go.mod file. To run the server locally
follow the steps below.

Install the dependencies:
```shell
go mod download
```

Build the go package:
```shell
go build -ldflags="-w -s" -installsuffix cgo -o  main ./cmd
```

Run Postgres docker-compose:
```shell
docker-compose -f docker-compose-postgres.yml up
```

Run the go server:
```shell
./main
```

## Configuration through environment variables:
The following environment variables are used to configure the server:
- MIGRATIONS_PATH - Path to the folder container golang-migrate files (default: `file://pkg/database/migrations`)
- DATABASE_URL - Connection URL for database the server connects to (default: `postgres://user:Password123!@localhost:5432/organization_service?sslmode=disable`)
- PORT - Port in which the server will listen on (default: `8082`)

## Running unit tests:
```shell
go test ./...
```
