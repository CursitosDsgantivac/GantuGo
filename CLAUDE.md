# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go web server application demonstrating two HTTP server implementations: a standard library-based server and a Chi router-based server. The application uses command-line arguments to select which server implementation to run.

## Build and Run Commands

The project uses a Makefile with environment-specific configuration:

```bash
# Build and run (default uses ENV=test and ARG=chiServer)
make

# Run with specific environment
make ENV=prod
make ENV=test

# Run with specific server type
make ARG=chiServer    # Chi router server (default)
make ARG=server       # Standard library server
make ARG=utils        # Print environment variables

# Run without loading .env files
make USE_ENV_FILE=false

# Individual build steps
make fmt              # Format code
make vet              # Run static analysis
make build            # Compile binary (creates gantuProgram)
make run              # Build and execute
make clean            # Build, run, then remove binary
```

The Makefile chain: `clean -> run -> build -> vet -> fmt`. Default target is `clean`, which runs everything and removes the binary afterward.

## Docker Commands

```bash
# Run server in Docker
docker-compose up

# Stop server
docker-compose down

# Docker runs the chiServer by default on port 8080
```

## Environment Configuration

The application loads environment-specific variables from `.env.{ENV}` files:
- `.env.test` - Test environment
- `.env.prod` - Production environment

Key environment variables:
- `PORT` - Server port (defaults to 8080 if not set)
- `TEST_ENV` - Example environment variable

The Makefile automatically includes and exports variables from the appropriate `.env.{ENV}` file based on the `ENV` parameter.

## Architecture

### Entry Point ([main.go:1](main.go#L1))
The application requires a command-line argument to determine which mode to run:
- `chiServer` - Starts Chi router-based HTTP server
- `server` - Starts standard library HTTP server
- `utils` - Loads and prints environment variables

### Server Implementations

**Chi Router Server** ([internal/server/chiServer.go:1](internal/server/chiServer.go#L1))
- Uses `github.com/go-chi/chi/v5` router
- Supports route parameters (e.g., `/test/{id}`)
- Returns JSON responses with appropriate content-type headers
- No timeout configuration

**Standard Library Server** ([internal/server/server.go:1](internal/server/server.go#L1))
- Uses `http.ServeMux` from standard library
- Configured with timeouts (Read/Write/Idle: 10s)
- Method-specific route handling using `GET /path` and `POST /path` syntax
- The `{$}` suffix in routes demands exact match

### Utilities ([internal/utils/utils.go:1](internal/utils/utils.go#L1))
- `LoadPort()` - Retrieves PORT environment variable, defaults to "8080"
- `LoadEnvVariables()` - Prints TEST_ENV variable (for debugging)

## Testing

Use the REST Client extension in VSCode with [testRequest.http](testRequest.http):
- GET requests to `/test` endpoint
- POST requests with JSON body to `/test` and `/test/api`
- PUT request with route parameter to `/test/{id}`
- 404 test to `/invented` endpoint
- Default hostname: `http://localhost:8080`

## Module Structure

Module name: `main` (from [go.mod](go.mod))
Go version: 1.25.5
External dependency: `github.com/go-chi/chi/v5 v5.2.4`

Project structure:
- `main.go` - Application entry point with argument routing
- `internal/server/` - HTTP server implementations
- `internal/utils/` - Utility functions for configuration
