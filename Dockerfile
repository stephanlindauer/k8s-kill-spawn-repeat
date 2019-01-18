FROM golang:1.11-alpine as builder

RUN apk add --no-cache git make

# Configure Go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Install Go Tools
RUN go get -u golang.org/x/lint/golint
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/stephanlindauer/k8s-kill-spawn-repeat
WORKDIR /go/src/github.com/stephanlindauer/k8s-kill-spawn-repeat
RUN CGO_ENABLED=0 make install

FROM alpine:latest
COPY --from=builder /root/go/bin/k8s-kill-spawn-repeat /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/k8s-kill-spawn-repeat"]
