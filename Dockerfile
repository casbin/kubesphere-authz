FROM golang:1.16

COPY ./build /webhook
EXPOSE 8080
CMD cd /webhook && ./main