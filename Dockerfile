FROM golang:1.16

COPY . /webhook
RUN cd /webhook&&go env -w GO111MODULE=on\
    &&go env -w GOPROXY=https://goproxy.cn,direct \
    &&go mod tidy
EXPOSE 8080
CMD ["/bin/bash", "/webhook/start.sh"]