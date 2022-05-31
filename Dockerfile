FROM golang:1.18.2-alpine AS builder

RUN apk --no-cache add ca-certificates libspatialite build-base

COPY go.mod go.sum mada/

WORKDIR mada

RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/mada

FROM alpine:latest

RUN apk --no-cache add ca-certificates libspatialite

RUN ln -s /usr/lib/mod_spatialite.so.7 /usr/lib/mod_spatialite.so

COPY --from=builder /usr/local/bin/mada /usr/local/bin/mada
