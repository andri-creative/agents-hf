# TUGAS
Buatkan saya sebuah REST API lengkap, production-ready, menggunakan bahasa 
pemrograman Go (Golang). API ini berfungsi sebagai gateway untuk AI Generation 
yang terintegrasi dengan HuggingFace Inference API dan database PostgreSQL.

---

## 1. TEKNOLOGI YANG DIGUNAKAN (WAJIB)

| Komponen       | Library / Tools                                              |
|----------------|--------------------------------------------------------------|
| Web Framework  | github.com/gin-gonic/gin                                     |
| Database       | PostgreSQL + gorm.io/gorm + gorm.io/driver/postgres          |
| Auth (JWT)     | github.com/golang-jwt/jwt/v5                                 |
| Password Hash  | golang.org/x/crypto/bcrypt                                   |
| Validasi       | github.com/go-playground/validator/v10                       |
| ENV Loader     | github.com/joho/godotenv                                     |
| AI Provider    | HuggingFace Inference API (HTTP REST)                        |

---

## 2. STRUKTUR FOLDER PROJECT (WAJIB DIIKUTI)

project/
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
├── config/
│   ├── env.go
│   └── database.go
├── models/
│   └── user.go
├── dto/
│   ├── auth_dto.go
│   └── ai_dto.go
├── handlers/
│   ├── auth_handler.go
│   └── ai_handler.go
├── middleware/
│   └── auth_middleware.go
├── routes/
│   └── routes.go
├── services/
│   └── huggingface.go
└── utils/
    ├── jwt.go
    ├── hash.go
    └── response.go

Setiap folder HARUS dijelaskan fungsinya dalam komentar singkat di atas file.

---

## 3. FITUR YANG DIMINTA

### 🔹 FITUR 1 — REGISTER
- **Endpoint**: `POST /api/register`
- **Request Body (JSON)**:
  ```json
  {
    "full_name": "Budi Santoso",
    "username": "budi123",
    "password": "rahasia123"
  }
  ```
- **Behavior**:
  1. Validasi input (full_name min 3, username alphanumeric min 3, password min 6)
  2. Cek apakah username sudah terpakai
  3. Hash password dengan bcrypt
  4. **Generate API Token unik otomatis** dengan format: `hf_` + 64 karakter hex random
  5. Simpan user baru ke database (termasuk API Token)
  6. Return: `api_token` + data user (tanpa password)

---

### 🔹 FITUR 2 — LOGIN
- **Endpoint**: `POST /api/login`
- **Request Body (JSON)**:
  ```json
  {
    "username": "budi123",
    "password": "rahasia123"
  }
  ```
- **Behavior**:
  1. Cari user berdasarkan username
  2. Verifikasi password dengan bcrypt
  3. Generate JWT token (expired dari ENV, default 24 jam)
  4. Return: `jwt_token` + `api_token` + data user

---

### 🔹 FITUR 3 — PROFILE
- **Endpoint**: `GET /api/me`
- **Protected**: Ya, menggunakan JWT middleware
- **Behavior**: Return data user yang sedang login (berdasarkan `user_id` di JWT claims)

---

### 🔹 FITUR 4 — AI GENERATE (ENDPOINT UTAMA)
- **Endpoint**: `POST /v1/messages`
- **Protected**: Ya, menggunakan API Token middleware
- **Request Body (JSON)**:
  ```json
  {
    "model": "mistralai/Mistral-7B-Instruct-v0.2",
    "prompt": "Jelaskan apa itu Go dalam 1 paragraf"
  }
  ```
- **Behavior**:
  1. Validasi input
  2. Jika `model` kosong → fallback ke `HF_MODEL` dari ENV
  3. Forward request ke HuggingFace Inference API
  4. Parse response dan ambil field `generated_text`
  5. Return: `{ model, response }`

---

### 🔹 FITUR 5 — HEALTH CHECK
- **Endpoint**: `GET /health`
- **Behavior**: Return `{ "status": "ok" }` untuk cek server hidup

---

## 4. ATURAN & KETENTUAN (STRICT)

### 🔐 Keamanan
- Semua credential (DB, JWT_SECRET, HF_API_KEY, HF_MODEL) **WAJIB** disimpan di file `.env`
- Password **TIDAK BOLEH** muncul di JSON response (`json:"-"`)
- API Token di-generate menggunakan `crypto/rand`, **JANGAN** pakai `math/rand`
- JWT menggunakan algoritma `HS256` dan divalidasi method signingnya

### 🧪 Validasi
- Gunakan `go-playground/validator/v10` dengan tag `validate:"..."` di struct DTO
- Setiap field di struct DTO HARUS punya tag `json:"..."` dan `validate:"..."`
- Return error 400 jika validasi gagal, dengan detail field yang salah

### 📦 Format Response (KONSISTEN)
**Sukses:**
```json
{
  "success": true,
  "message": "Pesan sukses",
  "data": { ... }
}
```
**Gagal:**
```json
{
  "success": false,
  "message": "Pesan error",
  "error": "Detail error"
}
```

### 🌐 HTTP Status Code
| Kondisi                      | Status |
|------------------------------|--------|
| Sukses GET                   | 200    |
| Sukses POST (create)         | 201    |
| Bad Request / Validasi gagal | 400    |
| Unauthorized                 | 401    |
| Not Found                    | 404    |
| Conflict (username terpakai) | 409    |
| Internal Server Error        | 500    |
| Bad Gateway (HF error)       | 502    |

### 🧱 Arsitektur
- Gunakan **package separation** yang jelas (config, models, dto, handlers, middleware, routes, services, utils)
- Handler **HANYA** handle HTTP request/response — logic berat ditaruh di service
- Tidak ada **global variable** berlebihan, kecuali `DB` dan `Env`
- Setiap fungsi HARUS punya error handling yang jelas
- Berikan **komentar** pada bagian penting agar mudah dipahami

---

## 5. MODEL USER (DATABASE SCHEMA)

Tabel: `users`

| Kolom       | Tipe          | Constraint              |
|-------------|---------------|-------------------------|
| id          | uint          | primary key, auto inc   |
| full_name   | varchar(120)  | not null                |
| username    | varchar(60)   | unique index, not null  |
| password    | varchar(255)  | not null                |
| api_token   | varchar(64)   | unique index            |
| created_at  | timestamp     | auto                    |
| updated_at  | timestamp     | auto                    |

**Auto migrate** tabel User saat server start menggunakan `config.DB.AutoMigrate(&models.User{})`.

---

## 6. HUGGINGFACE INTEGRATION (DETAIL TEKNIS)

**Base URL**: `{HF_BASE_URL}/models/{model}` (default: `https://api-inference.huggingface.co`)

**Headers**:
```
Authorization: Bearer {HF_API_KEY}
Content-Type: application/json
```

**Request Body**:
```json
{
  "inputs": "prompt dari user",
  "parameters": {
    "max_new_tokens": 512,
    "temperature": 0.7,
    "return_full_text": false,
    "options": {
      "wait_for_model": true,
      "use_cache": false
    }
  }
}
```

**Behavior**:
- Timeout HTTP client: **60 detik**
- Response dari HF berbentuk array → ambil `result[0].generated_text`
- Jika status != 200 → return error dengan detail dari HF

---

## 7. ENV FILE (WAJIB ADA .env.example)

```env
# App
APP_PORT=8080
APP_ENV=development

# Database PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=ai_generate_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=super-secret-key-change-this
JWT_EXPIRED_HOURS=24

# HuggingFace
HF_API_KEY=hf_xxxxxxxxxxxxxxxxxxxxx
HF_MODEL=mistralai/Mistral-7B-Instruct-v0.2
HF_BASE_URL=https://api-inference.huggingface.co
```

Setiap variabel HARUS dibaca di `config/env.go` dengan fallback default yang wajar.
Jika `JWT_SECRET` kosong → **server harus gagal start** (log.Fatal).

---

## 8. MIDDLEWARE (WAJIB 2 BUAH)

### Middleware 1: `JWTAuth()`
- Cek header `Authorization: Bearer <token>`
- Parse JWT dengan secret dari ENV
- Simpan `user_id` dan `username` ke context
- Return 401 jika token invalid/expired

### Middleware 2: `APIKeyAuth()`
- Cek header `X-API-Key` ATAU `Authorization: Bearer <api_token>`
- Cari user di DB dengan `api_token` tersebut
- Simpan `user_id` dan `username` ke context
- Return 401 jika API key tidak valid
- **Khusus dipakai** di endpoint `/v1/messages`

---

## 9. OUTPUT YANG DIHARAPKAN

Tolong hasilkan **SEMUA FILE LENGKAP** dengan urutan berikut:

1. `.env.example` (dengan komentar)
2. `go.mod` (dengan daftar dependency)
3. `main.go`
4. `config/env.go`
5. `config/database.go`
6. `models/user.go`
7. `dto/auth_dto.go`
8. `dto/ai_dto.go`
9. `utils/response.go`
10. `utils/hash.go`
11. `utils/jwt.go`
12. `middleware/auth_middleware.go`
13. `handlers/auth_handler.go`
14. `handlers/ai_handler.go`
15. `services/huggingface.go`
16. `routes/routes.go`

Untuk setiap file:
- Berikan **header komentar** berisi nama file dan fungsinya
- Sertakan **import** yang lengkap
- Code **HARUS BISA LANGSUNG DIJALANKAN** (compilable, no pseudo-code)

---

## 10. CONTOH CARA PAKAI (SERTAKAN DI AKHIR)

Sertakan contoh **curl** untuk:

### a) Register
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Budi Santoso",
    "username": "budi123",
    "password": "rahasia123"
  }'
```

### b) Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"budi123","password":"rahasia123"}'
```

### c) Generate AI
```bash
curl -X POST http://localhost:8080/v1/messages \
  -H "Content-Type: application/json" \
  -H "X-API-Key: hf_xxxxxxxx" \
  -d '{
    "model": "mistralai/Mistral-7B-Instruct-v0.2",
    "prompt": "Jelaskan apa itu Go"
  }'
```

---

## 11. CATATAN AKHIR

- **JANGAN** pakai pseudo-code — semua harus real code.
- **JANGAN** skip file — hasilkan semuanya dari nomor 1 sampai 16.
- Gunakan **best practice Go**: error wrapping, konteks, package separation.
- Tambahkan **komentar** di bagian penting untuk memudahkan pemahaman.
- Pastikan `go.mod` module name: `ai-generate-api`.
- Jika ada ambiguitas, pilih pendekatan yang **paling aman & production-ready**.

---

## ✅ CHECKLIST AKHIR (VERIFIKASI SEBELUM SELESAI)

- [ ] Semua 16 file dihasilkan lengkap
- [ ] Password di-hide dari JSON response
- [ ] API Token di-generate pakai `crypto/rand`
- [ ] Middleware JWT & API Key terpisah
- [ ] Validasi input dengan validator/v10
- [ ] Format response konsisten (success/message/data)
- [ ] HTTP status code sesuai tabel
- [ ] Auto migrate tabel User
- [ ] ENV dibaca dari `.env` dengan fallback
- [ ] Contoh curl disertakan
- [ ] Code compilable tanpa error