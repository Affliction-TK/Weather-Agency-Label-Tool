# Deployment Guide

This guide covers deploying the Weather Agency Label Tool in production.

## Prerequisites

- Linux server (Ubuntu 20.04+ recommended)
- MySQL 8.0+
- Go 1.24+
- Node.js 20+
- Nginx (recommended for production)
- Domain name (optional)

## Deployment Steps

### 1. Server Setup

Update system packages:
```bash
sudo apt update && sudo apt upgrade -y
```

Install required packages:
```bash
sudo apt install -y git mysql-server nginx
```

### 2. Install Go

```bash
wget https://go.dev/dl/go1.24.10.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.10.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### 3. Install Node.js

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
```

### 4. Clone and Build Application

```bash
cd /opt
sudo git clone https://github.com/Affliction-TK/Weather-Agency-Label-Tool.git
cd Weather-Agency-Label-Tool

# Install dependencies and build
./setup.sh
```

### 5. Configure MySQL

```bash
sudo mysql_secure_installation

# Create database
sudo mysql -u root -p
```

In MySQL console:
```sql
CREATE DATABASE weather_label_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'weather_user'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT ALL PRIVILEGES ON weather_label_db.* TO 'weather_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

Import schema:
```bash
sudo mysql -u weather_user -p weather_label_db < schema.sql
```

### 6. Configure Environment

Create `.env` file:
```bash
sudo nano .env
```

Add:
```
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_NAME=weather_label_db
DB_USER=weather_user
DB_PASSWORD=your_secure_password
DB_PARAMS=parseTime=true&charset=utf8mb4&loc=Local
UPLOAD_DIR=/opt/Weather-Agency-Label-Tool/uploads
STATIC_DIR=/opt/Weather-Agency-Label-Tool/frontend/dist
```

### 7. Create Systemd Service

Create service file:
```bash
sudo nano /etc/systemd/system/weather-label.service
```

Add:
```ini
[Unit]
Description=Weather Agency Label Tool
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/Weather-Agency-Label-Tool
ExecStart=/opt/Weather-Agency-Label-Tool/server
Restart=on-failure
RestartSec=5s
Environment="PORT=8080"
Environment="DB_HOST=localhost"
Environment="DB_PORT=3306"
Environment="DB_NAME=weather_label_db"
Environment="DB_USER=weather_user"
Environment="DB_PASSWORD=your_secure_password"
Environment="DB_PARAMS=parseTime=true&charset=utf8mb4&loc=Local"
Environment="UPLOAD_DIR=/opt/Weather-Agency-Label-Tool/uploads"
Environment="STATIC_DIR=/opt/Weather-Agency-Label-Tool/frontend/dist"

[Install]
WantedBy=multi-user.target
```

Enable and start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable weather-label
sudo systemctl start weather-label
sudo systemctl status weather-label
```

### 8. Configure Nginx

Create Nginx config:
```bash
sudo nano /etc/nginx/sites-available/weather-label
```

Add:
```nginx
server {
    listen 80;
    server_name your-domain.com;  # Change to your domain or IP

    client_max_body_size 50M;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable site:
```bash
sudo ln -s /etc/nginx/sites-available/weather-label /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 9. SSL/HTTPS (Optional but Recommended)

Install Certbot:
```bash
sudo apt install -y certbot python3-certbot-nginx
```

Get SSL certificate:
```bash
sudo certbot --nginx -d your-domain.com
```

### 10. Firewall Configuration

```bash
sudo ufw allow 'Nginx Full'
sudo ufw allow OpenSSH
sudo ufw enable
```

## Maintenance

### View Logs
```bash
# Application logs
sudo journalctl -u weather-label -f

# Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Restart Service
```bash
sudo systemctl restart weather-label
```

### Update Application
```bash
cd /opt/Weather-Agency-Label-Tool
sudo git pull
cd frontend && sudo npm install && sudo npm run build
cd .. && sudo go build -o server main.go
sudo systemctl restart weather-label
```

### Backup Database
```bash
mysqldump -u weather_user -p weather_label_db > backup_$(date +%Y%m%d).sql
```

### Restore Database
```bash
mysql -u weather_user -p weather_label_db < backup_YYYYMMDD.sql
```

## Monitoring

### Check Service Status
```bash
sudo systemctl status weather-label
```

### Check Disk Space for Uploads
```bash
du -sh /opt/Weather-Agency-Label-Tool/uploads
```

### Monitor Database Size
```sql
SELECT 
    table_schema AS 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.tables
WHERE table_schema = 'weather_label_db'
GROUP BY table_schema;
```

## Troubleshooting

### Service won't start
```bash
sudo journalctl -u weather-label -n 50 --no-pager
```

### Database connection issues
- Check MySQL is running: `sudo systemctl status mysql`
- Verify credentials in `.env`
- Check MySQL user permissions

### File upload issues
- Check uploads directory permissions: `sudo chown -R www-data:www-data /opt/Weather-Agency-Label-Tool/uploads`
- Verify Nginx client_max_body_size is sufficient

### High memory usage
- Monitor with: `htop` or `free -m`
- Adjust MySQL buffer pool size if needed
- Consider using a CDN for image serving

## Security Recommendations

1. **Change default passwords** in MySQL and `.env`
2. **Enable firewall** (UFW or iptables)
3. **Use HTTPS** with Let's Encrypt
4. **Regular updates**: Keep OS, Go, Node.js, and dependencies updated
5. **Backup regularly**: Database and uploads folder
6. **Monitor logs** for suspicious activity
7. **Restrict database access** to localhost only
8. **Use strong passwords** for all accounts
9. **Consider rate limiting** in Nginx to prevent abuse
10. **Keep uploads directory** outside web root or restrict access
