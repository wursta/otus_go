FROM postgres

COPY migrations/0001_create_events_table.sql /docker-entrypoint-initdb.d/
COPY migrations/0002_alter_events_table_add_notified.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_USER calendar
ENV POSTGRES_PASSWORD calendar
ENV POSTGRES_DB calendar