# Deploy PKL Tracker — Single Binary (no Node.js on VPS)

## Perbandingan dengan Deploy Standar

| Aspek | Deploy Standar (DEPLOY.md) | Single Binary (ini) |
|---|---|---|
| Install di VPS | Go SDK + Node.js + Nginx + PostgreSQL | PostgreSQL + Nginx |
| Frontend | Dilayani Nginx dari `dist/` | Ditanam ke binary Go via `embed` |
| HTTPS | Nginx + Certbot | Nginx + Certbot (Let's Encrypt) |
| Proses | 2 proses (Go binary + Nginx) | 2 proses (binary + Nginx) |
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

## Langkah 3: Install Nginx + Certbot (Reverse Proxy + HTTPS)

```bash
ssh user@ip-vps-anda

sudo apt install -y nginx certbot python3-certbot-nginx
sudo systemctl enable --now nginx
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

> **PENTING — JWT_SECRET:** jangan tulis `$(openssl rand -hex 32)` langsung di file
> service. systemd **tidak** menjalankan perintah shell di dalam `Environment=`,
> sehingga secret-nya menjadi string literal dan bisa ditebak siapa pun yang melihat
> file/repo ini (token admin bisa dipalsukan). Generate dulu di terminal, lalu tempel
> **hasilnya** sebagai nilai tetap:

```bash
openssl rand -hex 32
# contoh output: 8f3c...a91  <- salin nilai ini ke Environment="JWT_SECRET=..."
```

Isi (ganti `password_anda` dan `<hasil-openssl-rand-hex-32>`):

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
Environment="JWT_SECRET=<hasil-openssl-rand-hex-32>"
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

## Langkah 6: Setup Nginx Reverse Proxy

Buat konfigurasi Nginx:

```bash
sudo nano /etc/nginx/sites-available/pkl-tracker
```

Isi (ganti `pkl.sekolah-anda.sch.id` dengan domain anda):

```nginx
server {
    listen 80;
    server_name pkl.sekolah-anda.sch.id;

    client_max_body_size 50M;

    location / {
        proxy_pass http://127.0.0.1:8082;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /uploads/ {
        proxy_pass http://127.0.0.1:8082;
        proxy_set_header Host $host;
    }
}
```

Aktifkan site:

```bash
sudo ln -sf /etc/nginx/sites-available/pkl-tracker /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

## Langkah 7: Aktifkan HTTPS dengan Certbot

```bash
sudo certbot --nginx -d pkl.sekolah-anda.sch.id
```

> Certbot otomatis mengambil sertifikat Let's Encrypt dan memperbarui konfigurasi Nginx. Sertifikat diperpanjang otomatis via cron/systemd timer.

---

## Langkah 8: Buka Firewall

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
/etc/nginx/sites-available/pkl-tracker
/etc/nginx/sites-enabled/pkl-tracker -> ../sites-available/pkl-tracker
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

### Nginx gagal

```bash
sudo systemctl status nginx
sudo journalctl -u nginx -f
sudo nginx -t
```

### SSL / Certbot gagal

```bash
sudo certbot --nginx -d pkl.sekolah-anda.sch.id --dry-run
sudo certbot renew --dry-run
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
3. **Lebih aman** — binary self-contained, tidak perlu expose port aplikasi langsung ke internet
4. **Atomic deploy** — binary baru langsung aktif setelah `systemctl restart`, tidak ada momen frontend/backend tidak sinkron
