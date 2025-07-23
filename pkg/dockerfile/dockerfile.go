package dockerfile

import "fmt"

func Generate(buildType, baseImage, version string) string {
	switch buildType {
	case "nightly":
		return generateNightly(baseImage, version)
	default:
		return generatePackage(baseImage, version)
	}
}

func generateNightly(baseImage, version string) string {
	switch baseImage {
	case "alpine":
		return fmt.Sprintf(`FROM alpine AS builder

LABEL maintainer="AnatoleLucet"

ARG BUILD_DEPS="autoconf automake cmake curl g++ gettext gettext-dev git libtool make ninja openssl pkgconfig unzip binutils wget"
ARG VERSION=%s

RUN apk add --no-cache ${BUILD_DEPS} && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
  git fetch --all --tags -f && \
  git checkout ${VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM alpine
COPY --from=builder /usr/local /usr/local/
# Required shared libraries
COPY --from=builder /lib/ld-musl-*.so.1 /lib/
COPY --from=builder /usr/lib/libgcc_s.so.1 /usr/lib/
COPY --from=builder /usr/lib/libintl.so.8 /usr/lib/

CMD ["/usr/local/bin/nvim"]
`, version)

	case "bookworm":
		return fmt.Sprintf(`FROM debian:bookworm AS builder

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive
ARG BUILD_DEPS="ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip git binutils wget"
ARG VERSION=%s

RUN apt update && apt upgrade -y && \
  apt install -y ${BUILD_DEPS} && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
  git fetch --all --tags -f && \
  git checkout ${VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM debian:bookworm
COPY --from=builder /usr/local /usr/local/

CMD ["/usr/local/bin/nvim"]
`, version)

	case "bullseye":
		return fmt.Sprintf(`FROM debian:bullseye AS builder

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive
ARG BUILD_DEPS="ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip git binutils wget"
ARG VERSION=%s

RUN apt update && apt upgrade -y && \
  apt install -y ${BUILD_DEPS} && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
  git fetch --all --tags -f && \
  git checkout ${VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM debian:bullseye
COPY --from=builder /usr/local /usr/local/

CMD ["/usr/local/bin/nvim"]
`, version)

	default:
		return generateNightly("alpine", version)
	}
}

func generatePackage(baseImage, version string) string {
	switch baseImage {
	case "alpine":
		return `FROM alpine

LABEL maintainer="AnatoleLucet"

RUN apk add --no-cache neovim

CMD ["/usr/bin/nvim"]
`

	case "bookworm":
		return `FROM debian:bookworm

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && apt upgrade -y && \
  apt install -y neovim && \
  apt clean && \
  rm -rf /var/lib/apt/lists/*

CMD ["/usr/bin/nvim"]
`

	case "bullseye":
		return `FROM debian:bullseye

LABEL maintainer="AnatoleLucet"

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && apt upgrade -y && \
  apt install -y neovim && \
  apt clean && \
  rm -rf /var/lib/apt/lists/*

CMD ["/usr/bin/nvim"]
`

	default:
		return generatePackage("alpine", version)
	}
}