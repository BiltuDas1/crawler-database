FROM golang:1.23.2-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go install
RUN go build -o cdb .


FROM alpine:3.20 AS main
COPY --from=builder /app/cdb /bin/
WORKDIR /data

ENV RunInDocker=true
ENV PORT=10000

LABEL maintainer="Biltu Das <billionto@gmail.com>"
LABEL org.opencontainers.image.version="0.0.1-alpha"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source=https://github.com/BiltuDas1/crawler-db
LABEL org.opencontainers.image.documentation=https://github.com/BiltuDas1/crawler-db/wiki

CMD [ "cdb" ]