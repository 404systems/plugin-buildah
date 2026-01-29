FROM golang:1.25.5 AS builder
WORKDIR /build
ENV GOCACHE=/go-cache GOMODCACHE=/gomod-cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache go mod download
COPY main.go ./
COPY internal/ ./internal/
RUN --mount=type=cache,target=/go-cache --mount=type=cache,target=/gomod-cache go build .

FROM quay.io/buildah/stable:v1.42

ENV BUILDAH_LAYERS=true BUILDAH_ISOLATION="rootless" HOME="/tmp" \
    PLUGIN_REGISTRIES_FILE="/tmp/registries.conf" PLUGIN_AUTHS_FILE="/tmp/auths.json"

WORKDIR /workspace

VOLUME [ "/workspace" ]

CMD [ "/plugin-buildah" ]

COPY --from=builder /build/plugin-buildah /