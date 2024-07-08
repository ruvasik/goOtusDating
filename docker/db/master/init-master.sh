#!/bin/bash
set -e

echo "Starting init-master.sh with DB_USER=${DB_USER}"

# Создание репликационного пользователя
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE ROLE ${DB_USER} WITH REPLICATION PASSWORD '${DB_PASSWORD}' LOGIN;
EOSQL

# Добавление репликационной конфигурации в postgresql.conf
echo "Appending replication configuration to postgresql.conf"
echo "listen_addresses = '*'" >> /var/lib/postgresql/data/postgresql.conf
echo "wal_level = replica" >> /var/lib/postgresql/data/postgresql.conf
echo "max_wal_senders = 2" >> /var/lib/postgresql/data/postgresql.conf
echo "max_replication_slots = 2" >> /var/lib/postgresql/data/postgresql.conf
echo "hot_standby = on" >> /var/lib/postgresql/data/postgresql.conf
echo "hot_standby_feedback = on" >> /var/lib/postgresql/data/postgresql.conf

# Проверка и логирование содержимого pg_hba.conf до изменений
echo "Contents of pg_hba.conf before modifications:"
cat /var/lib/postgresql/data/pg_hba.conf || echo "pg_hba.conf not found"

# Добавление конфигурации репликации в pg_hba.conf
echo "host replication ${DB_USER} 127.0.0.1/32 trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host replication ${DB_USER} db-master trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host replication ${DB_USER} db-slave trust" >> /var/lib/postgresql/data/pg_hba.conf

# Логирование содержимого pg_hba.conf после изменений
echo "Contents of pg_hba.conf after modifications:"
cat /var/lib/postgresql/data/pg_hba.conf || echo "pg_hba.conf not found"

echo "init-master.sh completed"
