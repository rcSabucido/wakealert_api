# WakeAlert API

WakeAlert API is a Go HTTP service for managing users, victims, contacts, medical information, alerts, and location data for the WakeAlert system.

The service uses PostgreSQL for persistence, sqlc for type-safe queries, and JWT-based authentication for login flows.

## Features

- Auth endpoints for web and mobile users
- Victim, contact, alert, and medical data management
- Philippine location lookup endpoints for regions, provinces, cities, and barangays
- PostgreSQL-backed storage with SQL migrations
- Type-safe query generation with sqlc

## Prerequisites

- Go 1.25 or later
- Docker and Docker Compose
- `migrate` CLI for running database migrations
- `sqlc` CLI if you need to regenerate query code

## Configuration

Create a `.env` file in the project root with at least the following values:

```env
ask the developers for the contents of the env
```

The server loads `.env` automatically when it starts.

## Local Setup

1. Start PostgreSQL:

```bash
make up
```

2. Run the database migrations:

```bash
make migrate-up
```

3. Generate sqlc code if you changed any queries:

```bash
make sqlc
```

4. Start the API server:

```bash
make run
```

The service listens on `http://localhost:8080`.

## Useful Commands

```bash
make up           # Start the local PostgreSQL container
make down         # Stop the local PostgreSQL container
make reset        # Stop PostgreSQL and remove the volume, then start fresh
make migrate-up   # Apply all pending migrations
make migrate-down # Roll back one migration
make sqlc         # Regenerate sqlc output
make run          # Run the API server
make psql         # Open a psql shell using DATABASE_URL
```

## API Overview

### Authentication

- `POST /auth/login`
- `POST /mobile/auth/login`

### Mobile Users

- `POST /mobile_users`
- `GET /mobile_users/{id}`
- `GET /mobile_users/email/{email}`

### Medical Information

- `POST /medical_info/add`
- `PUT /medical_info/update/{medical_info_id}`
- `GET /medical_history/{medical_info_id}`
- `POST /medical_history`
- `DELETE /medical_history/{medical_info_id}`

### Contacts

- `POST /contacts/add`
- `POST /contacts/edit/{client_user_id}/{contact_id}`
- `POST /contacts/edit/{client_user_id}`
- `PUT /contacts/{client_user_id}/primary/clear`
- `GET /contacts/client/{client_user_id}`
- `DELETE /contacts/{client_user_id}/{contact_id}`
- `DELETE /contacts/by_details`

### Alerts

- `GET /alerts`
- `POST /alerts`
- `GET /alerts/{id}`
- `PUT /alerts/{id}`
- `DELETE /alerts/{id}`
- `GET /alerts/victim/{victim_id}`

### Victims

- `GET /victims`
- `POST /victims/add`
- `POST /victims/update/{victim_id}`
- `GET /victims/{id}`
- `GET /victims/mobile_user/{mobile_user_id}`
- `GET /victims/address_id/{victim_id}`
- `PUT /victims/address_id/{victim_id}`

### Addresses

- `GET /addresses/regions`
- `GET /addresses/regions/{region_psgc}/provinces`
- `GET /addresses/provinces/{province_or_huc_psgc}/cities`
- `GET /addresses/cities/{city_mun_psgc}/barangays`
- `POST /addresses/lines`
- `PUT /addresses/lines/{address_id}`
- `GET /addresses/lines/{address_id}`

## Database

Database migrations live in `internal/db/migrations`. sqlc query definitions live in `internal/db/queries`, and generated code is written to `internal/db/sqlc`.

If you add or change SQL, update the query files first and run `make sqlc` to regenerate the Go bindings.

## Developers

- Ryz Clyd Sabucido
- Juan Miguel Agunod
- Dan Blair Bapilar

