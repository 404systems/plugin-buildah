FROM golang AS builder
WORKDIR /build
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache go mod download
COPY main.go ./
COPY internal/ ./internal/
RUN --mount=type=cache,target=/go-cache --mount=type=cache,target=/gomod-cache go build .

FROM quay.io/buildah/stable:v1.42

ENV HOME=/tmp BUILDAH_LAYERS=true BUILDAH_ISOLATION=rootless \
    PLUGIN_REGISTRIES_FILE=/tmp/registries.conf PLUGIN_AUTHS_FILE=/tmp/auths.json \
    CI_WORKSPACE="/workspace" PLUGIN_CONTAINERFILE="Containerfile"

WORKDIR ${CI_WORKSPACE}

CMD [ "/plugin-buildah" ]

COPY --from=builder /build/plugin-buildah /