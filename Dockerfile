FROM golang:alpine as builder
WORKDIR /app
# Copy all files into the image
COPY . .
# Run go mod
RUN go mod tidy
# Build Go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /build ./cmd

FROM debian AS runner
WORKDIR /
# COPY executable file from previous builder stage
COPY --from=builder /build /build
COPY .env .env
# Expose ports
EXPOSE 8081
# Run Go program, just like locally
ENTRYPOINT ["/build"]