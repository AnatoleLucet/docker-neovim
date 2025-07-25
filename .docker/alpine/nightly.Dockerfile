FROM alpine AS builder

LABEL maintainer="AnatoleLucet"

ARG VERSION=stable
ARG BUILD_DEPS="build-base cmake libtool autoconf automake gettext-dev curl unzip"

COPY ./build-neovim.sh /tmp/

RUN apk update && apk add --no-cache ${BUILD_DEPS} && sh /tmp/build-neovim.sh

FROM alpine
COPY --from=builder /usr/local /usr/local/

# Required shared libraries
COPY --from=builder /lib/ld-musl-*.so.1 /lib/
COPY --from=builder /usr/lib/libgcc_s.so.1 /usr/lib/
COPY --from=builder /usr/lib/libintl.so.8 /usr/lib/

CMD ["/usr/local/bin/nvim"]
