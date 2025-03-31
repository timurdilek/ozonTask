FROM golang:alpine AS builder

WORKDIR /ozon
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /ozon/main .
#COPY --from=builder /ozon/config ./config/

EXPOSE 8080

CMD ["./main"]