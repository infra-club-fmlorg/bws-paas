FROM node:18-buster-slim

COPY application.zip /application

RUN apt update && apt install -y unzip &&\
    unzip /application &&\
    cd vite-project &&\
    npm i &&\
    npm run build

WORKDIR /vite-project

ENTRYPOINT [ "npm", "run", "start" ]

