# Scripts PowerShell para testar a API Ativvo

# 1. Health Check
function Test-Health {
    Write-Host "🏥 Testing Health Check..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri "http://localhost:8080/api/health" -Method GET | 
        Select-Object -ExpandProperty Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
}

# 2. Register User
function Register-User {
    param(
        [string]$Email = "test@example.com",
        [string]$Password = "senha123",
        [string]$Name = "Test User"
    )
    
    Write-Host "📝 Registering new user..." -ForegroundColor Cyan
    $body = @{
        email = $Email
        password = $Password
        name = $Name
    } | ConvertTo-Json
    
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/auth/register" `
        -Method POST `
        -ContentType "application/json" `
        -Body $body
    
    $response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
}

# 3. Login
function Login-User {
    param(
        [string]$Email = "test@example.com",
        [string]$Password = "senha123"
    )
    
    Write-Host "🔐 Logging in..." -ForegroundColor Cyan
    $body = @{
        email = $Email
        password = $Password
    } | ConvertTo-Json
    
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/auth/login" `
        -Method POST `
        -ContentType "application/json" `
        -Body $body
    
    $result = $response.Content | ConvertFrom-Json
    $global:Token = $result.token
    Write-Host "✅ Token saved to `$global:Token" -ForegroundColor Green
    $result | ConvertTo-Json -Depth 10
}

# 4. Get My Profile
function Get-MyProfile {
    if (-not $global:Token) {
        Write-Host "❌ You need to login first! Run: Login-User" -ForegroundColor Red
        return
    }
    
    Write-Host "👤 Getting profile..." -ForegroundColor Cyan
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/me" `
        -Method GET `
        -Headers @{
            "Authorization" = "Bearer $global:Token"
        }
    
    $response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
}

# 5. Update Profile
function Update-Profile {
    param(
        [string]$Name,
        [string]$Gender,
        [double]$Height,
        [double]$Weight,
        [double]$BodyFat,
        [int]$WeeklyWorkouts,
        [int]$CardioTime,
        [string]$Goal
    )
    
    if (-not $global:Token) {
        Write-Host "❌ You need to login first! Run: Login-User" -ForegroundColor Red
        return
    }
    
    Write-Host "✏️ Updating profile..." -ForegroundColor Cyan
    
    $body = @{}
    if ($Name) { $body.name = $Name }
    if ($Gender) { $body.gender = $Gender }
    if ($Height) { $body.height = $Height }
    if ($Weight) { $body.weight = $Weight }
    if ($BodyFat) { $body.body_fat = $BodyFat }
    if ($WeeklyWorkouts) { $body.weekly_workouts = $WeeklyWorkouts }
    if ($CardioTime) { $body.cardio_time = $CardioTime }
    if ($Goal) { $body.goal = $Goal }
    
    $jsonBody = $body | ConvertTo-Json
    
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/profile" `
        -Method PUT `
        -ContentType "application/json" `
        -Headers @{
            "Authorization" = "Bearer $global:Token"
        } `
        -Body $jsonBody
    
    $response.Content | ConvertFrom-Json | ConvertTo-Json -Depth 10
}

# Exemplos de uso
Write-Host @"

╔════════════════════════════════════════════════════╗
║           🏋️ Ativvo API Test Scripts 🏋️            ║
╚════════════════════════════════════════════════════╝

Comandos disponíveis:

1️⃣  Test-Health
    Testa o health check da API

2️⃣  Register-User -Email "seu@email.com" -Password "senha" -Name "Seu Nome"
    Registra um novo usuário

3️⃣  Login-User -Email "seu@email.com" -Password "senha"
    Faz login e salva o token

4️⃣  Get-MyProfile
    Busca seus dados (precisa estar logado)

5️⃣  Update-Profile -Height 175 -Weight 70 -Goal "gain_muscle"
    Atualiza seu perfil

Exemplo completo:
-----------------
Register-User -Email "victor@ativvo.com" -Password "123456" -Name "Victor"
Login-User -Email "victor@ativvo.com" -Password "123456"
Get-MyProfile
Update-Profile -Height 180 -Weight 75 -WeeklyWorkouts 5 -Goal "gain_muscle"
Get-MyProfile

"@ -ForegroundColor Yellow
