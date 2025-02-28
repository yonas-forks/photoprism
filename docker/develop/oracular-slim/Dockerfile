#### Base Image: Ubuntu 24.10 (Oracular Oriole)
FROM ubuntu:oracular

# Copyright © 2018 - 2025 PhotoPrism UG. All rights reserved.
#
# Questions? Email us at hello@photoprism.app or visit our website to learn
# more about our team, products and services: https://www.photoprism.app/

# Add Open Container Initiative (OCI) annotations.
# See: https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="PhotoPrism® Base Image (Ubuntu 24.10)"
LABEL org.opencontainers.image.description="Ubuntu 24.10 (Oracular Oriole)"
LABEL org.opencontainers.image.url="https://hub.docker.com/repository/docker/photoprism/develop"
LABEL org.opencontainers.image.source="https://github.com/photoprism/photoprism"
LABEL org.opencontainers.image.documentation="https://docs.photoprism.app/developer-guide/setup/"
LABEL org.opencontainers.image.authors="PhotoPrism UG <hello@photoprism.app>"
LABEL org.opencontainers.image.vendor="PhotoPrism UG"

# Declare build parameters.
ARG TARGETARCH
ARG BUILD_TAG

# Set environment variables, see https://docs.photoprism.app/getting-started/config-options/.
ENV PHOTOPRISM_ARCH=$TARGETARCH \
    DOCKER_TAG=$BUILD_TAG \
    DOCKER_ENV="prod" \
    PS1="\u@$BUILD_TAG:\w\$ " \
    PATH="/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin:/scripts:/opt/photoprism/bin" \
    LD_LIBRARY_PATH="/usr/local/lib:/usr/lib" \
    TMPDIR="/tmp" \
    DEBIAN_FRONTEND="noninteractive" \
    TF_CPP_MIN_LOG_LEVEL="2" \
    MALLOC_ARENA_MAX="4" \
    PROG="photoprism" \
    S6_KEEP_ENV=1 \
    S6_LOGGING=0

# Copy scripts and package sources config.
COPY --chown=root:root --chmod=755 /scripts/dist/ /scripts/

# Update base image and add dependencies.
RUN echo 'APT::Acquire::Retries "3";' > /etc/apt/apt.conf.d/80retries && \
    echo 'APT::Install-Recommends "false";' > /etc/apt/apt.conf.d/80recommends && \
    echo 'APT::Install-Suggests "false";' > /etc/apt/apt.conf.d/80suggests && \
    echo 'APT::Get::Assume-Yes "true";' > /etc/apt/apt.conf.d/80forceyes && \
    echo 'APT::Get::Fix-Missing "true";' > /etc/apt/apt.conf.d/80fixmissing && \
    echo 'force-confold' > /etc/dpkg/dpkg.cfg.d/force-confold && \
    apt-get update && apt-get -qq dist-upgrade && \
    apt-get -qq install \
        libc6 ca-certificates bash sudo nano tzdata gpg make zip unzip wget curl rsync \
        xz-utils avahi-utils jq lsof lshw libebml5 libgav1-bin libatomic1 exiftool sqlite3 \
        ffmpeg imagemagick libvips-dev rawtherapee libjxl-dev libjxl-tools libffmpeg-nvenc-dev librav1e-dev \
        libswscale-dev libavfilter-extra libavformat-extra libavcodec-extra x264 x265 libde265-dev libaom-dev \
        libvpx-dev libwebm-dev libjpeg-dev libmatroska-dev libdvdread-dev libdav1d-dev libsharpyuv0 \
        iputils-ping dnsutils \
    && \
    /scripts/install-mariadb.sh mariadb-client && \
    /scripts/install-darktable.sh && \
    /scripts/install-libheif.sh && \
    echo 'alias ll="ls -alh"' >> /etc/skel/.bashrc && \
    echo 'export PS1="\u@$DOCKER_TAG:\w\$ "' >> /etc/skel/.bashrc && \
    echo "ALL ALL=(ALL) NOPASSWD:SETENV: /scripts/entrypoint-init.sh" >> /etc/sudoers.d/init && \
    /scripts/install-dircolors.sh && \
    cp /etc/skel/.bashrc /root/.bashrc && \
    /scripts/create-users.sh && \
    install -d -m 0777 -o 1000 -g 1000 \
        /photoprism/originals \
        /photoprism/import \
        /photoprism/storage \
        /photoprism/storage/sidecar \
        /photoprism/storage/albums \
        /photoprism/storage/backups \
        /photoprism/storage/config \
        /photoprism/storage/cache && \
    /scripts/install-s6.sh && \
    /scripts/cleanup.sh

# Set default working directory.
WORKDIR /photoprism

# Expose default HTTP and HTTPS ports.
EXPOSE 2342 2442 2443

# Set default entrypoint and command.
ENTRYPOINT ["/init"]
CMD ["/scripts/cmd.sh", "tail", "-f", "/dev/null"]
