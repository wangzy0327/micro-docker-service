FROM scratch
ENV NODE_NAME 10.130.152.232
ENV SUBSCRIBE_PATH /home/wzy/micro-docker-service/subscribe/subscribe.txt
ADD scratch2 /
CMD ["/scratch2","10"]
