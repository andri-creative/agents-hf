# TUGAS: UPGRADE PROJECT GO API YANG SUDAH ADA

Saya SUDAH PUNYA project Go API yang berjalan dengan struktur dan fitur berikut:

## KONDISI PROJECT SAAT INI (JANGAN DIUBAH STRUKTURNYA)

```
project/
├── .env
├── .env.example
├── go.mod
├── go.sum
├── main.go
├── config/
│   ├── env.go
│   └── database.go
├── models/
│   └── user.go          → struct User { ID, FullName, Username, Password, APIToken, CreatedAt, UpdatedAt }
├── dto/
│   ├── auth_dto.go
│   └── ai_dto.go
├── handlers/
│   ├── auth_handler.go  → Register, Login, Me
│   └── ai_handler.go    → GenerateMessage
├── middleware/
│   └── auth_middleware.go → JWTAuth(), APIKeyAuth()
├── routes/
│   └── routes.go
├── services/
│   └── huggingface.go
└── utils/
    ├── jwt.go
    ├── hash.go
    └── response.go
```

**Fitur yang sudah jalan saat ini:**
- Register (full_name, username, password) → auto-generate 1 API token di kolom `users.api_token`
- Login → return JWT + API token (yang lama, dari kolom `users.api_token`)
- Endpoint `POST /v1/messages` — diproteksi API key (cek kolom `users.api_token`)
- HuggingFace config dibaca dari `.env` (HF_API_KEY, HF_MODEL, HF_BASE_URL)
- Model HF hardcode / dari ENV saja (tidak ada tabel model)

---

## 🎯 TUGAS SAYA: TAMBAHKAN 3 UPGRADE BERIKUT TANPA REWRITE DARI NOL

Saya ingin kamu **MEMODIFIKASI & MENAMBAHKAN** file pada project yang sudah ada. 
**JANGAN** buat project baru. **JANGAN** ubah struktur folder yang sudah ada kecuali 
menambah file baru. Semua file lama yang tidak disebut → biarkan apa adanya.

Berikut 3 upgrade yang saya minta:

---

### 🆙 UPGRADE 1 — PISAHKAN API TOKEN KE TABEL `api_tokens` (RELASI)

**Tujuan**: 
- API token **dipisah** dari tabel `users` ke tabel baru `api_tokens`
- 1 user bisa punya **BANYAK token** (multi-device)
- Login **TIDAK menimpa** token lama
- Setiap login berhasil → **generate token BARU** di tabel `api_tokens`

#### 🔧 Yang harus dilakukan:

**1. Buat file baru: `models/api_token.go`**
```go
type APIToken struct {
    ID         uint       `gorm:"primaryKey" json:"id"`
    UserID     uint       `gorm:"index;not null" json:"user_id"`
    Token      string     `gorm:"type:varchar(80);uniqueIndex;not null" json:"token"`
    Name       string     `gorm:"type:varchar(100)" json:"name"`
    IsActive   bool       `gorm:"default:true" json:"is_active"`
    LastUsedAt *time.Time `json:"last_used_at"`
    ExpiresAt  *time.Time `json:"expires_at"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
}
```

**2. Modifikasi `models/user.go`**
- **HAPUS** kolom `APIToken` dari struct User
- **TAMBAHKAN** relasi:
  ```go
  APITokens []APIToken `gorm:"foreignKey:UserID" json:"api_tokens,omitempty"`
  ```

**3. Modifikasi `handlers/auth_handler.go`**
- **Register**: 
  - Buat user dulu → lalu buat 1 token di tabel `api_tokens` (name: "Default Token")
  - Return data user + token baru
- **Login**:
  - Verifikasi user
  - **Generate token BARU** setiap login (format: `hf_` + 64 hex)
  - Simpan ke tabel `api_tokens` (name: "Login dari <timestamp>")
  - **JANGAN hapus** token lama
  - Return: `jwt_token` + `new_api_token` + data user

**4. Buat file baru: `dto/token_dto.go`**
```go
type CreateTokenRequest struct {
    Name      string     `json:"name" validate:"required,min=1,max=100"`
    ExpiresAt *time.Time `json:"expires_at"`
}

type UpdateTokenStatusRequest struct {
    IsActive bool `json:"is_active"`
}
```

**5. Buat file baru: `handlers/token_handler.go`**
Endpoint baru (protected JWT):
- `GET    /api/tokens`              → list semua token milik user yang login
- `POST   /api/tokens`              → generate token baru
- `GET    /api/tokens/:id`          → detail token
- `PUT    /api/tokens/:id/status`   → aktif/nonaktifkan token
- `DELETE /api/tokens/:id`          → hapus token (hard delete)

Aturan: user **hanya bisa** kelola token **miliknya sendiri** (cek `user_id` di context JWT).

**6. Modifikasi `middleware/auth_middleware.go` → fungsi `APIKeyAuth()`**
- Cari token di tabel `api_tokens` (bukan lagi di `users.api_token`)
- Kondisi query:
  - `token = ?`
  - `is_active = true`
  - `expires_at IS NULL OR expires_at > NOW()`
- Update `last_used_at = NOW()`
- Simpan `user_id` ke context

**7. Modifikasi `main.go`**
- Update `AutoMigrate` → tambahkan `&models.APIToken{}`

**8. Update `routes/routes.go`**
- Tambahkan group `/api/tokens` dengan middleware JWT

---

### 🆙 UPGRADE 2 — TABEL `ai_models` (CRUD KHUSUS ADMIN)

**Tujuan**:
- Model HuggingFace disimpan di database, bukan hardcode di ENV
- Hanya admin yang bisa tambah/edit/hapus model
- User biasa cuma bisa pakai model yang `is_active = true`

#### 🔧 Yang harus dilakukan:

**1. Modifikasi `models/user.go`**
- Tambahkan kolom `Role`:
  ```go
  Role string `gorm:"type:varchar(20);default:'user'" json:"role"`
  ```

**2. Buat file baru: `models/ai_model.go`**
```go
type AIModel struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"type:varchar(120);not null" json:"name"`
    Slug        string    `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
    Description string    `gorm:"type:text" json:"description"`
    IsActive    bool      `gorm:"default:true" json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**3. Buat file baru: `dto/model_dto.go`**
```go
type CreateModelRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=120"`
    Slug        string `json:"slug" validate:"required,min=2,max=200"`
    Description string `json:"description" validate:"max=1000"`
    IsActive    bool   `json:"is_active"`
}

type UpdateModelRequest struct {
    Name        string `json:"name" validate:"omitempty,min=2,max=120"`
    Slug        string `json:"slug" validate:"omitempty,min=2,max=200"`
    Description string `json:"description" validate:"omitempty,max=1000"`
    IsActive    *bool  `json:"is_active"`
}
```

**4. Buat file baru: `middleware/admin_middleware.go`**
```go
func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != "admin" {
            utils.Error(c, http.StatusForbidden, "Akses hanya untuk admin", nil)
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**5. Modifikasi `utils/jwt.go`**
- Tambahkan `Role` ke struct `Claims`
- Update `GenerateJWT(userID uint, username, role string)`

**6. Modifikasi `handlers/auth_handler.go`**
- Register: set `Role = "user"` (default)
- Login: masukkan `role` ke JWT claims

**7. Buat file baru: `handlers/model_handler.go`**
Endpoint admin (protected JWT + AdminOnly):
- `GET    /api/admin/models`
- `POST   /api/admin/models`
- `GET    /api/admin/models/:id`
- `PUT    /api/admin/models/:id`
- `DELETE /api/admin/models/:id`

**8. Buat file baru: `seeders/seeder.go`**
- Buat admin default (dari ENV: ADMIN_USERNAME, ADMIN_PASSWORD)
- Buat model default dari ENV `HF_MODEL` (slug = nilai HF_MODEL, name = "Default Model")

**9. Modifikasi `services/huggingface.go`**
- Tambahkan validasi: cek apakah `model` (slug) ada di tabel `ai_models` dan `is_active = true` sebelum call HF
- Jika model tidak ada di DB → return error 400

**10. Modifikasi `main.go`**
- Tambahkan `AutoMigrate(&models.AIModel{})`
- Panggil `seeders.SeedAll()` setelah migrate

**11. Update `routes/routes.go`**
- Tambahkan group `/api/admin` dengan `JWTAuth() + AdminOnly()`

---

### 🆙 UPGRADE 3 — TABEL `configs` (KONFIG HF DARI ADMIN)

**Tujuan**:
- HF API Key, HF Model, dan HF URL bisa diatur dari **DB** (via admin panel)
- Fallback ke ENV jika config tidak ada di DB
- Admin bisa ubah tanpa restart server

#### 🔧 Yang harus dilakukan:

**1. Buat file baru: `models/config.go`**
```go
type Config struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Key         string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
    Value       string    `gorm:"type:text;not null" json:"value"`
    Description string    `gorm:"type:varchar(255)" json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**2. Buat file baru: `dto/config_dto.go`**
```go
type CreateConfigRequest struct {
    Key         string `json:"key" validate:"required,min=2,max=100"`
    Value       string `json:"value" validate:"required"`
    Description string `json:"description" validate:"max=255"`
}

type UpdateConfigRequest struct {
    Value       string `json:"value" validate:"omitempty"`
    Description string `json:"description" validate:"omitempty,max=255"`
}
```

**3. Buat file baru: `services/config_service.go`**
```go
// GetConfig — ambil config dari DB, fallback ke ENV
func GetConfig(key string) string {
    var cfg models.Config
    err := config.DB.Where("key = ?", key).First(&cfg).Error
    if err == nil && cfg.Value != "" {
        return cfg.Value
    }
    // fallback ke ENV
    switch key {
    case "HF_API_KEY":
        return config.Env.HFAPIKey
    case "HF_MODEL":
        return config.Env.HFModel
    case "HF_BASE_URL":
        return config.Env.HFBaseURL
    }
    return ""
}

// GetHuggingFaceConfig — helper ambil 3 config sekaligus
func GetHuggingFaceConfig() (apiKey, model, baseURL string) {
    return GetConfig("HF_API_KEY"),
           GetConfig("HF_MODEL"),
           GetConfig("HF_BASE_URL")
}
```

**4. Buat file baru: `handlers/config_handler.go`**
Endpoint admin (protected JWT + AdminOnly):
- `GET    /api/admin/configs`
- `POST   /api/admin/configs`
- `GET    /api/admin/configs/:id`
- `PUT    /api/admin/configs/:id`
- `DELETE /api/admin/configs/:id`

**5. Modifikasi `services/huggingface.go`**
- **Ganti** pembacaan `config.Env.HFAPIKey` dsb → jadi pakai `config_service.GetHuggingFaceConfig()`
- Jadi tiap request AI, config dibaca dari DB (fallback ENV)

**6. Modifikasi `seeders/seeder.go`**
- Tambahkan seeding config default dari ENV:
  - `HF_API_KEY` = `config.Env.HFAPIKey`
  - `HF_MODEL`   = `config.Env.HFModel`
  - `HF_BASE_URL` = `config.Env.HFBaseURL`
- **Hanya insert kalau belum ada** (idempotent, bisa dijalankan berulang)

**7. Modifikasi `main.go`**
- Tambahkan `AutoMigrate(&models.Config{})`

**8. Update `routes/routes.go`**
- Tambahkan route config di dalam group `/api/admin`

---

## 📌 ATURAN PENTING (WAJIB DIIKUTI)

1. **JANGAN rewrite** file yang tidak perlu diubah — cukup tambah/modifikasi
2. **PERTAHANKAN** semua endpoint lama tetap berjalan (backward compatible)
3. **GUNAKAN** pola yang sudah ada di project (utils.Error, utils.Success, dll)
4. **KONSISTEN** dengan format response: `{ success, message, data }` / `{ success, message, error }`
5. **MIGRASI AMAN**: karena kolom `users.api_token` akan dihapus, buatlah migrasi yang aman:
   - Sebelum hapus, **copy data** `users.api_token` yang lama ke tabel `api_tokens` (name: "Migrated Token")
   - Beri komentar di kode migrasi
6. **SEEDER idempotent**: kalau data sudah ada, jangan insert ulang (gunakan `FirstOrCreate`)
7. **VALIDASI** tetap pakai validator/v10 di semua DTO baru
8. **HTTP status code** sesuai standar (200/201/400/401/403/404/409/500)
9. **JANGAN** munculkan password/token di log yang tidak perlu
10. **BERI KOMENTAR** di setiap file yang diubah: 
    ```
    // [UPGRADE v3] — <alasan perubahan>
    ```

---

## 📂 OUTPUT YANG DIHARAPKAN

Tolong hasilkan dalam bentuk berikut:

### A. DAFTAR PERUBAHAN
Tabel ringkasan:
| File | Status | Keterangan |
|------|--------|------------|
| models/user.go | MODIFIED | hapus kolom APIToken, tambah relasi + Role |
| models/api_token.go | NEW | ... |
| ... | ... | ... |

### B. FILE YANG BARU DIBUAT
Hasilkan **full code** untuk file baru:
1. `models/api_token.go`
2. `models/ai_model.go`
3. `models/config.go`
4. `dto/token_dto.go`
5. `dto/model_dto.go`
6. `dto/config_dto.go`
7. `middleware/admin_middleware.go`
8. `handlers/token_handler.go`
9. `handlers/model_handler.go`
10. `handlers/config_handler.go`
11. `services/config_service.go`
12. `seeders/seeder.go`

### C. FILE YANG DIMODIFIKASI
Hasilkan **full code versi BARU** (bukan diff/patch) untuk file:
1. `models/user.go`
2. `utils/jwt.go`
3. `middleware/auth_middleware.go`
4. `handlers/auth_handler.go`
5. `services/huggingface.go`
6. `routes/routes.go`
7. `main.go`

### D. FILE YANG TIDAK BERUBAH
Cukup sebutkan (tidak usah tampilkan code).

### E. MIGRASI DATA (PENTING)
Sertakan function migrasi di `main.go` atau file terpisah `migrations/migrate.go` yang:
- Baca semua `users.api_token` yang tidak kosong
- Insert ke tabel `api_tokens` sebagai token migrasi (jangan sampai hilang)
- Baru setelah itu kolom bisa di-drop (atau biarkan sebagai legacy)

### F. CONTOH CURL BARU
Sertakan curl untuk:
- `POST /api/tokens` (generate token baru)
- `GET /api/tokens` (list token)
- `DELETE /api/tokens/:id`
- `POST /api/admin/models`
- `PUT /api/admin/configs/:id`
- Update `POST /v1/messages` (pastikan model yang dipakai ada di tabel `ai_models`)

---

## ✅ CHECKLIST AKHIR

- [ ] Tidak ada project baru yang dibuat, hanya upgrade
- [ ] Endpoint lama tetap jalan (backward compatible)
- [ ] Kolom `users.api_token` di-migrate ke tabel `api_tokens` sebelum dihapus
- [ ] Login menghasilkan token BARU di tabel `api_tokens` (tidak menimpa)
- [ ] User bisa kelola token miliknya sendiri
- [ ] Role admin & user berjalan
- [ ] Middleware AdminOnly() terpasang
- [ ] CRUD ai_models (admin only)
- [ ] CRUD configs (admin only)
- [ ] `/v1/messages` baca config dari DB (fallback ENV)
- [ ] `/v1/messages` validasi model dari tabel ai_models
- [ ] Seeder idempotent (aman dijalankan berkali-kali)
- [ ] Komentar `[UPGRADE v3]` di bagian yang diubah
- [ ] Semua DTO divalidasi dengan validator/v10
- [ ] Format response konsisten
- [ ] Compilable tanpa error

---

## 🎬 MULAI DARI MANA?

**Langkah 1**: Tampilkan dulu **DAFTAR PERUBAHAN** (bagian A) supaya saya review.
**Langkah 2**: Tunggu saya bilang "lanjut".
**Langkah 3**: Baru hasilkan file-file baru (bagian B).
**Langkah 4**: Lanjut file modifikasi (bagian C).
**Langkah 5**: Terakhir contoh curl (bagian F).