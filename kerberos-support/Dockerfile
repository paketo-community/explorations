FROM ubuntu:bionic

RUN echo "debconf debconf/frontend select noninteractive" | debconf-set-selections \
  && export DEBIAN_FRONTEND=noninteractive \
  && apt-get -y update \
  && apt-get -y upgrade \
  && apt-get -y install locales \
  && locale-gen en_US.UTF-8 \
  && update-locale LANG=en_US.UTF-8 LANGUAGE=en_US.UTF-8 LC_ALL=en_US.UTF-8 \
  && apt-get -y install \
    bison \
    curl \
    gcc \
    make \
  && rm -rf /var/lib/apt/lists/* /tmp/* /etc/apt/preferences
