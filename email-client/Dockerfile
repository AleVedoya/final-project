FROM node:18.17.1-alpine

RUN npm install -g http-server

RUN mkdir app

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

CMD [ "http-server", "dist"]
