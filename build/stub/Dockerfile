FROM golang:1.18 as builder

WORKDIR /go/src/stub

COPY ../../go.mod ./go.mod
COPY ../../internal ./internal
COPY ../../cmd/stub/main.go ./main.go

RUN go build && go install

FROM alpine:latest
RUN apk add --no-cache libc6-compat && \
    apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /go/bin/multiplexer /app/stub

CMD ["/app/stub"]
