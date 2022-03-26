FROM node:16-alpine as build
WORKDIR /app
COPY . .
RUN npm i
RUN npm run build

FROM nginx:stable-alpine
EXPOSE 3000
COPY --from=build /app/build /usr/share/nginx/html
COPY --from=build /app/nginx/default.conf /etc/nginx/conf.d/default.conf