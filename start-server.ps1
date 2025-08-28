#!/usr/bin/env pwsh

# Script để chạy backend server
Write-Host "Starting Backend Server..." -ForegroundColor Green

# Kiểm tra port 3000 có bị sử dụng không
$portCheck = netstat -ano | findstr :3000
if ($portCheck) {
    Write-Host "Port 3000 is already in use:" -ForegroundColor Yellow
    Write-Host $portCheck
    $continue = Read-Host "Do you want to continue anyway? (y/N)"
    if ($continue -ne "y" -and $continue -ne "Y") {
        exit 1
    }
}

# Build và chạy application
Write-Host "Building application..." -ForegroundColor Cyan
go build -o bin/main.exe cmd/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Server will be available at: http://localhost:3000" -ForegroundColor Cyan
    Write-Host "Health check: http://localhost:3000/health" -ForegroundColor Cyan
    Write-Host "API endpoints:" -ForegroundColor Cyan
    Write-Host "   POST /api/auth/register" -ForegroundColor Gray
    Write-Host "   POST /api/auth/login" -ForegroundColor Gray
    Write-Host "   GET  /api/auth/profile" -ForegroundColor Gray
    Write-Host "   GET  /api/admin/users" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Yellow
    Write-Host ""
    
    # Chạy server
    .\bin\main.exe
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
