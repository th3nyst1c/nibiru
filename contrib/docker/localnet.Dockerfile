FROM golang:1.19 AS builder

WORKDIR /nibiru

COPY go.sum go.mod ./
RUN go mod download
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  --mount=type=cache,target=/nibiru/temp \
  make build

FROM alpine:latest
WORKDIR /root

RUN apk --no-cache add \
  ca-certificates \
  build-base \
  curl \
  jq

COPY --from=builder /nibiru/build/nibid /usr/local/bin/nibid

COPY ./contrib/scripts/configure-localnet.sh ./
RUN chmod +x ./configure-localnet.sh
ARG MNEMONIC
ARG CHAIN_ID
RUN MNEMONIC=${MNEMONIC} CHAIN_ID=${CHAIN_ID} ./configure-localnet.sh

ENTRYPOINT ["nibid"]
CMD [ "start" ]
