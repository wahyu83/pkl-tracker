#!/bin/bash
set -e

# ============================================================
# PKL Tracker - Update Script (Single Binary Method)
#
# Jalankan dari LAPTOP (bukan di VPS):
#   ./update.sh [user@vps] [path_binary_di_vps] [ssh_port]
#
# Contoh:
#   ./update.sh root@203.0.113.10
#   ./update.sh deploy@vps.example.com /opt/pkl-tracker 30333
#
# Alur:
#   1. git pull (ambil update terbaru dari GitHub)
#   2. make build (npm build -> backend/public -> go build)
#   3. scp backend/pkl-server ke VPS
#   4. restart service pkl-tracker di VPS
# ============================================================

# --- Config ---
VPS="${1:?Usage: ./update.sh user@vps [remote_path] [ssh_port]}"
REMOTE_PATH="${2:-/opt/pkl-tracker}"
SSH_PORT="${3:-22}"
SERVICE_NAME="pkl-tracker"

echo "=== PKL Tracker Update ==="
echo "VPS target : $VPS (port $SSH_PORT)"
echo "Remote path: $REMOTE_PATH"
echo ""

# --- 1. Tarik perubahan terbaru ---
if git rev-parse --git-dir >/dev/null 2>&1; then
  echo "[1/4] git pull..."
  git pull --ff-only
else
  echo "[1/4] Bukan git repo, lewati 'git pull'."
  echo "       (Pastikan folder ini sudah versi terbaru!)"
fi
echo ""

# --- 2. Build ulang binary ---
echo "[2/4] make build..."
make build
echo ""

# --- 3. Upload binary ke VPS ---
echo "[3/4] Upload binary ke $VPS:$REMOTE_PATH/ (port $SSH_PORT)..."
# Upload ke file temp dulu, lalu mv (atomic replace) supaya tidak error
# saat service lama masih berjalan memegang file pkl-server.
scp -P "$SSH_PORT" backend/pkl-server "$VPS:$REMOTE_PATH/pkl-server.new"
ssh -p "$SSH_PORT" "$VPS" "mv -f $REMOTE_PATH/pkl-server.new $REMOTE_PATH/pkl-server && chmod +x $REMOTE_PATH/pkl-server"
echo ""

# --- 4. Restart service di VPS ---
echo "[4/4] Restart service $SERVICE_NAME di VPS..."
ssh -p "$SSH_PORT" "$VPS" "sudo systemctl restart $SERVICE_NAME && sudo systemctl status $SERVICE_NAME --no-pager | head -5"
echo ""

echo "=== Update selesai! ==="
echo "Cek aplikasi: curl https://domain-anda/api/login (atau buka di browser)"
echo ""
echo "Catatan:"
echo "  - Database tidak perlu diubah; AutoMigrate otomatis saat binary start."
echo "  - Port SSH default 22; untuk port lain: ./update.sh user@vps /opt/pkl-tracker <port>"
