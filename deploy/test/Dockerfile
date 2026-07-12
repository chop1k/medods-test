# syntax=docker/dockerfile:1

# Runs the e2e suite against an in-process Gin router (tests/e2e spins the
# server up via httptest), so no separate app container needs to be running.
FROM golang:1.26.4-alpine

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

ENTRYPOINT ["go", "test", "-v", "./tests/e2e/..."]