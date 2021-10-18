FROM debian:latest as controller
WORKDIR /workspace
COPY build/controller .
ENTRYPOINT ["/workspace/controller" ]

#external webhook
# FROM alpin as webhook
# WORKDIR /workspace
# COPY build/webhook .
# COPY config/config config/config

