FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/sint-new-api ./cmd/api

FROM alpine:3.21

RUN adduser -D -H -u 10001 appuser
WORKDIR /app
COPY --from=build /out/sint-new-api /app/sint-new-api
COPY migrations /app/migrations
USER appuser
EXPOSE 10000
ENTRYPOINT ["/app/sint-new-api"]
