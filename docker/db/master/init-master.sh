#!/bin/bash
set -e

echo "/////////////// Starting init-master.sh with user $POSTGRES_USER and database $POSTGRES_DB"

export PGPASSWORD=$POSTGRES_PASSWORD

cp /docker-entrypoint-initdb.d/postgresql.conf ${PGDATA}/postgresql.conf
cp /docker-entrypoint-initdb.d/pg_hba.conf ${PGDATA}/pg_hba.conf

echo "/////////////// cp ok"

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

echo "/////////////// Start"

# Создание пользователей и баз данных
psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER replicator REPLICATION LOGIN ENCRYPTED PASSWORD 'password';

    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        birth_date DATE,
        city VARCHAR(50)
    );
EOSQL

echo "/////////////// SQL ok"

#cp /docker-entrypoint-initdb.d/pg_hba.conf "/var/lib/postgresql/pg_hba.conf"

# Перезапуск сервера для применения конфигурации
pg_ctl -D "$PGDATA" restart

echo "/////////////// restart ok"

echo "init-master.sh completed"
