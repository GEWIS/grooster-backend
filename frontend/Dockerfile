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
EXPOSE 8080
CMD ["caddy", "file-server", "--root", "/srv", "--listen", ":8080"]
