FROM golang:1.24-alpine AS builder
LABEL authors="Dev Tobby"
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .
CMD ["./main", "-mode=server"]
