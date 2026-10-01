# AI Generate API - Backend

REST API untuk AI text generation menggunakan HuggingFace models dengan autentikasi JWT dan API token management.

## 🚀 Tech Stack

- **Go 1.26+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM untuk PostgreSQL
- **PostgreSQL** - Database
- **JWT** - Authentication
- **HuggingFace Router API** - AI model provider (OpenAI-compatible)

## 📁 Struktur Project

```
.
├── config/         # Konfigurasi database dan environment
├── handlers/       # HTTP request handlers
├── middleware/     # JWT auth, API key auth, admin middleware
├── models/         # Database models (User, APIToken, AIModel, Config)
├── routes/         # Route definitions
├── seeders/        # Database seeders (admin, models, configs)
├── utils/          # Helper functions (password hash, token generation)
├── main.go         # Entry point aplikasi
├── .env            # Environment variables (jangan commit!)
└── .env.example    # Template environment variables
```

## 🔧 Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd huggingface-ai
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Database

Buat database PostgreSQL atau gunakan service seperti Prisma.io:

```bash
# Contoh menggunakan PostgreSQL lokal
createdb ai_generate_db
```

### 4. Konfigurasi Environment

Salin file `.env.example` ke `.env` dan sesuaikan:

```bash
cp .env.example .env
```

Edit file `.env`:

```env
# App
APP_PORT=8080
APP_ENV=development

# Database
DATABASE_URL=postgres://user:password@host:5432/dbname?sslmode=require

# JWT - WAJIB diisi
JWT_SECRET=your-super-secret-key-change-this
JWT_EXPIRED_HOURS=24

# Admin default - untuk seeder
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123

# HuggingFace API
HF_API_KEY=hf_xxx
HF_MODEL=zai-org/GLM-5.3:novita
HF_BASE_URL=https://router.huggingface.co/v1
```

### 5. Run Server

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `POST /api/register` - Register user baru
- `POST /api/login` - Login dan dapatkan JWT token
- `GET /v1/models` - List AI models (OpenAI-compatible)
- `GET /models` - List AI models (alias)

### Protected Endpoints (JWT Required)

- `GET /api/me` - Get user profile
- `GET /api/tokens` - List API tokens milik user
- `POST /api/tokens` - Create API token baru
- `GET /api/tokens/:id` - Detail API token
- `PUT /api/tokens/:id/status` - Toggle status API token
- `DELETE /api/tokens/:id` - Delete API token

### Admin Endpoints (JWT + Admin Role Required)

- `GET /api/admin/models` - List semua AI models
- `POST /api/admin/models` - Create AI model baru
- `GET /api/admin/models/:id` - Detail AI model
- `PUT /api/admin/models/:id` - Update AI model
- `DELETE /api/admin/models/:id` - Delete AI model
- `GET /api/admin/configs` - List configs
- `POST /api/admin/configs` - Create config
- `GET /api/admin/configs/:id` - Detail config
- `PUT /api/admin/configs/:id` - Update config
- `DELETE /api/admin/configs/:id` - Delete config

### AI Generation Endpoints (API Token Required)

- `POST /v1/messages` - Generate text (custom format)
- `POST /v1/chat/completions` - Generate text (OpenAI format)
- `POST /chat/completions` - Generate text (OpenAI format, no prefix)

## 🔑 Authentication

### 1. JWT Authentication (untuk akses user endpoints)

Login untuk mendapatkan JWT token:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "full_name": "Administrator",
    "username": "admin",
    "role": "admin"
  }
}
```

Gunakan token untuk request berikutnya:

```bash
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 2. API Token Authentication (untuk AI generation)

Buat API token melalui endpoint `/api/tokens`, lalu gunakan untuk AI generation:

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxx" \
  -d '{
    "model": "zai-org/GLM-5.3:novita",
    "messages": [{"role":"user","content":"Hello"}],
    "max_tokens": 100
  }'
```

## 🗄️ Database Models

### User
- `id` (uint) - Primary key
- `full_name` (string) - Nama lengkap
- `username` (string) - Username (unique)
- `password` (string) - Hashed password
- `role` (string) - "user" atau "admin"
- `created_at`, `updated_at`

### APIToken
- `id` (uint) - Primary key
- `user_id` (uint) - Foreign key ke users
- `token` (string) - API token (unique, prefix: sk-)
- `name` (string) - Nama token
- `is_active` (bool) - Status aktif/nonaktif
- `created_at`, `updated_at`

### AIModel
- `id` (uint) - Primary key
- `name` (string) - Nama model
- `slug` (string) - Model slug (unique)
- `description` (string) - Deskripsi model
- `is_active` (bool) - Status aktif/nonaktif
- `created_at`, `updated_at`

### Config
- `id` (uint) - Primary key
- `key` (string) - Config key (unique)
- `value` (string) - Config value
- `description` (string) - Deskripsi config
- `created_at`, `updated_at`

## 🌱 Seeders

Seeders berjalan otomatis saat server start:

1. **Admin User Seeder** - Buat user admin dari `ADMIN_USERNAME` dan `ADMIN_PASSWORD`
2. **Model Seeder** - Buat model default dari `HF_MODEL`
3. **Config Seeder** - Buat config HF_API_KEY, HF_MODEL, HF_BASE_URL

## 🔒 Security

- Password di-hash menggunakan bcrypt
- JWT token expire sesuai `JWT_EXPIRED_HOURS`
- API token diawali prefix `sk-` dan random 32 karakter
- Middleware admin hanya izinkan role "admin"
- CORS dikonfigurasi untuk frontend `http://localhost:5173`

## 🧪 Testing

Login sebagai admin:

```bash
# Register user baru
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Test User","username":"test","password":"test123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Get profile
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer <jwt-token>"
```

## 📝 Development

```bash
# Run dengan auto-reload (install air)
go install github.com/cosmtrek/air@latest
air

# Build binary
go build -o ai-generate-api

# Run binary
./ai-generate-api
```

## 📄 License

MIT License

## 👨‍💻 Author

Andri Dev Code
