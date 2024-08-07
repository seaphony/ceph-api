FROM golang:1.21 AS builder
ARG TARGETOS
ARG TARGETARCH
# pacific, quincy, reef	
ENV CEPH_RELEASE="reef" 

RUN echo $TARGETARCH

WORKDIR /build

# install build dependecies
RUN apt-get update && apt-get install -y gcc g++ librbd-dev librados-dev libcephfs-dev linux-headers-generic ceph-common

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# build app
RUN CGO_ENABLED=1 go test ./test -c -o api-test

CMD ["./api-test", "-test.v"]
