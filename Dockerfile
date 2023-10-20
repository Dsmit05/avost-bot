ARG PATH=github.com/Dsmit05/avost-bot

FROM golang:1.20.4-alpine3.17 AS builder

RUN apk add --no-cache ca-certificates git make

WORKDIR /home/${PATH}

COPY . .

RUN make build

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /home/${PATH}/avost-bot .
COPY --from=builder /home/${PATH}/data data/

VOLUME ["/data"]

CMD ["./avost-bot"]