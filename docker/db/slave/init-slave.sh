#!/bin/bash
set -e

echo "++++++ Starting init-slave.sh"

# Wait for the PostgreSQL master to start
until pg_isready -h db-master -p 5432
do
  echo "++++++ Waiting for master to be ready..."
  sleep 1
done

# Ensure the data directory is empty
echo "++++++ Cleaning up existing data directory"
rm -rf /var/lib/postgresql/data/*

# Perform base backup
echo "++++++ Performing base backup"
export PGPASSWORD=$DB_PASSWORD
pg_basebackup -h db-master -D /var/lib/postgresql/data -U $DB_USER -vP

# Create standby.signal for replication
echo "++++++ Creating standby.signal for replication"
touch /var/lib/postgresql/data/standby.signal

# Configure replication settings in postgresql.auto.conf
cat <<EOF >> /var/lib/postgresql/data/postgresql.auto.conf
primary_conninfo = 'host=db-master port=5432 user=$DB_USER password=$DB_PASSWORD'
restore_command = 'cp /var/lib/postgresql/data/archive/%f %p'
EOF

chown postgres:postgres /var/lib/postgresql/data/standby.signal
chown postgres:postgres /var/lib/postgresql/data/postgresql.auto.conf

echo "++++++ Cleaning up shared memory segments"
ipcs -m | grep postgres | awk '{print $2}' | xargs -I {} ipcrm -m {}

echo "++++++ init-slave.sh completed"

# Start the PostgreSQL server
echo "++++++ Starting PostgreSQL server"
pg_ctl -D /var/lib/postgresql/data -o "-c listen_addresses='*'" -w start

# Keep the container running
tail -f /dev/null
