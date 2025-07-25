#!/bin/sh

set -xeu

echo "Installing Neovim version ${NEOVIM_VERSION}..."
wget https://github.com/neovim/neovim/releases/download/${NEOVIM_VERSION}/nvim-linux-x86_64.tar.gz -O /tmp/nvim.tar.gz
tar xzf /tmp/nvim.tar.gz -C /usr --strip-components=1

# strip /usr/bin/nvim
