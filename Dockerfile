# build stage
FROM golang:1.16 AS builder
WORKDIR /go/src/sad
COPY . .
RUN make install

# final stage
FROM ubuntu:20.04
WORKDIR /go/src/sad
COPY --from=builder /go/src/sad/SADBackend .
COPY --from=builder /go/src/sad/.env .env
RUN apt-get -qq update && apt-get -qq install -y --no-install-recommends ca-certificates curl
EXPOSE 8888
ENTRYPOINT ["./SADBackend"]
