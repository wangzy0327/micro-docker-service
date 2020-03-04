FROM fedora21-base
MAINTAINER wzy
RUN yum -y install java-1.8.0-openjdk && \
yum clean all

