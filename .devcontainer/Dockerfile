FROM docker.io/gotify/build:1.25.1-linux-amd64

ENV NVM_DIR /root/.nvm
ENV NODE_VERSION 22

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash \
    && . $NVM_DIR/nvm.sh \
    && nvm install 22 \
    && npm install --global yarn