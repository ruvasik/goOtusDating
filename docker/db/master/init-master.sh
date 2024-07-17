#!/bin/bash
set -e

echo "Starting init-master.sh with user $POSTGRES_USER and database $POSTGRES_DB"

export PGPASSWORD=$POSTGRES_PASSWORD

# Старт PostgreSQL
if pg_ctl -D "$PGDATA" status; then
    echo "PostgreSQL is already running"
else
    echo "Starting PostgreSQL"
    pg_ctl -D "$PGDATA" start
    sleep 5
fi

# Ожидание запуска сервера
sleep 5

# Создание пользователей и баз данных
psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
CREATE USER replicator REPLICATION LOGIN ENCRYPTED PASSWORD 'password';
EOSQL

# Перезапуск сервера для применения конфигурации
pg_ctl -D "$PGDATA" restart

echo "init-master.sh completed"
