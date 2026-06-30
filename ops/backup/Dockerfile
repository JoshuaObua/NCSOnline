FROM alpine:3.21
RUN apk add --no-cache postgresql16-client rclone tzdata
COPY scripts/backup-postgres.sh /usr/local/bin/backup-postgres
RUN chmod 0555 /usr/local/bin/backup-postgres
ENTRYPOINT ["/bin/sh", "-c", "while true; do /usr/local/bin/backup-postgres; sleep ${BACKUP_INTERVAL_SECONDS:-86400}; done"]
