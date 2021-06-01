# FIXME: switch to an alpine image
FROM ubuntu:20.04

LABEL maintainer="AnatoleLucet"

ARG BUILD_APT_DEPS="ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip git"
ARG DEBIAN_FRONTEND=noninteractive
ARG TARGET=stable

RUN apt update && apt upgrade -y && \
  apt install -y ${BUILD_APT_DEPS} && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
	git checkout tags/${TARGET} && \
  make CMAKE_BUILD_TYPE=Release && \
  make CMAKE_INSTALL_PREFIX=$HOME/nvim install && \
  cd ~ && \
  echo "\nPATH=\$PATH:\$HOME/nvim/bin" >> ~/.bashrc && \
  rm -rf /tmp/neovim && \
  apt purge -y ${BUILD_APT_DEPS} && apt autoremove -y --purge

WORKDIR /root

CMD ["/root/nvim/bin/nvim"]
