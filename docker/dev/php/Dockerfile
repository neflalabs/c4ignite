FROM composer:latest AS composer_stage

FROM php:8.4-fpm AS base

ENV COMPOSER_ALLOW_SUPERUSER=1 \
    APP_DIR=/var/www/html

WORKDIR ${APP_DIR}

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bash \
        git \
        unzip \
        libicu-dev \
        libzip-dev \
        libonig-dev \
        libxml2-dev \
        libssl-dev \
        libpq-dev \
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
    && pecl install xdebug && \
    docker-php-ext-enable opcache xdebug && \
    rm -rf /var/lib/apt/lists/*

COPY --from=composer_stage /usr/bin/composer /usr/bin/composer

RUN echo "xdebug.mode=off" > /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini && \
    echo "xdebug.start_with_request=no" >> /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini

ENV PATH="${APP_DIR}/vendor/bin:${PATH}"
