# Builder
FROM golang:alpine as builder

RUN apk --no-cache add gcc musl-dev
WORKDIR /go/src/node

COPY . .
RUN go build -mod=readonly .

# Result
FROM alpine

WORKDIR /app
COPY --from=builder /go/src/node/perun-polkadot-demo .

ENTRYPOINT [ "/app/perun-polkadot-demo" ]
