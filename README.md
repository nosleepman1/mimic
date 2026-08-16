# Mimic

<p align="center">
  <strong>High-Performance Intelligent API Mock Server & Contract Engine in Go.</strong>
</p>

<p align="center">
  <a href="README.fr.md">Lire en Francais</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#supported-sources">Supported Sources</a> •
  <a href="#cli-commands">CLI Commands</a>
</p>

---

## Overview

**Mimic** is an ultra-fast (< 5ms startup time), zero-dependency CLI tool written in Go that turns API contracts into realistic, stateful, in-memory HTTP mock servers for frontend, mobile (Flutter, React Native, iOS, Android), and QA teams.

---

## Features

- **Microsecond HTTP Latency**: Powered by native Go goroutines, serving 100k+ mock requests/second with zero configuration.
- **Universal Contract Ingestion**: Ingests OpenAPI 3.0/3.1, Swagger 2.0, Postman Collections v2.1, Laravel `routes/api.php`, or declarative `mimic.yaml` files.
- **Stateful In-Memory CRUD**: `POST`, `PUT`, `PATCH`, and `DELETE` requests actually mutate in-memory collections. `GET /resource/:id` retrieves created items dynamically.
- **Smart Faker Engine**: Generates realistic names, emails, UUIDs, ISO timestamps, and localized data (e.g., Senegalese phone numbers `+221` and West African CFA Franc `currency_xof`).
- **Network & Chaos Simulation**:
  - Artificial network latency (`--latency=200ms`)
  - Jitter simulation (`--jitter=50ms`)
  - Chaos 500 error injection (`--chaos=0.05` for 5% failure rate)
- **Built-in CORS**: Instant out-of-the-box support for browser and mobile apps.

---

## Quick Start

### 1. Initialize a Sample Contract
```bash
mimic init
```
This generates a starter `mimic.yaml` template with stateful CRUD, authentication, and localized payment examples.

### 2. Launch the Mock Server
```bash
mimic serve -f mimic.yaml --port=8080 --latency=200ms
```

### 3. Test Endpoints
```bash
# List users
curl http://localhost:8080/api/v1/users

# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Mamadou Diop", "email": "mamadou@example.com"}'

# Retrieve created user
curl http://localhost:8080/api/v1/users/1
```

---

## Supported Sources

| Source Format | Flag / File | Example |
| :--- | :--- | :--- |
| **Declarative YAML/JSON** | `-f mimic.yaml` | Native format with CRUD & Faker mapping |
| **OpenAPI 3.0 / 3.1** | `-f openapi.yaml` | Industry-standard spec |
| **Swagger 2.0** | `-f swagger.json` | Standard Swagger schemas |
| **Postman Collection v2.1** | `-f collection.json` | Direct export from Postman |
| **Laravel Routes** | `-f routes/api.php` | Direct AST parsing of Laravel route files |

---

## CLI Commands

```text
Usage:
  mimic [command]

Available Commands:
  init        Generate a starter mimic.yaml template file
  serve       Start the mock HTTP server from an API contract file
  validate    Check the syntax and validity of an API contract file
  version     Display Mimic version
```

### Serve Flags
- `-f, --file string` : Path to contract file (default `mimic.yaml`).
- `-p, --port int` : Port to listen on (default `8080`).
- `--host string` : Host address to bind to (default `0.0.0.0`).
- `--latency string` : Artificial network latency (e.g., `250ms`, `1s`).
- `--jitter string` : Simulated network jitter (e.g., `50ms`).
- `--chaos float` : Failure rate for 500 error injection (`0.0` to `1.0`).
- `--cors bool` : Enable CORS headers (default `true`).

---

## License

Distributed under the [MIT License](LICENSE).
