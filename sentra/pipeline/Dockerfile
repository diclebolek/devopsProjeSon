FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG CMD=ingest
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/${CMD}

FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/app /usr/local/bin/sentra
USER nobody
ENTRYPOINT ["/usr/local/bin/sentra"]
