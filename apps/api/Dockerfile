FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/pawsear-api ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/pawsear-api /usr/local/bin/pawsear-api
COPY db/migrations ./db/migrations

EXPOSE 8080
CMD ["pawsear-api"]
