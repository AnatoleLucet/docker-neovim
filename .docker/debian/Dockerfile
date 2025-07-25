ARG DISTRO_VERSION=bookworm
ARG DEBIAN_FRONTEND=noninteractive

FROM debian:${DISTRO_VERSION} AS builder

LABEL maintainer="AnatoleLucet"

ARG NEOVIM_VERSION=stable
ARG BUILD_DEPS="ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip curl wget doxygen"

COPY ./build-neovim.sh /tmp/

RUN apt update && apt install -y ${BUILD_DEPS} && sh /tmp/build-neovim.sh

FROM debian:${DISTRO_VERSION}
COPY --from=builder /tmp/nvim/usr/ /usr/

CMD ["/usr/bin/nvim"]
