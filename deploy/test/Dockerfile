# syntax=docker/dockerfile:1

# ---- Build stage -----------------------------------------------------------
FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/trackmytasks ./cmd/track-my-tasks
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go test -c -trimpath -ldflags="-s -w" -o /out/e2e ./tests/e2e/...

# ---- Runtime stage ----------------------------------------------------------
FROM alpine AS runtime

WORKDIR /app

COPY --from=builder /out/trackmytasks /app/trackmytasks
COPY --from=builder /out/e2e /app/e2e
COPY --from=builder /src/migrations /app/migrations
COPY --from=builder /src/seeds /app/seeds

USER nobody:nogroup

ENTRYPOINT ["sh", "-c", "/app/trackmytasks migrate && /app/trackmytasks seed && /app/e2e"]