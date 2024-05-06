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
RUN apk add --no-cache --update libpng-dev libjpeg-turbo-dev giflib-dev tiff-dev autoconf automake make gcc g++ wget

RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.2.0.tar.gz && \
tar -xvzf libwebp-1.2.0.tar.gz && \
mv libwebp-1.2.0 libwebp && \
rm libwebp-1.2.0.tar.gz && \
cd /libwebp && \
./configure && \
make && \
make install && \
rm -rf libwebp

COPY --from=builder /app/moria /usr/local/bin/moria

ENV SKIP_DOWNLOAD=true
ENV VENDOR_PATH=/usr/local/bin
ENTRYPOINT [ "moria" ]