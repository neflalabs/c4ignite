#!/usr/bin/env sh
set -eu

# ensure runtime directories exist
mkdir -p /run/nginx /run/php

php-fpm --allow-to-run-as-root --daemonize

exec nginx -g "daemon off;"
