###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    upx-ucl

WORKDIR /build
COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 go build -o ./bin/gsrv \
    -ldflags='-w -s -extldflags "-static"' \
    ./server

###########
# release #
###########

FROM gcr.io/distroless/static-debian11:latest AS release

COPY --from=builder /build/bin/gsrv /bin/
WORKDIR /workdir

# sample
COPY dummy.txt dummy.txt

CMD ["/bin/gsrv", "."]
