# PowerShell script for Backend API management

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

switch ($Command.ToLower()) {
    "build" {
        Write-Host "🔨 Building application..." -ForegroundColor Yellow
        go build -o bin/main.exe cmd/main.go
        go build -o bin/setup.exe cmd/setup/main.go
        Write-Host "✅ Build completed!" -ForegroundColor Green
    }
    
    "run" {
        Write-Host "🚀 Starting application..." -ForegroundColor Yellow
        go run cmd/main.go
    }
    
    "setup" {
        Write-Host "🔧 Setting up database..." -ForegroundColor Yellow
        go run cmd/setup/main.go
    }
    
    "clean" {
        Write-Host "🧹 Cleaning build artifacts..." -ForegroundColor Yellow
        if (Test-Path "bin") { Remove-Item -Recurse -Force "bin" }
        if (Test-Path "tmp") { Remove-Item -Recurse -Force "tmp" }
        Write-Host "✅ Clean completed!" -ForegroundColor Green
    }
    
    "test" {
        Write-Host "🧪 Running tests..." -ForegroundColor Yellow
        go test ./...
    }
    
    "deps" {
        Write-Host "📦 Installing dependencies..." -ForegroundColor Yellow
        go mod download
        go mod tidy
        Write-Host "✅ Dependencies installed!" -ForegroundColor Green
    }
    
    "createdb" {
        Write-Host "🗄️ Creating database..." -ForegroundColor Yellow
        createdb backend
        Write-Host "✅ Database created!" -ForegroundColor Green
    }
    
    "help" {
        Write-Host "Backend API Management Script" -ForegroundColor Cyan
        Write-Host "=============================" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Available commands:" -ForegroundColor White
        Write-Host "  build    - Build the application" -ForegroundColor Gray
        Write-Host "  run      - Run the application" -ForegroundColor Gray
        Write-Host "  setup    - Setup database and create admin user" -ForegroundColor Gray
        Write-Host "  clean    - Clean build artifacts" -ForegroundColor Gray
        Write-Host "  test     - Run tests" -ForegroundColor Gray
        Write-Host "  deps     - Install dependencies" -ForegroundColor Gray
        Write-Host "  createdb - Create database" -ForegroundColor Gray
        Write-Host "  help     - Show this help" -ForegroundColor Gray
        Write-Host ""
        Write-Host "Usage: .\manage.ps1 <command>" -ForegroundColor Yellow
        Write-Host "Example: .\manage.ps1 run" -ForegroundColor Yellow
    }
    
    default {
        Write-Host "❌ Unknown command: $Command" -ForegroundColor Red
        Write-Host "Run '.\manage.ps1 help' to see available commands." -ForegroundColor Yellow
    }
}
