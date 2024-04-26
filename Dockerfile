FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o moria .

FROM alpine:latest as runtime

RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

COPY --from=build /app/moria /sr/local/bin/moria

ENTRYPOINT [ "moria" ]