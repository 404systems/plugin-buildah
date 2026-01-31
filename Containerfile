FROM golang:1.25.5 AS builder
WORKDIR /build
ENV GOCACHE=/go-cache GOMODCACHE=/gomod-cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache go mod download
COPY main.go ./
COPY internal/ ./internal/
RUN --mount=type=cache,target=/go-cache --mount=type=cache,target=/gomod-cache go build .

FROM quay.io/buildah/stable:v1.42

LABEL org.opencontainers.image.source = "https://github.com/404systems/plugin-buildah"

ENV BUILDAH_LAYERS=true BUILDAH_ISOLATION="rootless" HOME="/tmp" \
    REGISTRIES_FILE="/tmp/registries.conf" AUTHS_FILE="/tmp/auths.json"

WORKDIR /workspace

VOLUME [ "/workspace" ]

USER 1000:1000

CMD [ "/plugin-buildah" ]

COPY --from=builder /build/plugin-buildah /