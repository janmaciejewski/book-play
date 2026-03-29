# BookPlay Full App Startup Script for Windows
# Run: ./start.ps1
# This starts everything: Database, Backend API, and Frontend Web

Write-Host "🚀 Starting BookPlay Full Application..." -ForegroundColor Cyan
Write-Host ""

# Check if Docker is running
Write-Host "Checking Docker..." -ForegroundColor Yellow
$dockerRunning = docker ps 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Docker is not running. Please start Docker Desktop first." -ForegroundColor Red
    exit 1
}

# Start database services
Write-Host "Starting PostgreSQL and Redis..." -ForegroundColor Yellow
docker-compose up -d postgres redis

# Wait for PostgreSQL to be ready
Write-Host "Waiting for PostgreSQL to be ready..." -ForegroundColor Yellow
$maxAttempts = 30
$attempt = 0
do {
    $attempt++
    $result = docker exec bookplay-postgres pg_isready -U bookplay 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ PostgreSQL is ready!" -ForegroundColor Green
        break
    }
    if ($attempt -ge $maxAttempts) {
        Write-Host "❌ PostgreSQL failed to start within 30 seconds" -ForegroundColor Red
        exit 1
    }
    Start-Sleep -Seconds 1
    Write-Host "." -NoNewline
} while ($true)

# Wait for Redis to be ready
Write-Host "Waiting for Redis to be ready..." -ForegroundColor Yellow
$attempt = 0
do {
    $attempt++
    $result = docker exec bookplay-redis redis-cli ping 2>$null
    if ($result -eq "PONG") {
        Write-Host "✅ Redis is ready!" -ForegroundColor Green
        break
    }
    if ($attempt -ge $maxAttempts) {
        Write-Host "❌ Redis failed to start within 30 seconds" -ForegroundColor Red
        exit 1
    }
    Start-Sleep -Seconds 1
    Write-Host "." -NoNewline
} while ($true)

Write-Host ""
Write-Host "✅ Database services are running!" -ForegroundColor Green
Write-Host ""

# Install frontend dependencies if needed
if (-not (Test-Path "apps/web/node_modules")) {
    Write-Host "Installing frontend dependencies..." -ForegroundColor Yellow
    npm install
}

Write-Host "Starting Backend (port 8080) and Frontend (port 3000)..." -ForegroundColor Cyan
Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  🌐 Frontend:  http://localhost:3000" -ForegroundColor White
Write-Host "  🔌 Backend:   http://localhost:8080" -ForegroundColor White
Write-Host "  📊 Database:  localhost:5432" -ForegroundColor White
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press Ctrl+C to stop all services" -ForegroundColor Yellow
Write-Host ""

# Run both backend and frontend concurrently with hot reload
# Backend: air (hot reload for Go)
# Frontend: npm run dev:web (hot reload for Nuxt)
npx concurrently --names "Backend,Frontend" --prefix-colors "magenta,cyan" "air" "npm run dev:web"
