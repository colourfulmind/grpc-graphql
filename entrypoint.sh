#!/bin/sh

until pg_isready -h $DB_HOST -p $DB_PORT; do
  echo "Waiting for database at $DB_HOST:$DB_PORT..."
  sleep 2
done

./migrator -migrate up

exec "$@"
