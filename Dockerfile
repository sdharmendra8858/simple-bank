# Build Stage
FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env /app/

EXPOSE 3000
CMD [ "/app/main" ]