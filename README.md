# Docker NeoVim

An up-to-date, ready-to-use NeoVim image.

## Features

- supports every builds (nightly, stable, individual versions)
- automatically updated
- easily extensible

## How to use

You can use this image anywhere you want with something like:
```
docker run -it -v `pwd`:/mnt/volume --workdir=/mnt/volume anatolelucet/neovim:stable
```

You can also extend this image in a Dockerfile to make your own (possibly containing your personal setup).

## Releases

You can find every releases here: https://hub.docker.com/r/anatolelucet/neovim/tags
