FROM fedora21-base
MAINTAINER wzy
RUN yum install -y python-devel && \
    yum install -y python-setuptools && \
    yum install -y curl  && \
    yum install -y python-pip && \
    python -m pip install --upgrade --force pip && \
    pip install --default-timeout=100 flask && \
yum clean all
