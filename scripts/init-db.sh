#!/bin/bash

set -e

if [ -n "${APP_DB_USER}" ] && [ -n "${APP_DB_PASSWORD}" ] && [ -n "${APP_DB_NAME}" ]; then
  echo "Creating database '${APP_DB_NAME}' and user '${APP_DB_USER}'..."

  psql \
    -v ON_ERROR_STOP=1 \
    --username "${POSTGRES_USER}" \
    -d "${POSTGRES_DB:-postgres}" <<-END
    CREATE USER ${APP_DB_USER} WITH PASSWORD '${APP_DB_PASSWORD}';
    CREATE DATABASE ${APP_DB_NAME} OWNER ${APP_DB_USER};
    GRANT ALL PRIVILEGES ON DATABASE ${APP_DB_NAME} TO ${APP_DB_USER};
END

  echo "Creating tables in database '${APP_DB_NAME}'..."

  psql \
    -v ON_ERROR_STOP=1 \
    --username "${APP_DB_USER}" \
    -d "${APP_DB_NAME}" \
    < /docker-entrypoint-initdb.d/init.sql
fi
