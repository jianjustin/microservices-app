FROM golang:1.20-alpine AS builder

# Set destination for COPY
ENV GOROOT /usr/local/go
ADD . /go/src/calculate-service
WORKDIR /go/src/calculate-service
#WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN GOOS=linux go build -o /server

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 50051
EXPOSE 8080

# Run
CMD ["/server"]