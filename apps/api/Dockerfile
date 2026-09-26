# Build from the repo root:
#   docker build -f apps/api/Dockerfile -t 415-api .
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY pkg ./pkg
COPY apps/api ./apps/api

ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown
ARG AUTHOR=unknown
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -ldflags "-X 'github.com/ofgrenudo/415/pkg/version.Version=${VERSION}' \
              -X 'github.com/ofgrenudo/415/pkg/version.Commit=${COMMIT}' \
              -X 'github.com/ofgrenudo/415/pkg/version.Date=${DATE}' \
              -X 'github.com/ofgrenudo/415/pkg/version.Author=${AUTHOR}'" \
    -o /out/api ./apps/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
CMD ["serve"]
