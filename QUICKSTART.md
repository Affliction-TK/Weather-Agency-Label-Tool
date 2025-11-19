# Quick Start Guide

Get the Weather Agency Label Tool up and running in 5 minutes!

## Option 1: Quick Start with Docker (Recommended)

This is the fastest way to get started:

```bash
# 1. Clone the repository
git clone https://github.com/Affliction-TK/Weather-Agency-Label-Tool.git
cd Weather-Agency-Label-Tool

# 2. Run the automated setup
./setup.sh

# 3. Start MySQL with Docker
docker-compose up -d

# 4. Wait a few seconds for MySQL to initialize (first time only)
sleep 10

# 5. Copy and edit environment file
cp .env.example .env
# Edit .env if needed (defaults should work with Docker)

# 6. Build and run
make build
./server
```

That's it! Open http://localhost:8080 in your browser.

## Option 2: Manual Setup

If you prefer manual setup or don't have Docker:

### 1. Install Prerequisites

- Go 1.24+
- Node.js 20+
- MySQL 8.0+

### 2. Clone and Setup

```bash
git clone https://github.com/Affliction-TK/Weather-Agency-Label-Tool.git
cd Weather-Agency-Label-Tool
```

### 3. Setup Database

```bash
# Login to MySQL
mysql -u root -p

# Create database and user
CREATE DATABASE weather_label_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'weather_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON weather_label_db.* TO 'weather_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;

# Import schema
mysql -u weather_user -p weather_label_db < schema.sql
```

### 4. Configure Application

复制示例配置：

```bash
cp .env.example .env
```

编辑 `.env`，按需更新 `DB_HOST/PORT/NAME/USER/PASSWORD` 等字段（也可通过 `DB_DSN` 自定义完整连接串）。

### 5. Install Dependencies

```bash
# Backend
go mod download

# Frontend
cd frontend
npm install
cd ..
```

### 6. Build and Run

```bash
# Build frontend
cd frontend
npm run build
cd ..

# Build backend
go build -o server main.go

# Run
./server
```

Open http://localhost:8080 in your browser.

## Using the Application

### 1. Upload Images

- Click the "上传图片" (Upload Images) tab
- Drag and drop images or click to select files
- Click "上传" (Upload) button

### 2. Annotate Images

- Images appear in the left sidebar
- Unannotated images are marked with orange badge
- Click an image to annotate it
- Fill in all required fields:
  - **天气类型** (Weather Type): 大雾, 结冰, or 积劳
  - **严重等级** (Severity): 无, 轻度, 中度, or 重度
  - **观测时间** (Observation Time): Date and time
  - **地点** (Location): Location name
  - **经度** (Longitude): Longitude coordinate
  - **维度** (Latitude): Latitude coordinate
  - **监测点** (Station): Auto-selected based on coordinates
- Click "保存标注" (Save Annotation)

### 3. Navigate Images

- Click any image in the left sidebar to view/edit
- System automatically moves to next unannotated image after saving
- Annotated images show green badge
- When all images are annotated, you'll see a completion message

## Keyboard Shortcuts

- Tab: Navigate between form fields
- Enter: Submit form when focused on a button
- Esc: Close dropdowns/dialogs (if implemented)

## Troubleshooting

### "Database connection failed"
- Check MySQL is running: `systemctl status mysql` (Linux) or check MySQL in services (Windows)
- Verify credentials in `.env` file
- Test connection: `mysql -u weather_user -p -h localhost weather_label_db`

### "Cannot upload images"
- Check `uploads/` directory exists
- Verify directory permissions: `chmod 755 uploads/`
- Check disk space: `df -h`

### "Frontend not loading"
- Ensure frontend is built: `cd frontend && npm run build`
- Check browser console for errors (F12)
- Verify server is running on port 8080

### "Port already in use"
- Change PORT in `.env` file
- Or stop the process using port 8080: `lsof -ti:8080 | xargs kill`

## Development Mode

For active development with hot reload:

### Terminal 1 - Backend
```bash
go run main.go
```

### Terminal 2 - Frontend
```bash
cd frontend
npm run dev
```

Frontend dev server runs on http://localhost:5173 with hot reload.

## Next Steps

- Read [README.md](README.md) for detailed documentation
- Check [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment
- Use [Makefile](Makefile) commands for common tasks

## Getting Help

If you encounter issues:

1. Check logs: `journalctl -u weather-label` (if using systemd)
2. Check MySQL logs: `tail -f /var/log/mysql/error.log`
3. Enable debug mode: Add `set -x` to scripts
4. Create an issue on GitHub with:
   - Your OS and versions (Go, Node.js, MySQL)
   - Error messages
   - Steps to reproduce

## Common Commands

```bash
# Using Makefile
make help          # Show all available commands
make setup         # Initial setup
make build         # Build application
make run           # Build and run
make clean         # Clean build artifacts
make docker-up     # Start MySQL with Docker
make docker-down   # Stop MySQL

# Manual commands
./setup.sh         # Automated setup
./dev.sh           # Development mode
go run main.go     # Run backend
./server           # Run built server
```

Enjoy using the Weather Agency Label Tool! 🌤️
