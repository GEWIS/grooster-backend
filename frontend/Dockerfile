FROM node:22 AS build

WORKDIR /app

COPY ./package.json ./
COPY ./package-lock.json ./

RUN npm install

COPY . ./

RUN npm run build

#Sever
FROM caddy:2-alpine
COPY --from=build /app/dist /srv
COPY ./caddy/Caddyfile /etc/caddy/Caddyfile
EXPOSE 8080
