#!/bin/zsh

docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 -t meschbach/vault-k8s-unsealer:$1 --push  .
