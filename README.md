# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Project Structure

1. **cmd/**: Contains the main application and seed command entry points.

   - `server/main.go`: The main application entry point, serves the REST API.
   - `seed/main.go`: Command to seed the database with initial product data.

2. **app/**: Contains the application logic.
3. **sql/**: Contains a very simple database migration scripts setup.
4. **models/**: Contains the data models and repositories used in the application.
5. `.env`: Environment variables file for configuration.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Important makefile targets:
  - `make tidy`: will install all dependencies.
  - `make docker-up`: will start the required infrastructure services via docker containers.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Will run the tests.
  - `make run`: Will start the application.
  - `make docker-down`: Will stop the docker containers.

Assignment: [ASSIGNMENT.md](ASSIGNMENT.md)

## API

Server listens on `http://localhost:8484`.

- `GET /catalog` — `offset` (default 0), `limit` (default 10, max 100), `category`, `priceLessThan`
- `GET /catalog/{code}` — product and variants. Missing variant price uses the product price.
- `GET /categories`
- `POST /categories` — `{"code","name"}`. 201 on create, 409 if the code already exists.

```bash
curl "http://localhost:8484/catalog?category=clothing&limit=5"
curl "http://localhost:8484/catalog/PROD001"
curl "http://localhost:8484/categories"
curl -X POST "http://localhost:8484/categories" \
  -H "Content-Type: application/json" \
  -d '{"code":"electronics","name":"Electronics"}'
```
