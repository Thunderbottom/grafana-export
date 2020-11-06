FROM golang:alpine as builder

RUN adduser \
        --disabled-password \
        --gecos "" \
        --no-create-home \
        --shell "/sbin/nologin" \
        --uid 10001 \
        go

WORKDIR /build

COPY . .

RUN apk add --no-cache make ca-certificates
RUN make dist

FROM scratch

WORKDIR /

COPY --from=builder /build/grafana-export .
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/shadow /etc/shadow
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER go

ENTRYPOINT ["/grafana-export"]
