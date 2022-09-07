FROM debian:bullseye-slim

COPY {{.ApplicationPath}} /application

ENTRYPOINT /application
