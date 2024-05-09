FROM ubuntu:latest
LABEL authors="emgus"

ENTRYPOINT ["top", "-b"]