FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник (путь изменился)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o socks5-proxy ./cmd/socks5-proxy

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/socks5-proxy .
COPY env.properties .

# Открываем порт
EXPOSE 5431

# Запускаем (имя бинарника изменилось)
CMD ["./socks5-proxy"]
