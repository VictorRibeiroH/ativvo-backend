# Ativvo Backend

Backend da aplicação Ativvo - Sistema de gerenciamento de treinos e fitness.

## 🚀 Tecnologias

- **Go 1.25+** - Linguagem principal
- **Fiber v2** - Framework web (rápido e moderno)
- **GORM** - ORM para PostgreSQL
- **PostgreSQL** - Banco de dados (via Supabase)
- **JWT** - Autenticação com tokens
- **bcrypt** - Hash seguro de senhas
- **UUID** - Identificadores únicos para entidades

## 📁 Estrutura

```
ativvo-backend/
├── main.go                 # Entrada da aplicação
├── go.mod                  # Dependências
├── .env                    # Variáveis de ambiente (não commitado)
├── .env.example            # Template de variáveis
└── internal/
    ├── config/            # Configurações e env vars
    ├── database/          # Conexão e migrações
    ├── models/            # Modelos de dados (User, Workout)
    ├── handlers/          # Controllers/Handlers
    └── middleware/        # Middlewares (auth, etc)
```

## ⚙️ Setup

### 1. Clonar e instalar dependências

```bash
cd ativvo-backend
go mod download
```

### 2. Configurar variáveis de ambiente

Copie `.env.example` para `.env` e preencha com suas credenciais:

```bash
cp .env.example .env
```

**Obtenha a Database URL do Supabase:**
1. Acesse https://supabase.com
2. Vá em **Settings > Database**
3. Copie a **Connection string (URI)**
4. Cole no `.env` em `DATABASE_URL`

**Gere um JWT Secret forte:**
```bash
# No PowerShell (Windows)
-join ((48..57) + (65..90) + (97..122) | Get-Random -Count 64 | ForEach-Object {[char]$_})
```

Cole o resultado em `JWT_SECRET` no `.env`.

### 3. Rodar o servidor

```bash
go run main.go
```

O servidor estará rodando em `http://localhost:8080` 🎉

## 📡 Endpoints

### Públicos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/health` | Health check |
| POST | `/api/auth/register` | Criar nova conta |
| POST | `/api/auth/login` | Fazer login |

### Protegidos (requer token JWT)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/me` | Dados do usuário logado |
| PUT | `/api/profile` | Atualizar perfil |

## 🔐 Autenticação

### Registro

```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "usuario@exemplo.com",
  "password": "senha123",
  "name": "João Silva"
}
```

**Resposta:**
```json
{
  "message": "User created successfully",
  "user": {
    "id": "uuid-aqui",
    "email": "usuario@exemplo.com",
    "name": "João Silva"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Login

```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "usuario@exemplo.com",
  "password": "senha123"
}
```

### Usando o token

Para rotas protegidas, envie o token no header:

```bash
GET /api/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 🗄️ Models

### User

```go
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "gender": "string",        // male, female, other
  "height": float64,         // cm
  "weight": float64,         // kg
  "body_fat": float64,       // %
  "weekly_workouts": int,    // treinos/semana
  "cardio_time": int,        // minutos
  "goal": "string",          // lose_weight, gain_muscle, maintain
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Workout (futuro)

```go
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "string",
  "description": "string",
  "type": "string",          // strength, cardio, flexibility
  "duration": int,           // minutos
  "calories": int,
  "date": "timestamp"
}
```

## 🛠️ Desenvolvimento

### Rodar em modo dev

```bash
go run main.go
```

### Build para produção

```bash
go build -o ativvo-backend.exe
./ativvo-backend.exe
```

### Executar testes

```bash
go test ./...
```

## 🔒 Segurança

- ✅ Senhas hasheadas com bcrypt
- ✅ JWT para autenticação stateless
- ✅ Variáveis sensíveis em `.env` (não commitadas)
- ✅ CORS configurado apenas para frontend autorizado
- ✅ Soft delete em models (dados não são deletados permanentemente)
- ✅ Validação de dados com `validator`

## 📝 Próximos passos

- [ ] Endpoints de Workouts (CRUD)
- [ ] Estatísticas e gráficos
- [ ] Upload de imagens de perfil
- [ ] Reset de senha
- [ ] Refresh tokens
- [ ] Rate limiting
- [ ] Testes unitários e de integração

## 📄 Licença

MIT

---

Desenvolvido com ❤️ usando Go e Fiber
