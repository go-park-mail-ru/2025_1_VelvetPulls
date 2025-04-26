FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o csat ./services/csat_service/cmd/csat.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/csat .

ENTRYPOINT ["./csat"]
