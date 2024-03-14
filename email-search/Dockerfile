# Golang image
FROM golang:1.21.1-alpine

# Env vars
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0
    # GO_ENV=production

## PROJECT SETUP

# Create app directory inside container
RUN mkdir app

# cd /app
WORKDIR /app

COPY go.mod .

# Download all dependencies
RUN go mod download

# Declare custom args
ARG ZINC_SEARCH_SERVER_URL
ENV ZINC_SEARCH_SERVER_URL=$ZINC_SEARCH_SERVER_URL

# Copies all files from local to container
COPY . .

# Build
RUN go build -o main .

# Run
CMD ["/app/main"]
