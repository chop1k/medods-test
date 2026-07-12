# syntax=docker/dockerfile:1

# ---- Build stage -----------------------------------------------------------
FROM golang:1.26.4-alpine AS builder

WORKDIR /src

# Cached separately from source so dependency downloads aren't repeated on
# every code change.
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/trackmytasks ./cmd/track-my-tasks

# ---- Runtime stage ----------------------------------------------------------
FROM alpine AS runtime

WORKDIR /app
COPY --from=builder /out/trackmytasks /app/trackmytasks

COPY migrations /app/migrations

USER nobody:nogroup

# All server/database options are configurable via flags or env vars - see
# internal/config. Example:
#   docker run app serve --http-port=8080
#   docker run -e DB_HOST=postgres app migrate
#   docker run -e DB_HOST=postgres app migrate-and-serve
ENTRYPOINT ["/app/trackmytasks"]
CMD ["serve"]