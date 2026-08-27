# syntax=docker/dockerfile:1.7

FROM composer:latest AS vendor
WORKDIR /app
COPY src/composer.json src/composer.lock ./
RUN --mount=type=cache,target=/tmp/composer \
    composer install --no-dev --prefer-dist --no-interaction --no-progress --ignore-platform-reqs

FROM php:8.4-fpm AS runtime

WORKDIR /var/www/html

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        git \
        unzip \
        nginx \
        libicu-dev \
        libzip-dev \
        libonig-dev \
        libxml2-dev \
        libssl-dev \
        libpq-dev \
        libmariadb-dev \
        libcurl4-openssl-dev \
        zlib1g-dev \
    && docker-php-ext-configure intl && \
    docker-php-ext-install \
        intl \
        mbstring \
        mysqli \
        pdo_mysql \
        pdo_pgsql \
        zip \
        bcmath \
        curl \
    && docker-php-ext-enable opcache && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /run/nginx

COPY --from=vendor /app/vendor ./vendor
COPY src ./
COPY docker/prod/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY docker/prod/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh && \
    chown -R www-data:www-data writable && \
    rm -f /etc/nginx/sites-enabled/default && \
    mkdir -p /var/log/nginx

ENV CI_ENVIRONMENT=production \
    PHP_OPCACHE_VALIDATE_TIMESTAMPS=0 \
    PHP_OPCACHE_MAX_ACCELERATED_FILES=20000 \
    PHP_OPCACHE_MEMORY_CONSUMPTION=256 \
    PHP_OPCACHE_INTERNED_STRINGS_BUFFER=16

LABEL org.opencontainers.image.authors="Nefla <neflaprojekt@gmail.com>" \
      org.opencontainers.image.vendor="C4ignite by NeflaLabs" \
      org.opencontainers.image.source="https://github.com/neflalabs/c4ignite" \
      org.opencontainers.image.title="C4ignite App" \
      org.opencontainers.image.description="Codeigniter 4 with C4ignite builder"

EXPOSE 8000

CMD ["/entrypoint.sh"]
