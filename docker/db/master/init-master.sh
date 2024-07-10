#!/bin/bash
set -e

echo "///////////////// Starting init-master.sh with user $DB_USER and database $DB_DB"

export PGPASSWORD=$DB_PASSWORD

# Добавление репликационной конфигурации в postgresql.conf
echo "///////////////// Appending replication configuration to postgresql.conf"
echo "listen_addresses = '*'" >> /var/lib/postgresql/data/postgresql.conf
echo "wal_level = replica" >> /var/lib/postgresql/data/postgresql.conf
echo "max_wal_senders = 2" >> /var/lib/postgresql/data/postgresql.conf
echo "max_replication_slots = 2" >> /var/lib/postgresql/data/postgresql.conf
echo "hot_standby = on" >> /var/lib/postgresql/data/postgresql.conf
echo "hot_standby_feedback = on" >> /var/lib/postgresql/data/postgresql.conf

# Проверка и логирование содержимого pg_hba.conf до изменений
echo "///////////////// Contents of pg_hba.conf before modifications:"
cat /var/lib/postgresql/data/pg_hba.conf || echo "pg_hba.conf not found"

# Добавление конфигурации репликации в pg_hba.conf
echo "host replication $DB_USER 127.0.0.1/32 trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host replication $DB_USER db-master trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host replication $DB_USER db-slave trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host replication $DB_USER 0.0.0.0/0 trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host all all 0.0.0.0/0 md5" >> /var/lib/postgresql/data/pg_hba.conf

# Перезапуск PostgreSQL для применения изменений
pg_ctl -D "$PGDATA" -m fast -w stop
pg_ctl -D "$PGDATA" -o "-c listen_addresses='*'" -w start

echo "///////////////// Configured"

# Создание пользователя и базы данных
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" <<-EOSQL
    CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
    ALTER USER $DB_USER WITH SUPERUSER;
    CREATE DATABASE $DB_NAME OWNER $DB_USER;
EOSQL

echo "///////////////// Created user & db"

# Создание таблицы users в базе данных
psql -v ON_ERROR_STOP=1 --username "$DB_USER" --dbname "$DB_NAME" <<-EOSQL
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        birth_date DATE,
        gender CHAR(1),
        interests TEXT,
        city VARCHAR(50),
        username VARCHAR(50),
        password VARCHAR(100)
    );
EOSQL

echo "///////////////// init-master.sh completed"
