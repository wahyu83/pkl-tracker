#!/bin/bash
set -e

echo "=== PKL Tracker Deployment ==="

# --- Config ---
APP_DIR="/opt/pkl-tracker"
DOMAIN="${1:-pkl.example.com}"
DB_PASS="${2:-100%Bisa}"
JWT_SECRET="${3:-$(openssl rand -hex 32)}"

echo "Domain: $DOMAIN"
echo "App directory: $APP_DIR"

# --- Install dependencies ---
echo "Installing dependencies..."
sudo apt update -y
sudo apt install -y nginx postgresql git curl wget

# --- Setup PostgreSQL ---
echo "Setting up PostgreSQL..."
sudo -u postgres psql -c "CREATE USER pkl_user WITH PASSWORD '$DB_PASS';" 2>/dev/null || true
sudo -u postgres psql -c "CREATE DATABASE pkl_db OWNER pkl_user;" 2>/dev/null || true
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pkl_db TO pkl_user;" 2>/dev/null || true

# --- Create app directory ---
sudo mkdir -p $APP_DIR
sudo chown $USER:$USER $APP_DIR

# --- Build & deploy backend ---
echo "Building backend..."
cd $APP_DIR
cp -r /home/$USER/pkl-tracker/backend/* $APP_DIR/backend/ 2>/dev/null || git clone https://github.com/Tarilusiana/pkl-tracker.git .
cd $APP_DIR/backend
go build -o pkl-server .

# --- Backend systemd service ---
echo "Creating systemd service..."
sudo tee /etc/systemd/system/pkl-tracker.service > /dev/null << SERVICE
[Unit]
Description=PKL Tracker Backend
After=network.target postgresql.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$APP_DIR/backend
Environment="DB_HOST=127.0.0.1"
Environment="DB_PORT=5432"
Environment="DB_USER=pkl_user"
Environment="DB_PASS=$DB_PASS"
Environment="DB_NAME=pkl_db"
Environment="JWT_SECRET=$JWT_SECRET"
Environment="SERVER_PORT=8082"
ExecStart=$APP_DIR/backend/pkl-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
SERVICE

sudo systemctl daemon-reload
sudo systemctl enable pkl-tracker
sudo systemctl start pkl-tracker

# --- Build frontend ---
echo "Building frontend..."
cd $APP_DIR/frontend
npm install
npm run build

# --- Nginx config ---
echo "Setting up Nginx..."
sudo tee /etc/nginx/sites-available/pkl-tracker > /dev/null << NGINX
server {
    listen 80;
    server_name $DOMAIN;

    root $APP_DIR/frontend/dist;
    index index.html;

    # API reverse proxy
    location /api/ {
        proxy_pass http://127.0.0.1:8082;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;

        # Increase upload size for photos/CSV
        client_max_body_size 50m;
    }

    # SPA fallback
    location / {
        try_files \$uri \$uri/ /index.html;
    }

    # Gzip
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript image/svg+xml;
}
NGINX

sudo ln -sf /etc/nginx/sites-available/pkl-tracker /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t && sudo systemctl reload nginx

# --- Firewall ---
sudo ufw allow 80/tcp 2>/dev/null || true
sudo ufw allow 443/tcp 2>/dev/null || true
sudo ufw allow 22/tcp 2>/dev/null || true

echo ""
echo "=== Deployment complete! ==="
echo "Frontend:  http://$DOMAIN"
echo "Backend:   http://$DOMAIN/api"
echo ""
echo "Test accounts:"
echo "  Admin: NIS=ADM-001  Pass=admin123"
echo "  Guru:  NIP=19850101 Pass=guru123"
echo "  Siswa: NIS=20230001 Pass=siswa123"
echo "  DUDI:  NIK=D-001    Pass=dudi123"
echo ""
echo "Next: Run 'sudo certbot --nginx -d $DOMAIN' for HTTPS"
