# Stage 1: Build the Go application using golang on Bookworm
FROM golang:1.24-bookworm as builder

# Install curl using apt
RUN apt-get update && apt-get install -y curl

WORKDIR /app

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy application source
COPY . .

# Install templ for code generation
RUN go install github.com/a-h/templ/cmd/templ@latest

# Download Tailwind CSS binary and make it executable
RUN curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
    -o /usr/local/bin/tailwindcss && \
    chmod +x /usr/local/bin/tailwindcss

# Generate files using templ and compile Tailwind CSS
RUN templ generate
RUN tailwindcss -i cmd/web/styles/input.css -o cmd/web/assets/css/output.css

# Build the Go application binary
RUN go build -o main cmd/app/main.go


# Stage 2: Create a smaller runtime image using Debian Bookworm-slim
FROM debian:bookworm-slim as runner

# Install minimal certificates if needed
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN addgroup --system myuser && adduser --system --ingroup myuser myuser

# Copy the binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/main

# Change ownership of the binary
RUN chown myuser:myuser /usr/local/bin/main

USER myuser

EXPOSE 8080

CMD ["/usr/local/bin/main"]

