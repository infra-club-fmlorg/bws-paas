FROM nginx:latest

COPY {{ .ApplicationPath }} /application
COPY default.conf /etc/nginx/conf.d/default.conf

RUN apt update && apt install -y unzip &&\
    unzip /application &&\
    rm -r /usr/share/nginx/html  &&\
    mv /{{ .ApplicationName }} /usr/share/nginx/html
