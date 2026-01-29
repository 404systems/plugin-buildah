FROM golang AS builder
WORKDIR /build
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache go mod download
COPY main.go ./
COPY internal/ ./internal/
RUN --mount=type=cache,target=/go-cache --mount=type=cache,target=/gomod-cache go build .

FROM debian:trixie-slim

ENV BUILDAH_LAYERS=true
ENV PLUGIN_REGISTRIES_FILE=/registries.conf
ENV PLUGIN_AUTHS_FILE=/auths.json
ENV CI_WORKSPACE="/workspace"
ENV PLUGIN_CONTAINERFILE="Containerfile"

WORKDIR ${CI_WORKSPACE}

CMD [ "/plugin-buildah" ]

ARG APT_CACHE_BUST=1
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked,id=apt-cache \
    --mount=type=cache,target=/var/lib/apt,sharing=locked,id=apt-lib \
    apt update && apt install -y buildah

COPY --from=builder /build/plugin-buildah /