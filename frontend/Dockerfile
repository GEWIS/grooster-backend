FROM node:22 AS build

WORKDIR /app

COPY ./package.json ./
COPY ./package-lock.json ./

RUN npm install

COPY . ./

RUN npm run build

FROM nginx:alpine AS target
WORKDIR /app
COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY --from=build --chown=nginx /app/dist /app
