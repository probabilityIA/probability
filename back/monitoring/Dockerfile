FROM golang:1.27-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o monitoring-api cmd/main.go

FROM alpine:3.22

RUN apk upgrade --no-cache && apk add --no-cache ca-certificates wget

WORKDIR /app

COPY --from=builder /app/monitoring-api .

EXPOSE 3070

CMD ["./monitoring-api"]
