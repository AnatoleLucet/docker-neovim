ARG DEBIAN_VERSION=bookworm

FROM debian:${DEBIAN_VERSION} AS builder

LABEL maintainer="AnatoleLucet"

ARG NEOVIM_VERSION=stable

COPY ./install-neovim.sh /tmp/

RUN apt update && apt install -y wget && /tmp/install-neovim.sh

CMD ["/usr/bin/nvim"]
