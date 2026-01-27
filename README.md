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

