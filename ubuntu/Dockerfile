FROM ubuntu AS builder

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive
ARG BUILD_DEPS="ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip git binutils wget"
ARG VERSION=stable

RUN apt update && apt upgrade -y && \
  apt install -y ${BUILD_DEPS} && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
  git fetch --all --tags -f && \
  git checkout ${VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM ubuntu
COPY --from=builder /usr/local /usr/local/

CMD ["/usr/local/bin/nvim"]
