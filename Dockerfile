FROM node:16-alpine AS ui_builder

COPY ui/package.json ui/yarn.lock ui/

WORKDIR ui

RUN yarn install

COPY ui/ .

RUN yarn build

FROM golang:1.18.2-alpine AS builder

RUN apk --no-cache add ca-certificates libspatialite build-base

COPY go.mod go.sum mada/

WORKDIR mada

RUN go mod download

RUN go install github.com/rakyll/statik@latest

COPY . .

COPY --from=ui_builder /ui/build dist

RUN statik -src=dist && go build -o /usr/local/bin/mada

FROM alpine:3.16

RUN apk --no-cache add ca-certificates libspatialite

RUN ln -s /usr/lib/mod_spatialite.so.7 /usr/lib/mod_spatialite.so

COPY --from=builder /usr/local/bin/mada /usr/local/bin/mada

EXPOSE 8010

CMD [ "mada", "ui" ]