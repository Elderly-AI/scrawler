FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go build ./cmd/main.go

EXPOSE 8090

RUN chmod +x /app/start.sh
CMD ["sh","/app/start.sh"]