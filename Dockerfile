FROM golang:alpine AS builder
RUN mkdir /build/
WORKDIR /build/
COPY . /build/
ENV CGO_ENABLED=0
RUN go get -d -v
RUN go build -o /go/bin/agent-mysql *.go
FROM alpine:latest
COPY --from=builder /go/bin/agent-mysql /agent-mysql
ENTRYPOINT ["/agent-mysql"]
