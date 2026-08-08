

# Docker Neovim

Una imagen de Neovim actualizada y lista para usar.

## Características

- soporta cada compilación (nocturna, estable, versiones individuales)
- se actualiza automáticamente
- fácilmente extensible

## Cómo usarlo

Puedes usar esta imagen donde quieras con algo como:
```
docker run -it -v `pwd`:/mnt/volume -w=/mnt/volume anatolelucet/neovim:latest
```

También puedes extender esta imagen en un Dockerfile para crear la tuya (posiblemente conteniendo tu configuración personal).

<details>
<summary>Ejemplo de imagen con configuración personal</summary>

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

## Etiquetas

> Nota: las siguientes etiquetas son ejemplos. Puedes encontrar las etiquetas más recientes en [DockerHub](https://hub.docker.com/r/anatolelucet/neovim/tags?ordering=last_updated).

Alpine: `:latest`, `:latest-alpine`, `:0-alpine`, `:0.11-alpine`, `:0.11.3-alpine`, `:nightly-alpine`

Debian Bookworm: `:latest-bookworm`, `:0-bookworm`, `:0.11-bookworm`, `:0.11.3-bookworm`, `:nightly-bookworm`

Debian Bullseye: `:latest-bullseye`, `:0-bullseye`, `:0.11-bullseye`, `:0.11.3-bullseye`, `:nightly-bullseye`
