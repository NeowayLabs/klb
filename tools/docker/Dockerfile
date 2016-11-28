FROM neowaylabs/nash
MAINTAINER Luciano A. Borguetti Faustino <lucianoborguetti@gmail.com>

# Set envs
ENV DEBIAN_FRONTEND noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN true

RUN apt-get -qqy update \
 && apt-get -qqy install awscli \
 && apt-get -qqy install nodejs-legacy npm \
 && npm install -g azure-cli \
 && git clone https://github.com/NeowayLabs/klb /root/.nash/lib/klb \
 && git clone https://github.com/NeowayLabs/nashlib /root/.nash/lib/nashlib \
 && echo "#!/usr/bin/env nash\n\nimport \"nashlib/all\"" > /root/.nash/init \
 && apt-get -qqy autoremove \
 && apt-get -qqy clean \
 && rm -rf /var/lib/apt/*
