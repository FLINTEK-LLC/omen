# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY web ./web

# modernc.org/sqlite is pure Go, so CGO_ENABLED=0 works and keeps the final
# image a static binary with no libc dependency.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/omen ./cmd/omen
RUN mkdir -p /out/data

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=builder /out/omen /app/omen
COPY config.sample.yaml /app/config.sample.yaml

# Pre-create data/ owned by the nonroot user (uid/gid 65532) so writes work
# whether this ends up as an anonymous volume or a bind mount. For a bind
# mount to a host directory, chown/chmod that host directory to 65532:65532
# first, or run the container with --user matching your host uid.
COPY --from=builder --chown=65532:65532 /out/data /app/data
VOLUME ["/app/data"]

EXPOSE 8080

ENTRYPOINT ["/app/omen"]
# Runs out of the box against config.sample.yaml's defaults (Shodan
# enrichment disabled). Mount your own config.yaml over /app/config.yaml and
# point -config at it, or set SHODAN_API_KEY, to customize.
CMD ["-config", "/app/config.sample.yaml"]
