
# Stage 1: Build the Go application
FROM golang:latest as go-builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o who-dat .

# Stage 2: Setup the final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=go-builder /go/src/app/who-dat .

# Set the entrypoint to the Go binary
ENTRYPOINT ["./who-dat"]
