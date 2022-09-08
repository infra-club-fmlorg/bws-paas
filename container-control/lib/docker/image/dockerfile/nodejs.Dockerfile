FROM node:18-buster-slim

COPY {{ .ApplicationPath }} /application

RUN apt update && apt install -y unzip &&\
    unzip /application &&\
    cd {{ .ApplicationName }} &&\
    npm i &&\
    npm run build

WORKDIR /{{ .ApplicationName }}

ENTRYPOINT [ "npm", "run", "start" ]

