# BUILDER
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main ./cmd/api/main.go

# RUNNER
FROM alpine:latest
# add user
RUN adduser -D -u 10001 appuser

WORKDIR /app

# certificates
# RUN apk --no-cache add ca-certificates tzdata

# copy binary from builder
COPY --from=builder /app/bin/main .
# copy config from builder
COPY ./config ./config
# copy migrations from builder
COPY --from=builder app/database/migrations ./database/migrations

USER appuser
EXPOSE 8080
CMD ["./main"]
