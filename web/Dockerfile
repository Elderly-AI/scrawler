FROM nginx:alpine

COPY .nginx .nginx

RUN rm -rf /usr/share/nginx/html/*

COPY /dist /usr/share/nginx/html/
COPY /.nginx/nginx.conf /etc/nginx/nginx.conf

ENTRYPOINT ["nginx", "-g", "daemon off;"]
