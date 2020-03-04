FROM arm64v8/openjdk:8
MAINTAINER wzy
ADD sum.jar /root/
WORKDIR /root/
CMD java -jar sum.jar 10 abc-edf-hig
