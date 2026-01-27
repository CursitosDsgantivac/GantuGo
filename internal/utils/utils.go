package utils

import (
	"fmt"
	"os"
)

func LoadPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func LoadEnvVariables() {
	fmt.Println("the env variable is " + os.Getenv("TEST_ENV"))

}
