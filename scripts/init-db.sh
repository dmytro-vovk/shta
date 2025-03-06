#!/bin/bash

set -e

if [ -n "${APP_DB_USER}" ] && [ -n "${APP_DB_PASSWORD}" ] && [ -n "${APP_DB_NAME}" ]; then
  echo "Creating database '${APP_DB_NAME}' and user '${APP_DB_USER}'..."

  psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" -d "${POSTGRES_DB:-postgres}" <<-END
    CREATE USER ${APP_DB_USER} WITH PASSWORD '${APP_DB_PASSWORD}';
    CREATE DATABASE ${APP_DB_NAME} OWNER ${APP_DB_USER};

    -- REVOKE CONNECT ON DATABASE ${APP_DB_NAME} FROM PUBLIC;
    -- REVOKE ALL ON SCHEMA public FROM PUBLIC;
    -- REVOKE ALL ON ALL TABLES IN SCHEMA public FROM PUBLIC;
END

  echo "Creating tables in database '${APP_DB_NAME}'..."

  psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" -d "${APP_DB_NAME}" < /docker-entrypoint-initdb.d/init.sql

  echo "Granting access to database '${APP_DB_NAME}' for user '${APP_DB_USER}'..."

  psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" -d "${POSTGRES_DB:-postgres}" <<-END
      GRANT CONNECT ON DATABASE ${APP_DB_NAME} TO ${APP_DB_USER};
      GRANT USAGE ON SCHEMA public TO ${APP_DB_USER};
      GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO ${APP_DB_USER};
END

fi
