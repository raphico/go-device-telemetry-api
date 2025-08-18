# Go Device Telemetry API

A backend service built with **Go** for managing IoT devices, collecting telemetry data, and sending remote commands. Includes authentication, CRUD operations, and CI setup.

## Documentation

- [Project Motivation](./docs/motivation.md)
- [Database Design](./docs/database.md)

## Features

- **User Authentication** with JWT
- **Device Management** (CRUD)
- **Telemetry Collection** (time-series sensor data)
- **Command Dispatch** to devices
- **CI** pipeline (GitHub Actions)

## Tech Stack

- **Language:** Go (Golang)
- **Router:** Chi
- **Database:** PostgreSQL + pgx + Goose SQL migrations

## Quick Start

1. Clone the repository

```bash
git clone git@github.com:raphico/go-device-telemetry-api.git
cd go-device-telemetry-api
```

2. Set environmental variables

```bash
export DATABASE_URL="postgres://telemetry_user:telemetry_pass@localhost:5432/telemetry_db"
export HTTP_PORT="8080"
export JWT_SECRET=$(openssl rand -base64 32)
```

3. Run server

```bash
go run cmd/api/main.go
```

## License

Licensed under the MIT License. Check the [LICENSE](./LICENSE) file for details.
