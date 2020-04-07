# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.13

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/marypr/guruBackendService

# Changing working directory.
WORKDIR /go/src/github.com/marypr/guruBackendService

# Building application.
RUN go build -o guruBackendService main.go

# Expose a port
EXPOSE 8080

# Run application when container starts.
CMD ["/go/src/github.com/marypr/guruBackendService/guruBackendService"]