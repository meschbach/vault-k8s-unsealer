# Based on https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
# and https://stackoverflow.com/questions/71372947/is-it-possible-to-manually-build-multi-arch-docker-image-without-docker-buildx
############################
# Builder
############################
FROM --platform=$BUILDPLATFORM golang:1.17 AS builder
WORKDIR $GOPATH/src/github.com/meschbach/vault-k8s-unsealer/
ADD . .
ARG TARGETOS
ARG TARGETARCH
# Build the binary.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /vault-k8s-unsealer ./cmd/k8s-watcher

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /vault-k8s-unsealer /vault-k8s-unsealer
# Run the hello binary.
CMD ["/vault-k8s-unsealer"]
