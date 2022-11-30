FROM paketobuildpacks/build:1.1.30-base

RUN mkdir /watchexec
WORKDIR /watchexec

ADD "https://github.com/watchexec/watchexec/releases/download/cli-v1.17.1/watchexec-1.17.1-x86_64-unknown-linux-gnu.tar.xz" .

RUN tar xvf watchexec-1.17.1-x86_64-unknown-linux-gnu.tar.xz
RUN cp watchexec-1.17.1-x86_64-unknown-linux-gnu/watchexec /usr/local/bin/watchexec
RUN chmod +x /usr/local/bin/watchexec
