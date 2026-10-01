# AI Generate Dashboard - Frontend

Dashboard modern untuk mengelola AI Generation service menggunakan **Vite + Vue 3 + TypeScript + Tailwind CSS v4**.

## 🚀 Tech Stack

- **Vite** - Build tool
- **Vue 3** - Framework (Composition API + `<script setup>`)
- **TypeScript** - Type safety
- **Tailwind CSS v4** - Styling (dengan `@tailwindcss/vite` plugin)
- **Vue Router 4** - Routing
- **Pinia** - State management
- **Axios** - HTTP client
- **Lucide Vue Next** - Icons

## 📁 Struktur Project

```
ui/
├── src/
│   ├── components/
│   │   └── layout/          # AppLayout, Sidebar, dll
│   ├── views/
│   │   ├── auth/            # Login, Register
│   │   ├── dashboard/       # Dashboard utama
│   │   ├── playground/      # AI Playground
│   │   ├── tokens/          # Manajemen API tokens
│   │   ├── admin/           # Admin pages (Models, Configs)
│   │   └── profile/         # User profile
│   ├── services/            # API services
│   ├── stores/              # Pinia stores
│   ├── types/               # TypeScript types
│   ├── router/              # Vue Router config
│   ├── style.css            # Tailwind CSS v4 config
│   └── main.ts
├── .env                     # Environment variables
├── vite.config.ts
└── package.json
```

## ⚙️ Setup & Installation

### 1. Install Dependencies

```bash
cd ui
npm install
```

### 2. Setup Environment Variables

Buat file `.env` (sudah ada dari `.env.example`):

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_NAME=AI Generate Dashboard
```

### 3. Jalankan Development Server

```bash
npm run dev
```

Aplikasi akan berjalan di `http://localhost:5173`

### 4. Build untuk Production

```bash
npm run build
```

Output ada di folder `dist/`.

### 5. Preview Production Build

```bash
npm run preview
```

## 🎯 Fitur

### 1. **Authentication**
- Login & Register
- JWT token management
- Auto redirect jika sudah login

### 2. **Dashboard**
- Statistik: Total tokens, token aktif, total models
- Quick actions
- Welcome message

### 3. **AI Playground**
- Test AI generation dengan prompt
- Pilih model yang tersedia
- Copy hasil ke clipboard
- Error handling

### 4. **API Tokens Management**
- List semua token milik user
- Buat token baru (dengan expires_at opsional)
- Toggle status aktif/nonaktif
- Delete token
- Copy token ke clipboard
- Token masking untuk keamanan

### 5. **Admin - Models Management** (Admin only)
- CRUD AI models
- Toggle active/inactive
- Search & filter

### 6. **Admin - Configs Management** (Admin only)
- CRUD configs (HF_API_KEY, HF_MODEL, HF_BASE_URL)
- Sensitive value masking (untuk KEY/SECRET)
- Show/hide toggle untuk value sensitif

### 7. **Profile**
- Lihat info user
- Quick links ke fitur lain

### 8. **Dark Mode**
- Toggle light/dark theme
- Tersimpan di localStorage

## 🎨 Tailwind CSS v4

Project ini menggunakan **Tailwind CSS v4** dengan konfigurasi modern:

- ✅ **TIDAK ADA** `tailwind.config.js`
- ✅ **TIDAK ADA** `postcss.config.js`
- ✅ Menggunakan `@tailwindcss/vite` plugin
- ✅ Config di `src/style.css` dengan `@theme` dan `@variant`
- ✅ Dark mode dengan `@variant dark`
- ✅ Custom breakpoints: `tv:` (1920px), `3xl:` (2560px)

### Contoh `src/style.css`:

```css
@import "tailwindcss";

@variant dark (&:where(.dark, .dark *));

@theme {
  --color-primary: oklch(0.55 0.2 250);
  --color-background: oklch(1 0 0);
  --color-foreground: oklch(0.15 0 0);
  --breakpoint-tv: 1920px;
  --font-sans: 'Inter', ui-sans-serif, system-ui, sans-serif;
}
```

## 🔐 Authentication Flow

1. User login → dapat `jwt_token` + `api_token`
2. Token disimpan di Pinia store + localStorage
3. JWT dipakai untuk endpoint `/api/*` (Authorization: Bearer)
4. API token dipakai untuk `/v1/messages` (X-API-Key header)
5. Axios interceptor otomatis attach JWT ke setiap request
6. Jika 401 → logout otomatis & redirect ke `/login`

## 🛡️ Route Guard

```typescript
// meta.requiresAuth → butuh login
// meta.requiresAdmin → butuh role admin
```

- `/login`, `/register` → public
- `/dashboard`, `/playground`, `/tokens`, `/profile` → butuh login
- `/admin/models`, `/admin/configs` → butuh admin

## 📱 Responsive Design

- **Mobile** (<768px): Sidebar jadi drawer, stacked layout
- **Tablet** (768px-1023px): Sidebar collapsible
- **Desktop** (1024px+): Sidebar fixed, 2-kolom layout
- **TV/4K** (1920px+): Font lebih besar, konten max-width

## 🔧 Development Tips

### Hot Module Replacement (HMR)

Vite mendukung HMR, jadi perubahan kode langsung ter-refresh tanpa reload penuh.

### TypeScript Error

Jika ada error TypeScript tentang `baseUrl` deprecated:

```bash
# Ignore warning dan langsung build
npx vite build
```

### API Backend

Pastikan backend Go berjalan di `http://localhost:8080` sebelum menjalankan frontend.

### Debugging

- Vue DevTools: Install extension untuk debug Pinia & Router
- Network tab: Cek request/response API
- Console: Error handling sudah lengkap

## 📦 Dependencies Utama

```json
{
  "dependencies": {
    "vue": "^3.5.42",
    "vue-router": "^4.2.5",
    "pinia": "^2.1.7",
    "axios": "^1.6.0",
    "tailwindcss": "^4.3.3",
    "@tailwindcss/vite": "^4.3.3",
    "lucide-vue-next": "^0.292.0"
  }
}
```

## 🚧 Known Issues

1. **TypeScript 6 Warning**: `baseUrl` deprecated di TS 6.0, tapi masih berfungsi. Fix: ignore atau tunggu update `@vue/tsconfig`.
2. **Vite Warning**: `__dirname` deprecated, gunakan `import.meta.dirname` di masa depan.

## 📝 Todo / Future Improvements

- [ ] Add form validation dengan vee-validate + zod
- [ ] Add toast notification dengan vue-sonner
- [ ] Add streaming response untuk AI playground
- [ ] Add pagination untuk tabel
- [ ] Add search & filter di semua tabel
- [ ] Add charts untuk statistik
- [ ] Add unit tests

## 🤝 Contributing

1. Fork repository
2. Buat branch baru (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## 📄 License

MIT License - bebas dipakai untuk project pribadi maupun komersial.

---

**Dibuat dengan ❤️ menggunakan Vite + Vue 3 + Tailwind CSS v4**
