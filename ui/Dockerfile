FROM nginx:1.13.9-alpine
COPY ui/build /usr/share/nginx/html
COPY ui/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80 443