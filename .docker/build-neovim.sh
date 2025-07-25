#!/bin/sh

set -xeu

echo "Fetching sources for Neovim version ${NEOVIM_VERSION}..."
cd /tmp
wget https://github.com/neovim/neovim/archive/refs/tags/${NEOVIM_VERSION}.tar.gz
tar -xzf ${NEOVIM_VERSION}.tar.gz

echo "Building Neovim version ${NEOVIM_VERSION}..."
cd neovim-*
make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/tmp/nvim/usr
make install

echo "Cleaning up build artifacts..."
strip /tmp/nvim/usr/bin/nvim
