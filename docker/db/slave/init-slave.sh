#!/bin/bash
set -e

echo "Starting init-slave.sh"

#cp /docker-entrypoint-initdb.d/pg_hba.conf "$PGDATA/pg_hba.conf"

until psql -h "db-master" -U "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c '\q'; do
  >&2 echo "Master is unavailable - sleeping"
  sleep 5
done

PGPASSWORD=password pg_basebackup -h db-master -D ${PGDATA} -U replicator -P -R

#echo "primary_conninfo = 'host=db-master port=5432 user=${DB_USER}_replicator password=${DB_PASSWORD}'" >> ${PGDATA}/postgresql.conf

cp /docker-entrypoint-initdb.d/postgresql.conf ${PGDATA}/postgresql.conf
#cp /docker-entrypoint-initdb.d/pg_hba.conf ${PGDATA}/pg_hba.conf

pg_ctl -D "${PGDATA}" -o "-c listen_addresses='*'" -w start

