# MCP Server Documentation

## Overview

The MCP (Model Context Protocol) Server is a lightweight HTTP server implementation that provides a standardized protocol for MCP-compatible clients to interact with the GantuGo service.

## Starting the Server

To start the MCP server, run:

```bash
go run main.go mcpServer
```

Or using the compiled binary:

```bash
./main mcpServer
```

The server will start on the port specified in your environment configuration (loaded via `PORT` environment variable).

## Architecture

The MCP server is implemented in `internal/server/mcpServer.go` and follows the Model Context Protocol specification. It provides a clean JSON-based API for client-server communication.

### Key Components

- **MCPRequest**: Standard request structure containing method, parameters, and optional ID
- **MCPResponse**: Standard response structure with result or error, and request ID
- **MCPError**: Error structure following JSON-RPC 2.0 error codes

## Endpoints

### 1. Health Check

**Endpoint**: `GET /health`

**Description**: Health check endpoint to verify server status

**Response**:
```json
{
  "status": "healthy",
  "service": "MCP Server",
  "time": 1234567890
}
```

**Example**:
```bash
curl http://localhost:8080/health
```

---

### 2. Server Info

**Endpoint**: `GET /info`

**Description**: Provides information about the MCP server

**Response**:
```json
{
  "name": "GantuGo MCP Server",
  "version": "1.0.0",
  "protocol": "MCP",
  "description": "Model Context Protocol server implementation"
}
```

**Example**:
```bash
curl http://localhost:8080/info
```

---

### 3. MCP Protocol Endpoint

**Endpoint**: `POST /mcp`

**Description**: Main MCP protocol endpoint for handling method calls

**Request Format**:
```json
{
  "method": "methodName",
  "params": {
    "key": "value"
  },
  "id": "optional-request-id"
}
```

**Response Format (Success)**:
```json
{
  "result": {
    "data": "response data"
  },
  "id": "optional-request-id"
}
```

**Response Format (Error)**:
```json
{
  "error": {
    "code": -32601,
    "message": "Method not found: methodName"
  },
  "id": "optional-request-id"
}
```

## Supported Methods

### 1. ping

**Description**: Simple ping method to test connectivity

**Request**:
```json
{
  "method": "ping",
  "id": "1"
}
```

**Response**:
```json
{
  "result": {
    "message": "pong"
  },
  "id": "1"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method":"ping","id":"1"}'
```

---

### 2. echo

**Description**: Echoes back the parameters sent in the request

**Request**:
```json
{
  "method": "echo",
  "params": {
    "message": "Hello, MCP!",
    "timestamp": 1234567890
  },
  "id": "2"
}
```

**Response**:
```json
{
  "result": {
    "message": "Hello, MCP!",
    "timestamp": 1234567890
  },
  "id": "2"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method":"echo","params":{"message":"Hello"},"id":"2"}'
```

---

### 3. getContext

**Description**: Returns context information along with the provided parameters

**Request**:
```json
{
  "method": "getContext",
  "params": {
    "key": "value"
  },
  "id": "3"
}
```

**Response**:
```json
{
  "result": {
    "context": "MCP server context",
    "params": {
      "key": "value"
    }
  },
  "id": "3"
}
```

**Example**:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method":"getContext","params":{"user":"test"},"id":"3"}'
```

## Error Codes

The MCP server follows JSON-RPC 2.0 error code conventions:

| Code | Message | Description |
|------|---------|-------------|
| -32700 | Parse error | Invalid JSON was received by the server |
| -32601 | Method not found | The requested method does not exist |

## Server Configuration

### Timeouts

The server is configured with the following timeouts:

- **Read Timeout**: 15 seconds
- **Write Timeout**: 15 seconds
- **Idle Timeout**: 60 seconds

### Port Configuration

The server port is configured via environment variables. Ensure the `PORT` environment variable is set in your `.env.test`, `.env.prod`, or `.env.local` file.

## Testing

You can test the MCP server using the provided `testRequest.http` file with the REST Client extension in VSCode, or use `curl` commands as shown in the examples above.

### Example Test Sequence

1. Check server health:
```bash
curl http://localhost:8080/health
```

2. Get server info:
```bash
curl http://localhost:8080/info
```

3. Test ping method:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method":"ping","id":"1"}'
```

4. Test echo method:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method":"echo","params":{"test":"data"},"id":"2"}'
```

5. Test getContext method:
```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method":"getContext","params":{"user":"example"},"id":"3"}'
```

## Integration

To integrate the MCP server into your application:

1. Import the server package:
```go
import "main/internal/server"
```

2. Start the server:
```go
server.Start()
```

The server will run as a blocking call, so you may want to run it in a goroutine if you need to perform other operations:

```go
go server.Start()
```

## Future Enhancements

Potential areas for expansion:

- Additional MCP methods for domain-specific functionality
- Authentication and authorization
- WebSocket support for real-time communication
- Request logging and metrics
- Rate limiting
- API versioning

## References

- Model Context Protocol (MCP) specification
- JSON-RPC 2.0 specification
