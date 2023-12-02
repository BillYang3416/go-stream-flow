FROM golang:1.21.3 AS build
WORKDIR /app
COPY . .

WORKDIR /app/cmd/go-file-gate

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o go-file-gate-app .

# Using a Debian base for troubleshooting
FROM debian:buster-slim
# Copy the go-file-gate-app binary
COPY --from=build /app/cmd/go-file-gate/go-file-gate-app /go-file-gate-app

# Copy the static files
COPY --from=build /app/internal/static /internal/static

# Copy the configuration file
COPY --from=build /app/config/config.dev.yml /config/config.dev.yml

COPY ./wait-for-it.sh /wait-for-it.sh

# Make the script executable
RUN chmod +x /wait-for-it.sh

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*

CMD ["/wait-for-it.sh", "db:5432", "--", "/go-file-gate-app"]
