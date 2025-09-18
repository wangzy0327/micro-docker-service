FROM scratch
ENV NODE_NAME 10.208.16.106
ENV SUBSCRIBE_PATH /home/wzy/micro-docker-service/subscribe/subscribe.txt
ADD scratch1 /
CMD ["/scratch1","10"]
