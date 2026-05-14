# Deploy PKL Tracker di Ubuntu 22.04

## Prasyarat

- VPS Ubuntu 22.04 (minimal 1 GB RAM, 10 GB disk)
- Domain yang sudah diarahkan ke IP VPS (contoh: `pkl.sekolah-anda.sch.id`)
- Akses SSH ke VPS

---

## Langkah 1: Masuk ke VPS

```bash
ssh user@ip-vps-anda
```

---

## Langkah 2: Install Go

```bash
sudo snap install go --classic
go version
# Harus muncul: go version go1.22.x ...
```

---

## Langkah 3: Install Node.js & npm

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
node --version
# Harus muncul: v20.x.x
```

---

## Langkah 4: Install PostgreSQL

```bash
sudo apt update
sudo apt install -y postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

---

## Langkah 5: Install Nginx

```bash
sudo apt install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

---

## Langkah 6: Setup Database

Ganti `password_anda` dengan password PostgreSQL yang diinginkan:

```bash
sudo -u postgres psql <<EOF
CREATE USER pkl_user WITH PASSWORD 'password_anda';
CREATE DATABASE pkl_db OWNER pkl_user;
GRANT ALL PRIVILEGES ON DATABASE pkl_db TO pkl_user;
\q
EOF
```

Verifikasi:

```bash
PGPASSWORD='password_anda' psql -U pkl_user -h 127.0.0.1 -d pkl_db -c "SELECT 1;"
# Harus muncul: 1
```

---

## Langkah 7: Clone Repository

```bash
cd /opt
sudo mkdir -p pkl-tracker
sudo chown $USER:$USER pkl-tracker
git clone https://github.com/Tarilusiana/pkl-tracker.git pkl-tracker
cd pkl-tracker
```

---

## Langkah 8: Build Backend (Go)

```bash
cd /opt/pkl-tracker/backend
go build -o pkl-server .
ls -la pkl-server
# Harus ada file: pkl-server
```

---

## Langkah 9: Setup Systemd Service (Backend)

Buat file service:

```bash
sudo nano /etc/systemd/system/pkl-tracker.service
```

Isi dengan (ganti `password_anda`):

```ini
[Unit]
Description=PKL Tracker Backend
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/pkl-tracker/backend
Environment="DB_HOST=127.0.0.1"
Environment="DB_PORT=5432"
Environment="DB_USER=pkl_user"
Environment="DB_PASS=password_anda"
Environment="DB_NAME=pkl_db"
Environment="JWT_SECRET=pkl-tracker-secret-key-2026"
Environment="SERVER_PORT=8082"
ExecStart=/opt/pkl-tracker/backend/pkl-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Simpan (Ctrl+O, Enter, Ctrl+X), lalu jalankan:

```bash
sudo systemctl daemon-reload
sudo systemctl enable pkl-tracker
sudo systemctl start pkl-tracker
```

Verifikasi backend:

```bash
sudo systemctl status pkl-tracker
# Harus: active (running)

curl http://localhost:8082/api/login -X POST -H 'Content-Type: application/json' -d '{"nis_nip_nik":"ADM-001","password":"admin123"}'
# Harus muncul response JSON
```

---

## Langkah 10: Build Frontend (Vue.js)

```bash
cd /opt/pkl-tracker/frontend
npm install
npm run build
ls dist/
# Harus ada: index.html, assets/, sw.js, ...
```

---

## Langkah 11: Setup Nginx

```bash
sudo nano /etc/nginx/sites-available/pkl-tracker
```

Isi (ganti `pkl.sekolah-anda.sch.id` dengan domain Anda):

```nginx
server {
    listen 80;
    server_name pkl.sekolah-anda.sch.id;

    root /opt/pkl-tracker/frontend/dist;
    index index.html;

    # API reverse proxy ke backend Go
    location /api/ {
        proxy_pass http://127.0.0.1:8082;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        client_max_body_size 50m;
    }

    # Uploaded files (foto absensi, dokumentasi jurnal)
    location /uploads/ {
        proxy_pass http://127.0.0.1:8082;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
    }

    # SPA fallback (Vue Router)
    location / {
        try_files $uri $uri/ /index.html;
    }

    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml text/javascript;
}
```

Aktifkan:

```bash
sudo ln -sf /etc/nginx/sites-available/pkl-tracker /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
# Harus: syntax is ok, test is successful
sudo systemctl reload nginx
```

---

## Langkah 12: Buka Firewall

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

---

## Langkah 13: Verifikasi Akses

Buka browser, akses:

```
http://pkl.sekolah-anda.sch.id/
```

Login dengan akun test:

| Role  | NIS/NIP/NIK | Password    |
|-------|-------------|-------------|
| Admin | `ADM-001`   | `admin123`  |
| Guru  | `19850101`  | `guru123`   |
| Siswa | `20230001`  | `siswa123`  |
| DUDI  | `D-001`     | `dudi123`   |

---

## Langkah 14: HTTPS (SSL Gratis)

```bash
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot --nginx -d pkl.sekolah-anda.sch.id
# Ikuti petunjuk, pilih redirect HTTP ke HTTPS
```

Setelah selesai, akses via HTTPS:

```
https://pkl.sekolah-anda.sch.id/
```

---

## Troubleshooting

### Backend tidak berjalan

```bash
sudo systemctl status pkl-tracker
sudo journalctl -u pkl-tracker -f
```

### Nginx error

```bash
sudo nginx -t
sudo tail -f /var/log/nginx/error.log
```

### Database tidak terkoneksi

```bash
sudo systemctl status postgresql
PGPASSWORD='password_anda' psql -U pkl_user -h 127.0.0.1 -d pkl_db -c "SELECT count(*) FROM users;"
```

### Update aplikasi

```bash
cd /opt/pkl-tracker
git pull

# Rebuild backend
cd backend && go build -o pkl-server . && cd ..
sudo systemctl restart pkl-tracker

# Rebuild frontend
cd frontend && npm install && npm run build && cd ..
```

---

## Struktur Setelah Deploy

```
/opt/pkl-tracker/
├── backend/
│   ├── pkl-server          ← binary Go
│   ├── main.go
│   └── ...
├── frontend/
│   ├── dist/               ← hasil build Vue.js
│   │   ├── index.html
│   │   ├── assets/
│   │   └── sw.js
│   └── ...
└── prd.md

/etc/nginx/sites-available/pkl-tracker   ← config Nginx
/etc/systemd/system/pkl-tracker.service   ← service backend
```

---

## Integrasi Google Drive (Foto Absensi & Dokumentasi)

### 1. Buat Service Account

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Buat project baru atau pilih existing project
3. **APIs & Services → Library** → cari **Google Drive API** → Enable
4. **APIs & Services → Credentials** → Create Credentials → **Service Account**
5. Isi nama (contoh: `pkl-tracker-uploader`), klik Done
6. Klik service account yang baru dibuat → tab **Keys** → Add Key → Create New Key → JSON → Download
7. Upload file JSON ke VPS: `scp service-account.json user@vps:/opt/pkl-tracker/backend/`

### 2. Buat Folder Google Drive

1. Buka [Google Drive](https://drive.google.com/)
2. Buat folder baru (contoh: `PKL_Photos`)
3. Buka folder, lihat URL: `https://drive.google.com/drive/folders/1ABC123...` → copy ID setelah `folders/`
4. Klik kanan folder → Share → tambahkan email service account (ada di file JSON, field `client_email`) sebagai Editor

### 3. Konfigurasi Backend

Edit service file:

```bash
sudo nano /etc/systemd/system/pkl-tracker.service
```

Tambahkan dua environment variable di bagian `[Service]`:

```ini
Environment="GDRIVE_CREDENTIALS=/opt/pkl-tracker/backend/service-account.json"
Environment="GDRIVE_FOLDER_ID=1ABC123..."    ← ganti dengan ID folder
```

Restart:

```bash
sudo systemctl daemon-reload
sudo systemctl restart pkl-tracker
```

Cek log:

```bash
sudo journalctl -u pkl-tracker -f | grep Storage
# Harus muncul: [Storage] Google Drive connected
```

### 4. Tanpa Google Drive (Fallback)

Jika environment `GDRIVE_CREDENTIALS` tidak diset, foto otomatis disimpan ke folder lokal `backend/uploads/`. Nginx sudah dikonfigurasi untuk menyajikan file dari path `/uploads/`.

---

## Catatan

- JWT Secret di `JWT_SECRET` sebaiknya diganti dengan string random panjang (`openssl rand -hex 32`).
- Backup database PostgreSQL secara berkala: `pg_dump -U pkl_user pkl_db > backup.sql`
