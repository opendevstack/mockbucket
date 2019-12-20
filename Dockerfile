# BUILDING THE APP
FROM golang:1.13-alpine AS builder

COPY go.* /apiserver/
COPY main.go /apiserver/main.go
COPY api /apiserver/api
COPY utils /apiserver/utils

WORKDIR /apiserver
RUN env CGO_ENABLED=0 go build -o apiserver

# DEPLOYING THE APP
FROM nginx:1.17-alpine
RUN apk \
    --update \
    --no-cache \
    --virtual build-dependencies \
    add \
    git \
    apache2-utils \
    git-daemon \
    fcgiwrap \
    spawn-fcgi \
    jq \
    bash \
    curl

COPY www /www
COPY scripts /scripts
COPY conf/default.conf.nginx /etc/nginx/conf.d/default.conf
COPY --from=builder /apiserver/apiserver /apiserver/apiserver

RUN chmod +x /apiserver/apiserver \
    && chmod +x /scripts/* \
    && chgrp -R 0 /var/log/nginx \
    && chmod -R g+rwX /var/log/nginx \
    && chgrp -R 0  /var/cache/nginx \
    && chmod -R g+rwX /var/cache/nginx \
    && chgrp -R 0  /var/run \
    && chmod -R g+rwX /var/run \
    && chgrp -R 0 /etc/nginx \
    && chmod -R g+rwX /etc/nginx \
    && chgrp -R 0 /apiserver \
    && chmod -R g+rwX /apiserver \
    && mkdir /scm \
    && chgrp -R 0 /scm \
    && chmod -R g+rwX /scm \
    && chgrp -R 0 /www \
    && chmod -R g+rwX /www

EXPOSE 8080

USER 1001

ENTRYPOINT ["/scripts/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]