#!/bin/bash
#wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

HOST="${REDIS_HOST:-redis}"
PORT="${REDIS_PORT:-6379}"
PASSWORD="${REDIS_DB_PASSWORD}"

until nc -z "$host" 5432;
do
    >&2 echo "Postgres is unavailable - sleeping"
    sleep 1
done

>&2 echo "Postgres is up - executing command"

until nc -z "$HOST" "$PORT";
do
    >&2 echo "Redis is unavailable - sleeping"
    sleep 1
done

>&2 echo "Redis is up - executing command"

exec $cmd