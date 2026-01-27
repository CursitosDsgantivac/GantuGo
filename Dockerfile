# syntax=docker/dockerfile:1

FROM golang:1.25.5

WORKDIR /app

COPY go.mod go.sum ./

COPY *.go ./

COPY internal/ ./

# CGO_ENABLED=0 means: compile pure Go only, no C dependencies
# GOOS=linux means: compile for Linux

RUN CGO_ENABLED=0 GOOS=linux go build -o gantuProgram

EXPOSE 8080

CMD ["./gantuProgram", "chiServer"]
