# Docker Neovim

An up-to-date, ready-to-use Neovim image.

## Features

- supports every build (nightly, stable, individual versions)
- automatically updated
- easily extensible

## How to use

You can use this image anywhere you want with something like:

```
docker run -it -v `pwd`:/mnt/volume -w=/mnt/volume anatolelucet/neovim:latest
```

You can also extend this image in a Dockerfile to make your own (possibly containing your personal setup).

<details>
<summary>Image with personal setup example</summary>

```dockerfile
FROM anatolelucet/neovim:latest

# copy your personal neovim config
COPY init.lua /root/.config/nvim/init.lua
COPY lua/ /root/.config/nvim/lua/

# install your favorite tools!
RUN apk add git ripgrep

CMD ["/usr/bin/nvim"]
```

</details>

## Tags

> Note: you can find the latest tags on [DockerHub](https://hub.docker.com/r/anatolelucet/neovim/tags?ordering=last_updated).

Alpine: `:latest`, `:latest-alpine`, `:<major>-alpine`, `:<major>.<minor>-alpine`, `:<major>.<minor>.<patch>-alpine`, `:nightly-alpine`

Debian Trixie: `:latest-trixie`, `:<major>-trixie`, `:<major>.<minor>-trixie`, `:<major>.<minor>.<patch>-trixie`, `:nightly-trixie`

Debian Bookworm: `:latest-bookworm`, `:<major>-bookworm`, `:<major>.<minor>-bookworm`, `:<major>.<minor>.<patch>-bookworm`, `:nightly-bookworm`
