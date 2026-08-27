#!/bin/bash
set -e

# ============================================================
# PKL Tracker - Backup Database & Uploads
#
# Jalankan di VPS:
#   ./backup.sh                  # backup manual
#   ./backup.sh --install-cron   # jadwalkan backup harian jam 02:00
#
# Restore database:
#   gunzip -c pkl_db_*.sql.gz | \
#     PGPASSWORD='password_anda' psql -h 127.0.0.1 -U pkl_user -d pkl_db
#
# Konfigurasi bisa di-override via env:
#   DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME,
#   BACKUP_DIR, RETENTION_DAYS, INCLUDE_UPLOADS, UPLOADS_DIR
# ============================================================

# --- Config (default sama dengan deploy) ---
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-pkl_user}"
DB_PASS="${DB_PASS:-100%Bisa}"
DB_NAME="${DB_NAME:-pkl_db}"
BACKUP_DIR="${BACKUP_DIR:-/opt/pkl-tracker/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-14}"
INCLUDE_UPLOADS="${INCLUDE_UPLOADS:-true}"
UPLOADS_DIR="${UPLOADS_DIR:-/opt/pkl-tracker/uploads}"

# --- Install cron backup harian (jam 02:00) ---
if [ "$1" = "--install-cron" ]; then
  SCRIPT_PATH="$(cd "$(dirname "$0")" && pwd)/backup.sh"
  CRON_LINE="0 2 * * * $SCRIPT_PATH >> /var/log/pkl-backup.log 2>&1"
  ( crontab -l 2>/dev/null | grep -v 'pkl-backup' ; echo "$CRON_LINE" ) | crontab -
  echo "Cron terpasang:"
  echo "  $CRON_LINE"
  echo "Cek: crontab -l"
  exit 0
fi

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DB_FILE="$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.sql.gz"
UPLOADS_FILE="$BACKUP_DIR/uploads_${TIMESTAMP}.tar.gz"

mkdir -p "$BACKUP_DIR"
chmod 700 "$BACKUP_DIR"

echo "=== PKL Tracker Backup ==="
echo "Database : $DB_NAME @ $DB_HOST:$DB_PORT (user $DB_USER)"
echo "Target   : $BACKUP_DIR"
echo ""

# --- 1. Dump database (gzip) ---
echo "[1/3] pg_dump -> $DB_FILE"
PGPASSWORD="$DB_PASS" pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
  --no-owner --no-privileges | gzip > "$DB_FILE"
chmod 600 "$DB_FILE"
echo ""

# --- 2. Backup folder uploads (foto lokal, opsional) ---
if [ "$INCLUDE_UPLOADS" = "true" ] && [ -d "$UPLOADS_DIR" ]; then
  echo "[2/3] tar uploads -> $UPLOADS_FILE"
  tar -czf "$UPLOADS_FILE" -C "$(dirname "$UPLOADS_DIR")" "$(basename "$UPLOADS_DIR")"
  chmod 600 "$UPLOADS_FILE"
  echo ""
else
  echo "[2/3] Lewati backup uploads (folder tidak ada / INCLUDE_UPLOADS=false)"
  echo ""
fi

# --- 3. Hapus backup lama (retensi) ---
echo "[3/3] Hapus backup lebih dari $RETENTION_DAYS hari..."
find "$BACKUP_DIR" -type f \( -name "*.sql.gz" -o -name "*.tar.gz" \) -mtime +"$RETENTION_DAYS" -delete

echo ""
echo "=== Backup selesai! ==="
ls -lh "$BACKUP_DIR" | tail -n +2
echo ""
echo "Restore DB : gunzip -c <file>.sql.gz | PGPASSWORD='...' psql -h $DB_HOST -U $DB_USER -d $DB_NAME"
