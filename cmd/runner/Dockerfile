FROM ubuntu:20.10 as build-base

RUN mkdir -p /build \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
        golang git \
        libgstreamer1.0 libgstreamer1.0-dev \
        libgstreamer-plugins-bad1.0-dev libgstreamer-plugins-base1.0-dev

COPY go.mod /build/go.mod
RUN cd /build && go mod download

##

FROM build-base as plugin-build

COPY gst/plugins/go.mod /build/plugins/go.mod
RUN cd /build/plugins && go mod download

COPY gst/plugins/minio /build/plugins/minio
RUN cd /build/plugins/minio && go build -o libgstminio.so -buildmode c-shared .

##

FROM build-base as runner-build

COPY apis /build/apis
COPY pkg /build/pkg
COPY cmd/runner /build/runner
RUN cd /build/runner && go build -o runner .

##

FROM ubuntu:20.10

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
        libgstreamer1.0 gstreamer1.0-plugins-base \
        gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly \
        gstreamer1.0-libav gstreamer1.0-tools

COPY --from=plugin-build /build/plugins/minio/libgstminio.so /usr/lib/x86_64-linux-gnu/gstreamer-1.0/libgstminio.so
COPY --from=runner-build /build/runner/runner /usr/local/bin/runner

CMD /usr/local/bin/runner
