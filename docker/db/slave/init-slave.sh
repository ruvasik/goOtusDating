#!/bin/bash
set -e

echo "+++++++++ Starting init-slave.sh"

pg_ctl -D "$PGDATA" -w stop
rm -rf "$PGDATA"/*

echo "+++++++++ Clean ok"

# Ожидание доступности мастера
until psql -h "db_master" -U "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c '\q'; do
  >&2 echo "Master is unavailable - sleeping"
  sleep 5
done

echo "+++++++++ Slave connected to master"

ls -la "$PGDATA"

# Выполнение pg_basebackup
PGPASSWORD=password pg_basebackup -h db_master -D ${PGDATA} -U replicator -P -R

echo "+++++++++ Backup ok"

# Применение конфигурационных файлов
cp /docker-entrypoint-initdb.d/postgresql.conf ${PGDATA}/postgresql.conf
#cp /docker-entrypoint-initdb.d/pg_hba.conf ${PGDATA}/pg_hba.conf

echo "+++++++++ cp ok"

# Старт PostgreSQL
pg_ctl -D "${PGDATA}" -o "-c listen_addresses='*'" -w start

echo "init-slave.sh completed"
