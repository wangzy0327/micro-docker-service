FROM ubuntu:16.04
COPY sources.list sources.list
RUN mv /etc/apt/sources.list /etc/apt/sources.list.bak \
    && mv sources.list /etc/apt/
RUN apt-get update && apt-get install -y --no-install-recommends \
    python2.7 \ 
    python-setuptools \
    python-pip \ 
    curl && \
    python -m pip install --upgrade --force pip && \
    pip install --default-timeout=100 flask && \
rm -rf /var/lib/apt/lists/*
