# syntax=docker/dockerfile:1.7

FROM composer:2 AS vendor
WORKDIR /app
COPY src/composer.json src/composer.lock ./
RUN --mount=type=cache,target=/tmp/composer \
    composer install --no-dev --prefer-dist --no-interaction --no-progress --ignore-platform-reqs

FROM php:8.2-fpm-alpine AS runtime

WORKDIR /var/www/html

RUN apk add --no-cache bash git nginx icu-data-full icu-dev libzip-dev oniguruma-dev libxml2-dev curl-dev openssl-dev zlib-dev postgresql-dev mariadb-connector-c-dev && \
    docker-php-ext-install intl mbstring mysqli pdo_mysql pdo_pgsql zip bcmath && \
    mkdir -p /run/nginx

COPY --from=vendor /app/vendor ./vendor
COPY src ./
COPY docker/prod/nginx/default.conf /etc/nginx/http.d/default.conf
COPY docker/prod/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh && \
    chown -R www-data:www-data writable && \
    mkdir -p /var/log/nginx

ENV CI_ENVIRONMENT=production \
    PHP_OPCACHE_VALIDATE_TIMESTAMPS=0 \
    PHP_OPCACHE_MAX_ACCELERATED_FILES=20000 \
    PHP_OPCACHE_MEMORY_CONSUMPTION=256 \
    PHP_OPCACHE_INTERNED_STRINGS_BUFFER=16

EXPOSE 8080

CMD ["/entrypoint.sh"]
