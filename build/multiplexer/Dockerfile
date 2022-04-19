FROM golang:1.18 as builder

WORKDIR /go/src/multiplexer

COPY ../../go.mod ./go.mod
COPY ../../internal ./internal
COPY ../../cmd/multiplexer/main.go ./main.go

RUN go build && go install

FROM alpine:latest
RUN apk add --no-cache libc6-compat && \
    apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /go/bin/multiplexer /app/multiplexer

CMD ["/app/multiplexer"]
