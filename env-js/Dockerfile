FROM node:9.8.0-alpine as builder
WORKDIR /tmp/alpaca
ADD package.json .
ADD env.js .
RUN npm install && npm run build

FROM nginx:1.13.9-alpine
COPY --from=builder /tmp/alpaca/env.min.js /usr/share/nginx/html/env.js
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80 443