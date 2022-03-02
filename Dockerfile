FROM golang:alpine

MAINTAINER poembro@126.com 
LABEL version="0.0.1" description="这是一个golang im服务器"

WORKDIR /webser/go_wepapp/golang-im

ENV APP_ENV=local 
ENV GRPC_CONNECT_ADDR=192.168.83.165:50000
ENV GRPC_LOGIC_ADDR=192.168.83.165:50100

RUN mkdir -p /webser/go_wepapp/golang-im/cmd/logic
RUN mkdir -p /webser/go_wepapp/golang-im/cmd/connect
RUN mkdir -p /webser/logs/ 

COPY ./cmd/logic/logic /webser/go_wepapp/golang-im/cmd/logic/
COPY ./cmd/connect/connect /webser/go_wepapp/golang-im/cmd/connect/
COPY ./docker-start.sh /webser/go_wepapp/golang-im/

RUN chmod +rwx /webser/go_wepapp/golang-im/cmd/logic/logic
RUN chmod +rwx /webser/go_wepapp/golang-im/cmd/connect/connect
RUN chmod +x /webser/go_wepapp/golang-im/docker-start.sh

EXPOSE 8090
EXPOSE 7923
EXPOSE 50000
EXPOSE 50100
 
ENTRYPOINT ["sh", "/webser/go_wepapp/golang-im/docker-start.sh"]