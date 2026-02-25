## Running commands

- make: Run the default build target
    - make ENV=prod: Build or run the application for the production environment
    - make ENV=test: Build or run the application for the test environment
    - make clean: Remove compiled binaries and build artifacts
    - make fmt: Format the source code according to style guidelines
    - make vet: Run static analysis to find potential issues in the code
    - make build: Compile the source code into an executable binary
    - make USE_ENV_FILE=false: Run the application without environment variables from .env files


## Environment files

The program uses environment variables to configure its behavior. The following environment files are used:

- .env.test: The environment variables to use for the test environment
- .env.prod: The environment variables to use for the production environment
- .env.local: The environment variables to use for the local environment


## Environment variables

The program uses environment variables to configure its behavior. The following environment variables are used:

- TEST_ENV: The environment variable to use for the test environment

# Docker commands

- docker-compose up: Run the application in a Docker container
- docker-compose down: Stop the application in a Docker container

# Project folder structure

- internal: Contains the source code of the application
- cmd: Contains the main entry point of the application

# Testing commands

check the file testRequest.http for more information you should use the extension REST Client in VSCode to run the requests

# GitHub Actions

This project includes a GitHub Actions workflow that automatically builds cross-platform binaries for the MCP server.

## Automatic Builds

The workflow builds binaries for:
- **Windows** (amd64)
- **Linux** (amd64, arm64)
- **macOS** (Intel and Apple Silicon)

## Triggering Builds

### Option 1: Create a Release Tag

```bash
# Commit your changes
git add .
git commit -m "Your commit message"
git push

# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

This will automatically:
1. Build binaries for all platforms
2. Create a GitHub release
3. Attach all binaries to the release

### Option 2: Manual Trigger

1. Go to your GitHub repository
2. Click the **Actions** tab
3. Select **Build MCP Server Binaries**
4. Click **Run workflow**
5. Download artifacts from the workflow run

## Download Binaries

After a release is created, users can download the appropriate binary for their platform from the **Releases** page on GitHub.

# MCP Server Setup

This project includes an MCP (Model Context Protocol) server that can be used with Claude Desktop and Claude Code.

## What is the MCP Server?

The MCP server provides a `greet` tool that can be called by Claude to say hi to users. It runs over stdio (standard input/output) and communicates using JSON-RPC.

## Connecting to Claude Desktop / Claude Code

### Prerequisites

1. Download the appropriate binary for your OS from the [Releases](https://github.com/CursitosDsgantivac/GantuGo/releases) page
2. Make the binary executable (Linux/macOS only)

### Windows Setup

1. **Download the binary:**
   - Download `gantuProgram-windows-amd64.exe` from releases
   - Place it somewhere on your system (e.g., `C:\Tools\gantuProgram.exe`)

2. **Configure Claude Desktop:**
   - Open: `%APPDATA%\Claude\claude_desktop_config.json`
   - Add the following configuration:

```json
{
  "mcpServers": {
    "greeter": {
      "command": "C:\\Tools\\gantuProgram.exe",
      "args": ["mcpServerBinary"]
    }
  }
}
```

3. **Restart Claude Desktop** completely

### Windows with WSL

If your binary is in WSL but you're using Claude Desktop on Windows:

1. **Build or place the binary in WSL:**
   ```bash
   # In WSL
   chmod +x gantuProgram-linux-amd64
   ```

2. **Configure Claude Desktop on Windows:**
   - Open: `%APPDATA%\Claude\claude_desktop_config.json`
   - Add:

```json
{
  "mcpServers": {
    "greeter": {
      "command": "wsl.exe",
      "args": [
        "-e",
        "/home/gantua/gantuProyects/GoLang/gantuProgram-linux-amd64",
        "mcpServerBinary"
      ]
    }
  }
}
```

3. **Restart Claude Desktop**

### Linux Setup

1. **Download and prepare the binary:**
   ```bash
   # Download the binary
   wget https://github.com/YOUR_USERNAME/YOUR_REPO/releases/download/v1.0.0/gantuProgram-linux-amd64

   # Make it executable
   chmod +x gantuProgram-linux-amd64

   # Optionally move to a standard location
   sudo mv gantuProgram-linux-amd64 /usr/local/bin/gantuProgram
   ```

2. **Configure Claude Desktop:**
   - Open: `~/.config/Claude/claude_desktop_config.json`
   - Add:

```json
{
  "mcpServers": {
    "greeter": {
      "command": "/usr/local/bin/gantuProgram",
      "args": ["mcpServerBinary"]
    }
  }
}
```

3. **Restart Claude Desktop**

### macOS Setup

1. **Download and prepare the binary:**
   ```bash
   # For Apple Silicon (M1/M2/M3)
   curl -L -o gantuProgram https://github.com/YOUR_USERNAME/YOUR_REPO/releases/download/v1.0.0/gantuProgram-darwin-arm64

   # For Intel Macs
   # curl -L -o gantuProgram https://github.com/YOUR_USERNAME/YOUR_REPO/releases/download/v1.0.0/gantuProgram-darwin-amd64

   # Make it executable
   chmod +x gantuProgram

   # Move to a standard location
   sudo mv gantuProgram /usr/local/bin/
   ```

2. **Configure Claude Desktop:**
   - Open: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Add:

```json
{
  "mcpServers": {
    "greeter": {
      "command": "/usr/local/bin/gantuProgram",
      "args": ["mcpServerBinary"]
    }
  }
}
```

3. **Restart Claude Desktop**

## Testing the MCP Server

After configuration, you can test the MCP server by asking Claude:

> "Use the greet tool to say hi to John"

Claude should respond with a greeting using your MCP server.

## Running Locally for Development

To run the MCP server locally (without Claude Desktop):

```bash
# Build the project
make build

# Run the MCP server
./gantuProgram mcpServerBinary
```

The server will wait for JSON-RPC messages on stdin.

## Troubleshooting

- **Binary not found:** Ensure the path in the config file is absolute and correct
- **Permission denied (Linux/macOS):** Run `chmod +x gantuProgram`
- **Changes not applied:** Make sure to fully restart Claude Desktop
- **WSL issues:** Verify the WSL path is correct and the binary has execute permissions

