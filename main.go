package main

import (
	"fmt"
	httpServer "main/internal/server/http"
	mcpServer "main/internal/server/mcp"
	"main/internal/utils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments sent: program <arg1> <arg2>...")
		os.Exit(1)
	}
	var arg = os.Args[1]

	switch arg {
	case "server":
		utils.LoadEnvVariables()
		httpServer.Server()
	case "chiServer":
		utils.LoadEnvVariables()
		httpServer.CreateChiServer()
	case "mcpServer":
		utils.LoadEnvVariables()
		mcpServer.StartMCP()
	case "mcpServerHttp":
		mcpServer.StartHttpMCP()
	case "mcpServerBinary":
		mcpServer.StartMCPBinary()
	case "utils":
		utils.LoadEnvVariables()
	default:
		fmt.Println("Invalid argument: " + arg)
		os.Exit(1)
	}
}
