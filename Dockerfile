# Build the manager binary
FROM golang:1.21 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG GIT_COMMIT='not set'
ARG GIT_TAG=development
ENV GIT_COMMIT=$GIT_COMMIT
ENV GIT_TAG=$GIT_TAG
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
RUN CGO_ENABLED=1 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} GO111MODULE=on \
    go build -ldflags="-X 'main.version=$GIT_TAG' -X 'main.commit=$GIT_COMMIT'" ./cmd/ceph-api

FROM gcr.io/distroless/cc-debian12

# TODO: support multi-arch build
# copy shared libraries
COPY --from=builder /usr/lib/aarch64-linux-gnu/librbd.so* /usr/lib/aarch64-linux-gnu/
# COPY --from=builder /usr/lib/x86_64-linux-gnu/librbd.so* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/aarch64-linux-gnu/librados.so* /usr/lib/aarch64-linux-gnu/
# COPY --from=builder /usr/lib/x86_64-linux-gnu/librados.so* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/aarch64-linux-gnu/libceph.so* /usr/lib/aarch64-linux-gnu/
# COPY --from=builder /usr/lib/x86_64-linux-gnu/libceph.so* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/aarch64-linux-gnu/ceph/libceph-common.so* /usr/lib/aarch64-linux-gnu/ceph/
# COPY --from=builder /usr/lib/x86_64-linux-gnu/ceph/libceph-common.so* /usr/lib/x86_64-linux-gnu/ceph/
COPY --from=builder /lib/aarch64-linux-gnu/libfmt.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libboost_thread.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libboost_iostreams.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libblkid.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libudev.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libibverbs.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/librdmacm.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libz.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libbz2.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/liblzma.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libzstd.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libnl-route-3.so* /lib/aarch64-linux-gnu/
COPY --from=builder /lib/aarch64-linux-gnu/libnl-3.so* /lib/aarch64-linux-gnu/
# COPY --from=builder /lib/x86_64-linux-gnu/libfmt.so* /lib/x86_64-linux-gnu/

# copy app bianry
COPY --from=builder /build/ceph-api /bin/ceph-api
WORKDIR /bin

CMD ["ceph-api"]
