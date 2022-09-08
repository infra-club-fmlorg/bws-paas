FROM debian:bullseye-slim

COPY {{ .ApplicationPath }} /application

RUN chmod 111 /application

ENTRYPOINT /application
