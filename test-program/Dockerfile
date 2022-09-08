FROM node:18-buster-slim

COPY application /application

RUN apt update && apt install -y unzip &&\
    unzip /application &&\
    cd application &&\
    npm i &&\
    npm run build

WORKDIR /application

ENTRYPOINT npm start

