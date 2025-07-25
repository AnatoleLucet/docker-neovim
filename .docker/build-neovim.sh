#!/bin/sh

set -xeu

echo "Fetching sources for Neovim version ${VERSION}..."
cd /tmp
wget https://github.com/neovim/neovim/archive/refs/tags/${VERSION}.tar.gz
tar -xzf ${VERSION}.tar.gz

echo "Building Neovim version ${VERSION}..."
cd neovim-*
make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr
make install

echo "Cleaning up build artifacts..."
strip /usr/bin/nvim
