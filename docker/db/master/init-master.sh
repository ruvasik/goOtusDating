#!/bin/bash
set -e

echo "//////////// Starting init-master.sh with user $DB_USER and database $DB_DB"

export PGPASSWORD=$DB_PASSWORD

# Create archive directory
mkdir -p /var/lib/postgresql/data/archive
chown postgres:postgres /var/lib/postgresql/data/archive

# Add replication configuration to postgresql.conf
echo "//////////// Appending replication configuration to postgresql.conf"
echo "listen_addresses = '*'" >> /var/lib/postgresql/data/postgresql.conf
echo "wal_level = replica" >> /var/lib/postgresql/data/postgresql.conf
echo "max_wal_senders = 10" >> /var/lib/postgresql/data/postgresql.conf
echo "max_replication_slots = 2" >> /var/lib/postgresql/data/postgresql.conf
echo "hot_standby = on" >> /var/lib/postgresql/data/postgresql.conf
echo "hot_standby_feedback = on" >> /var/lib/postgresql/data/postgresql.conf
echo "wal_keep_size = 512MB" >> /var/lib/postgresql/data/postgresql.conf
echo "archive_mode = on" >> /var/lib/postgresql/data/postgresql.conf
echo "archive_command = 'cp %p /var/lib/postgresql/data/archive/%f'" >> /var/lib/postgresql/data/postgresql.conf

# Add replication configuration to pg_hba.conf
echo "host replication $DB_USER 127.0.0.1/32 trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host replication $DB_USER db-master trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host replication $DB_USER db-slave trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host replication $DB_USER 0.0.0.0/0 trust" >> /var/lib/postgresql/data/pg_hba.conf
echo "host all all 0.0.0.0/0 md5" >> /var/lib/postgresql/data/pg_hba.conf

# Restart PostgreSQL to apply changes
pg_ctl -D "$PGDATA" -m fast -w restart

echo "//////////// Configured"

# Create user and database
psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" <<-EOSQL
    CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
    ALTER USER $DB_USER WITH SUPERUSER;
    CREATE DATABASE $DB_NAME OWNER $DB_USER;
EOSQL

echo "//////////// Created user & db"

# Create table in the database
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

echo "//////////// init-master.sh completed"
