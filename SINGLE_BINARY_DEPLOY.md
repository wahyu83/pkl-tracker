# Deploy PKL Tracker — Single Binary (no Nginx, no Node.js on VPS)

## Perbandingan dengan Deploy Standar

| Aspek | Deploy Standar (DEPLOY.md) | Single Binary (ini) |
|---|---|---|
| Install di VPS | Go SDK + Node.js + Nginx + PostgreSQL | PostgreSQL saja |
| Frontend | Dilayani Nginx dari `dist/` | Ditanam ke binary Go via `embed` |
| HTTPS | Nginx + Certbot | Caddy (reverse proxy ringan) |
| Proses | 2 proses (Go binary + Nginx) | 1 proses (binary + Caddy) |
| VPS RAM | Minimal 1 GB | Minimal 512 MB |
| Update | git pull + go build + npm build | `make build` di laptop → scp binary |

## Prasyarat

- VPS Ubuntu 22.04 atau 24.04 (minimal 512 MB RAM, 10 GB disk)
- Domain yang sudah diarahkan ke IP VPS
- **Go SDK dan Node.js terinstall di laptop/CI** (bukan di VPS)
- Akses SSH ke VPS

---

## Langkah 1: Build Binary di Laptop

```bash
cd /path/ke/project/pkl-tracker
make build
```

Hasil: `backend/pkl-server` (~32 MB, static binary dengan frontend di dalamnya).

> `make build` menjalankan: `npm run build` → copy `dist/` ke `backend/public/` → `go build` dengan embed.

---

## Langkah 2: Install PostgreSQL di VPS

SSH ke VPS:

```bash
ssh user@ip-vps-anda

sudo apt update
sudo apt install -y postgresql
sudo systemctl enable --now postgresql

sudo -u postgres psql <<EOF
CREATE USER pkl_user WITH PASSWORD 'password_anda';
CREATE DATABASE pkl_db OWNER pkl_user;
GRANT ALL PRIVILEGES ON DATABASE pkl_db TO pkl_user;
EOF
```

---

## Langkah 3: Install Caddy (Reverse Proxy + HTTPS Otomatis)

```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install -y caddy
```

---

## Langkah 4: Upload Binary ke VPS

Dari laptop:

```bash
scp backend/pkl-server user@vps:/opt/pkl-tracker/
ssh user@vps "chmod +x /opt/pkl-tracker/pkl-server"
```

---

## Langkah 5: Setup systemd Service

SSH ke VPS, buat file service:

```bash
sudo nano /etc/systemd/system/pkl-tracker.service
```

Isi (ganti `password_anda` dan `JWT_SECRET`):

```ini
[Unit]
Description=PKL Tracker
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/pkl-tracker
Environment="DB_HOST=127.0.0.1"
Environment="DB_PORT=5432"
Environment="DB_USER=pkl_user"
Environment="DB_PASS=password_anda"
Environment="DB_NAME=pkl_db"
Environment="JWT_SECRET=$(openssl rand -hex 32)"
Environment="SERVER_PORT=8082"
ExecStart=/opt/pkl-tracker/pkl-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Jalankan:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now pkl-tracker
sudo systemctl status pkl-tracker
```

Verifikasi:

```bash
curl http://localhost:8082/api/login -X POST \
  -H 'Content-Type: application/json' \
  -d '{"nis_nip_nik":"ADM-001","password":"admin123"}'
```

---

## Langkah 6: Setup Caddy (HTTPS + Reverse Proxy)

```bash
sudo nano /etc/caddy/Caddyfile
```

Isi (ganti domain):

```
pkl.sekolah-anda.sch.id {
    reverse_proxy localhost:8082
}
```

> Caddy otomatis mengurus sertifikat SSL Let's Encrypt. Tidak perlu `certbot --nginx` lagi.

Reload:

```bash
sudo systemctl reload caddy
```

---

## Langkah 7: Buka Firewall

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

---

## Verifikasi

Buka `https://pkl.sekolah-anda.sch.id/` — halaman login muncul. Login dengan akun test:

| Role  | NIS/NIP/NIK | Password   |
|-------|-------------|------------|
| Admin | `ADM-001`   | `admin123` |
| Guru  | `19850101`  | `guru123`  |
| Siswa | `20230001`  | `siswa123` |
| DUDI  | `D-001`     | `dudi123`  |

---

## Update Aplikasi

Di laptop:

```bash
git pull
make build
scp backend/pkl-server user@vps:/opt/pkl-tracker/
ssh user@vps "sudo systemctl restart pkl-tracker"
```

---

## Struktur di VPS

```
/opt/pkl-tracker/
└── pkl-server          ← satu-satunya file (binary ~32 MB)

/etc/systemd/system/pkl-tracker.service
/etc/caddy/Caddyfile
```

---

## Troubleshooting

### Backend tidak berjalan

```bash
sudo systemctl status pkl-tracker
sudo journalctl -u pkl-tracker -f
```

### Database tidak terkoneksi

```bash
sudo systemctl status postgresql
PGPASSWORD='password_anda' psql -U pkl_user -h 127.0.0.1 -d pkl_db -c "SELECT 1;"
```

### Caddy gagal

```bash
sudo systemctl status caddy
sudo journalctl -u caddy -f
```

### Integrasi Google Drive

Tambahkan environment variable di `/etc/systemd/system/pkl-tracker.service`:

```ini
Environment="GDRIVE_CREDENTIALS=/opt/pkl-tracker/service-account.json"
Environment="GDRIVE_FOLDER_ID=1ABC123..."
```

Upload file credential:

```bash
scp service-account.json user@vps:/opt/pkl-tracker/
ssh user@vps "sudo systemctl daemon-reload && sudo systemctl restart pkl-tracker"
```

---

## Kenapa Single Binary?

1. **VPS lebih kecil** — tidak perlu install Go SDK (200+ MB) atau Node.js (100+ MB)
2. **Deploy lebih cepat** — 1 file SCP, bukan 3 langkah build di VPS
3. **Lebih aman** — tidak ada Nginx config yang bisa salah konfigurasi, binary self-contained
4. **Caddy > Certbot** — HTTPS otomatis, konfigurasi 2 baris
5. **Atomic deploy** — binary baru langsung aktif setelah `systemctl restart`, tidak ada momen frontend/backend tidak sinkron
