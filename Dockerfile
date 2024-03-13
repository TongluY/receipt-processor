# Build stage
FROM golang:1.22 AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o receipt-processor .

# Final stage
FROM alpine:latest  
WORKDIR /app
RUN apk --no-cache add ca-certificates && adduser -D myuser
COPY --from=builder /build/receipt-processor /app/
USER myuser
CMD ["./receipt-processor"]
