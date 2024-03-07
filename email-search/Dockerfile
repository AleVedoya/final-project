# Golang image
FROM golang:1.21.1-alpine

# Env vars
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

## PROJECT SETUP

# Create app directory inside container
RUN mkdir app
COPY . /app

# cd /app
WORKDIR /app

# Build
RUN go build -o main .

# Run
CMD ["/app/main"]

# docker build -t trucode/go .
# docker run -p 8080:8080 trucode/go