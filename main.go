package main

import (
	"fmt"
	"main/internal/server"
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
		server.Server()
	case "chiServer":
		server.CreateChiServer()
	case "utils":
		utils.LoadEnvVariables()
	default:
		fmt.Println("Invalid argument: " + arg)
		os.Exit(1)
	}
}
