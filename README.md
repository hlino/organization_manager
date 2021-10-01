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

## Running unit tests:
```shell
go test ./...
```
