# Company Service

This is a simple service that provides a REST API for managing companies.

## API

| Method | URL Pattern     | Action                                          |
| ------ | --------------- | ----------------------------------------------- |
| GET    | /v1/healthcheck | Show application health and version information |
| GET    | /v1/company/:id | Show Company information identified by ID       |
| PATCH  | /v1/company/:id | Patch Company information                       |
| DELETE | /v1/company/:id | Delete a Company                                |
| CREATE | /v1/company     | Create a Company                                |

## Database

Postgres is used as the database for this service. 
Extension uuid-ossp is used for generating UUIDs. 
https://github.com/golang-migrate/migrate is used for database migrations.
The migrations scripts are located in the ./migrations folder.

## Additional points

- [ ] On each mutating operation, an event should be produced.
- [ ] Dockerize the application to be ready for building the production docker image
- [ ] Use docker for setting up the external services such as the database
- [ ] REST is suggested, but GRPC is also an option
- [ ] JWT for authentication
- [ ] Kafka for events
- [ ] Integration tests are highly appreciated
- [ ] Linter
- [ ] Configuration file