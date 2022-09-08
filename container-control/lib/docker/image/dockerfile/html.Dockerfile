FROM nginx:latest

COPY {{ .ApplicationPath }} /application

RUN apt update && apt install -y unzip &&\
    unzip /application &&\
    rm -r /usr/share/nginx/html  &&\
    mv /application /usr/share/nginx/html
