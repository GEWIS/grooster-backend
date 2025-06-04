FROM ubuntu:latest
LABEL authors="sven"

ENTRYPOINT ["top", "-b"]