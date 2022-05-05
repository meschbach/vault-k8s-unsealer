# Based on https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
# and https://stackoverflow.com/questions/71372947/is-it-possible-to-manually-build-multi-arch-docker-image-without-docker-buildx
############################
# Builder
############################
FROM --platform=$BUILDPLATFORM golang:1.18 as builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
RUN uname -a
RUN echo $BUILDPLATFORM $TARGETPLATFORM $TARGETARCH $TARGETOS
ENV USER=appuser
ENV UID=9912
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR /app
RUN mkdir -p /app
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags='-w -s -extldflags "-static"' -o vault-k8s-unsealer ./cmd/k8s-watcher

############################
# STEP 2 build a small image
############################
FROM --platform=$TARGETPLATFORM scratch as final
WORKDIR /
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/vault-k8s-unsealer /vault-k8s-unsealer
USER appuser:appuser
CMD ["/vault-k8s-unsealer"]
