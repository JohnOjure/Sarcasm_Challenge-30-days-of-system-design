FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

#build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o worker-pool .

#run the application
FROM alpine:latest

WORKDIR /root/

#copy only the compiled binary from the builder stage
COPY --from=builder /app/worker-pool .

CMD ["./worker-pool"]

