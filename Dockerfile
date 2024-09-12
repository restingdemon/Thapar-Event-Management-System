# # syntax=docker/dockerfile:1

# FROM golang:1.23.0

# # Set destination for COPY
# WORKDIR /app

# # Download Go modules
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the source code. Note the slash at the end, as explained in
# # https://docs.docker.com/engine/reference/builder/#copy
# COPY *.go ./

# # Build
# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# # Optional:
# # To bind to a TCP port, runtime parameters must be supplied to the docker command.
# # But we can document in the Dockerfile what ports
# # the application is going to listen on by default.
# # https://docs.docker.com/engine/reference/builder/#expose
# EXPOSE 5112

# # Run
# CMD ["/docker-gs-ping"]


# Specify the base image
FROM golang:1.19-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the necessary Go modules
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go application
EXPOSE 5112

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping


# Optional:
# Run the application
CMD ["/docker-gs-ping"]
