# Stage 1: Build
FROM golang:1.23.4-alpine3.21 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o orchestrator ./cmd/orchestrator

# Stage 2: Image
FROM alpine:3.21

WORKDIR /app

COPY --from=build /app/orchestrator /app/orchestrator

EXPOSE 8080

ENTRYPOINT ["/app/orchestrator"]
