package utils

import (
	"fmt"
	"os"
)

func LoadEnvVariables() {
	fmt.Println("the env variable is " + os.Getenv("TEST_ENV"))

}
