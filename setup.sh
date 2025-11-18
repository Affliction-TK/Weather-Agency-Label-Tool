#!/bin/bash

# Weather Agency Label Tool - Setup Script

echo "================================"
echo "Weather Agency Label Tool Setup"
echo "================================"
echo ""

# Check prerequisites
echo "Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24 or higher."
    exit 1
fi
echo "✅ Go installed: $(go version)"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 20 or higher."
    exit 1
fi
echo "✅ Node.js installed: $(node --version)"

# Check npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed."
    exit 1
fi
echo "✅ npm installed: $(npm --version)"

# Check MySQL
if ! command -v mysql &> /dev/null; then
    echo "⚠️  MySQL client not found. You can use docker-compose to run MySQL."
    echo "   Run: docker-compose up -d"
fi

echo ""
echo "Step 1: Installing Go dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "❌ Failed to download Go dependencies"
    exit 1
fi
echo "✅ Go dependencies installed"

echo ""
echo "Step 2: Installing frontend dependencies..."
cd frontend
npm install
if [ $? -ne 0 ]; then
    echo "❌ Failed to install npm dependencies"
    exit 1
fi
echo "✅ Frontend dependencies installed"

echo ""
echo "Step 3: Building frontend..."
npm run build
if [ $? -ne 0 ]; then
    echo "❌ Failed to build frontend"
    exit 1
fi
echo "✅ Frontend built successfully"

cd ..

echo ""
echo "Step 4: Creating uploads directory..."
mkdir -p uploads
echo "✅ Uploads directory created"

echo ""
echo "Step 5: Creating .env file from example..."
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ .env file created. Please edit it with your database credentials."
else
    echo "⚠️  .env file already exists, skipping..."
fi

echo ""
echo "================================"
echo "Setup completed successfully! 🎉"
echo "================================"
echo ""
echo "Next steps:"
echo "1. Start MySQL database:"
echo "   - Option A: docker-compose up -d"
echo "   - Option B: Use your existing MySQL installation"
echo ""
echo "2. Initialize the database:"
echo "   mysql -u root -p < schema.sql"
echo ""
echo "3. Update .env file with your database credentials"
echo ""
echo "4. Run the application:"
echo "   go run main.go"
echo ""
echo "5. Open http://localhost:8080 in your browser"
echo ""
