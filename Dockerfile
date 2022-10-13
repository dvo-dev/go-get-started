# syntax=docker/dockerfile:1
# Build image
FROM golang:latest

# Set working directory
WORKDIR /go-get-started

#Populate module caches based on GO dependencies
COPY go.mod .
COPY go.sum .

# Get dependencies
RUN go mod download

# Get everything else
COPY . .

# Unit tests
RUN make test

# Build executable
RUN make build

# Expose port
EXPOSE 8080

CMD [ "./bin/go-get-started/webapp" ]
