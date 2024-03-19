# Set the base build image
ARG BASE_BUILD_IMAGE=golang:1.22-alpine
ARG BASE_IMAGE=alpine

# Build stage
FROM ${BASE_BUILD_IMAGE} as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copy all files from current directory to /app in the image
COPY . /app
RUN go build -o main

# Production stage
FROM ${BASE_IMAGE}
WORKDIR /app

RUN mkdir "file"

# Copy the built main file from build stage to production stage
COPY --from=build /app/main /app/main

CMD ["./main"]