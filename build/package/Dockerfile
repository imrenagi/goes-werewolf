FROM alpine:3.6 as alpine
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' appuser

FROM golang:1.12.0 as golang
ENV GOPATH /go
ENV GO111MODULE on
RUN mkdir -p /go/src/github.com/imrenagi/goes-werewolf
WORKDIR /go/src/github.com/imrenagi/goes-werewolf
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make test
RUN make build

FROM scratch
ENTRYPOINT []
WORKDIR /
COPY --from=alpine /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=golang /go/src/github.com/imrenagi/goes-werewolf/goes-werewolf .
CMD ["./goes-werewolf"]
