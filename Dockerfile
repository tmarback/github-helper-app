FROM golang:1.26-trixie AS build

WORKDIR /build
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go mod download
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -o bin/server ./cmd/server 

FROM debian:trixie-slim

# Root group for OpenShift compatibility
RUN useradd --uid 1000 --gid root --no-user-group --create-home runner

USER runner
WORKDIR /server

COPY --from=build --chown=root:root --chmod=755 /build/bin/server /opt/server

ENTRYPOINT [ "/opt/server" ]
