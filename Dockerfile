FROM golang:1.26-alpine@sha256:a2012b4a1d1f9f7e9f44e0d39b805a901e856b1dd80b4e6ee4338db547050308 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY web ./web
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/carva ./cmd/carva \
    && mkdir -p /out/data

FROM scratch
WORKDIR /app
COPY --from=build --chown=10001:10001 /out/carva /app/carva
COPY --from=build --chown=10001:10001 /out/data /data
USER 10001:10001
ENV CARVA_DB_PATH=/data/carva.db
EXPOSE 8080
ENTRYPOINT ["/app/carva"]
