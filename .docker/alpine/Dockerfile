FROM alpine AS builder

LABEL maintainer="AnatoleLucet"

ARG NEOVIM_VERSION=stable

COPY ./install-neovim.sh /tmp/

RUN /tmp/install-neovim.sh

CMD ["/usr/bin/nvim"]
