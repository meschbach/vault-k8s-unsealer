# Based on https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
############################
# Builder
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/meschbach/vault-k8s-unsealer/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /vault-k8s-unsealer

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /vault-k8s-unsealer /vault-k8s-unsealer
# Run the hello binary.
CMD ["/vault-k8s-unsealer"]
