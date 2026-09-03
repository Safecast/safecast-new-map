#!/bin/bash
set -euo pipefail

# Nightly backup script for the production postgresql-5432 service.
# Deployed to /usr/local/bin/postgresql-5432-backup on the prod host,
# triggered by /etc/cron.d/postgresql-5432-backup (0 2 * * *).
#
# Previously this ran `docker exec postgresql-5432 pg_dumpall`, but Postgres
# was migrated off Docker to a native systemd service; docker exec then
# failed on every run, silently producing near-empty .sql.gz files for
# months because the failure happened mid-pipeline after the output file
# was already created by the shell redirection, and the script had no
# output-size check to catch it. The cron trigger itself was also later
# lost during the same migration, so even the empty-file writes stopped.

PGUSER="postgres"
SERVICE_NAME="postgresql-5432"

BACKUP_BASE="/backup"
HOSTNAME_STR="$(hostname -s 2>/dev/null || hostname)"
BACKUP_DIR="${BACKUP_BASE}/${HOSTNAME_STR}"

RETENTION_DAYS=3
MIN_SIZE_BYTES=1000000  # sanity floor; a real dump of this DB is many MB

DATE_STR="$(date +%F_%H-%M-%S)"
OUT_FILE="${BACKUP_DIR}/${SERVICE_NAME}_${DATE_STR}.sql.gz"
TMP_FILE="${OUT_FILE}.tmp"

mkdir -p "${BACKUP_DIR}"

sudo -u "${PGUSER}" pg_dumpall -U "${PGUSER}" | gzip -9 > "${TMP_FILE}"

SIZE=$(stat -c%s "${TMP_FILE}")
if [ "${SIZE}" -lt "${MIN_SIZE_BYTES}" ]; then
	echo "Backup FAILED: ${TMP_FILE} is only ${SIZE} bytes (expected >= ${MIN_SIZE_BYTES}); leaving it for inspection, not rotating it in." >&2
	mv "${TMP_FILE}" "${OUT_FILE}.suspect"
	exit 1
fi

mv "${TMP_FILE}" "${OUT_FILE}"

# Purge old backups for this service (only reached on success)
find "${BACKUP_DIR}" -type f -name "${SERVICE_NAME}_*.sql.gz" -mtime +${RETENTION_DAYS} -delete

echo "Backup OK: ${OUT_FILE} (${SIZE} bytes)"
