FROM node:22 AS build

WORKDIR /app

COPY ./package.json ./
COPY ./package-lock.json ./

RUN npm install

COPY . ./

RUN npm run build
