FROM postgres:16.3

# Copy initialization scripts
COPY init-master.sh /docker-entrypoint-initdb.d/

# Ensure scripts are executable
RUN chmod +x /docker-entrypoint-initdb.d/init-master.sh
