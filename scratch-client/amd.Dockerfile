FROM scratch
ENV NODE_NAME 10.18.127.2
ENV SUBSCRIBE_PATH /home/wzy/micro-docker-service/subscribe/subscribe.txt
ADD scratch0 /
CMD ["/scratch0","10"]
