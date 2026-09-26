# Build from apps/website:
#   docker build -t website .
# The site is arch-independent, so build it once on the builder's own platform.
FROM --platform=$BUILDPLATFORM debian:bookworm-slim AS build
ARG HUGO_VERSION=0.166.0
ARG BUILDARCH
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl git libstdc++6 \
    && curl -fsSL "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-${BUILDARCH}.tar.gz" \
       | tar -xz -C /usr/local/bin hugo \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY . .
RUN hugo --gc --minify --destination /public

FROM nginxinc/nginx-unprivileged:stable-alpine
COPY --from=build /public /usr/share/nginx/html
EXPOSE 8080
