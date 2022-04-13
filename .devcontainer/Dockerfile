# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.224.2/containers/go/.devcontainer/base.Dockerfile

ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

COPY requirements.txt .

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    python3 python3-pip \
    bash-completion \
    lcov \
    # install mkdocs
    && pip3 install -r ./requirements.txt
