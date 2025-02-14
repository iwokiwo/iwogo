# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.12-alpine base image
FROM golang:alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
# RUN apk update && apk upgrade && \
#     apk add --no-cache bash git openssh
RUN apk update && apk add --no-cache git

# Add Maintainer Info
LABEL maintainer="iwokiwo <bayuiwo@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# # Copy go mod and sum files
# COPY go.mod go.sum ./

# # Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN mkdir -p storage/branches
RUN mkdir -p storage/store
RUN mkdir -p storage/item
RUN mkdir -p storage/gallery
# RUN touch storage/logs/swoole_http.log
# RUN touch storage/logs/laravel.log
# RUN touch storage/logs/crontab.log
RUN chmod -R 777 storage

RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8090

# Run the executable
CMD ["./main"]