FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY . .

# Get Module file
COPY go.mod .
RUN go mod download 

# Build the application
RUN go build -o proxy-server .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/proxy-server .

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["/dist/proxy-server"]
