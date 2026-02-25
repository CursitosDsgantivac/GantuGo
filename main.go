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
	utils.LoadEnvVariables()

	switch arg {
	case "server":
		httpServer.Server()
	case "chiServer":
		httpServer.CreateChiServer()
	case "mcpServer":
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
