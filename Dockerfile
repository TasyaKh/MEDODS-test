FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/app ./app
COPY .env* ./

CMD ["/app/app"] 