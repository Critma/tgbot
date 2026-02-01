# Build
FROM golang:alpine AS builder
WORKDIR /app
ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./build/main ./cmd/bot/*.go

#Run
FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/configs .
COPY --from=builder /app/build/main .
CMD ["./main"]