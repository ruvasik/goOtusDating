#!/bin/bash
set -e

echo "++++++++++++ Starting init-slave.sh"

# Wait for the PostgreSQL master to start
until pg_isready -h db-master -p 5432
do
  echo "++++++++++++ Waiting for master to be ready..."
  sleep 1
done

# Ensure the data directory is empty
echo "++++++++++++ Cleaning up existing data directory"
rm -rf /var/lib/postgresql/data/*

# Perform base backup
echo "++++++++++++ Performing base backup"
export PGPASSWORD=$DB_PASSWORD
pg_basebackup -h db-master -D /var/lib/postgresql/data -U $DB_USER -vP

# Create recovery.conf for replication
echo "++++++++++++ Creating recovery.conf for replication"
cat <<EOF > /var/lib/postgresql/data/recovery.conf
standby_mode = 'on'
primary_conninfo = 'host=db-master port=5432 user=$DB_USER password=$DB_PASSWORD'
trigger_file = '/tmp/postgresql.trigger.5432'
EOF

chown postgres:postgres /var/lib/postgresql/data/recovery.conf

echo "++++++++++++ init-slave.sh completed"

# Keep the container running
tail -f /dev/null
