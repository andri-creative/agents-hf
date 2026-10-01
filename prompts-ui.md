# TUGAS
Buatkan saya sebuah aplikasi FRONTEND modern menggunakan 
**Vite + Vue 3 + TypeScript + Tailwind CSS v4 + shadcn-vue** yang terhubung 
ke backend REST API Go (base URL: `http://localhost:8080`).

Aplikasi ini adalah **dashboard admin & user** untuk mengelola AI Generation 
service yang sudah berjalan di backend.

---

## 1. TEKNOLOGI FRONTEND (WAJIB)

| Komponen         | Tools / Library                                     |
|------------------|-----------------------------------------------------|
| Build Tool       | Vite                                                |
| Framework        | Vue 3 (Composition API + `<script setup>`)          |
| Bahasa           | TypeScript                                          |
| Styling          | **Tailwind CSS v4**                                 |
| Vite Plugin      | **@tailwindcss/vite** (BUKAN PostCSS)               |
| UI Components    | shadcn-vue (https://www.shadcn-vue.com)             |
| Icons            | lucide-vue-next                                     |
| Routing          | vue-router 4                                        |
| State Management | Pinia                                               |
| HTTP Client      | axios                                               |
| Form Validation  | vee-validate + zod                                  |
| Notifikasi       | vue-sonner (toast)                                  |
| Dark Mode        | Custom variant di Tailwind v4                       |

**PENTING — Tailwind v4**:
- ❌ JANGAN buat `tailwind.config.js`
- ❌ JANGAN buat `postcss.config.js`
- ❌ JANGAN install `autoprefixer` (sudah built-in)
- ❌ JANGAN pakai `@tailwind base/components/utilities`
- ✅ GUNAKAN `@tailwindcss/vite` plugin
- ✅ GUNAKAN `@import "tailwindcss"` di CSS
- ✅ GUNAKAN `@theme {}` untuk custom token (warna, breakpoint, dll)
- ✅ GUNAKAN `@variant dark` untuk dark mode

---

## 2. STRUKTUR FOLDER (WAJIB DIIKUTI)

```
frontend/
├── .env
├── .env.example
├── index.html
├── package.json
├── vite.config.ts
├── tsconfig.json
├── tsconfig.node.json
├── components.json              ← konfigurasi shadcn-vue
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── style.css                ← @import "tailwindcss" + @theme
│   ├── env.d.ts
│   ├── lib/
│   │   └── utils.ts             ← helper cn() dari shadcn
│   ├── assets/
│   ├── router/
│   │   └── index.ts             ← vue-router + route guard
│   ├── stores/
│   │   ├── auth.ts              ← Pinia store auth
│   │   └── theme.ts             ← dark/light mode
│   ├── services/
│   │   ├── api.ts               ← axios instance + interceptor
│   │   ├── auth.service.ts
│   │   ├── token.service.ts
│   │   ├── model.service.ts
│   │   ├── config.service.ts
│   │   └── ai.service.ts
│   ├── types/
│   │   ├── auth.ts
│   │   ├── token.ts
│   │   ├── model.ts
│   │   ├── config.ts
│   │   └── ai.ts
│   ├── composables/
│   │   ├── useAuth.ts
│   │   ├── useToast.ts
│   │   └── useResponsive.ts     ← helper deteksi device
│   ├── components/
│   │   ├── ui/                  ← shadcn-vue components (auto-generated)
│   │   ├── layout/
│   │   │   ├── AppLayout.vue
│   │   │   ├── Sidebar.vue
│   │   │   ├── Topbar.vue
│   │   │   ├── MobileNav.vue
│   │   │   └── ThemeToggle.vue
│   │   ├── shared/
│   │   │   ├── DataTable.vue
│   │   │   ├── ConfirmDialog.vue
│   │   │   ├── LoadingSpinner.vue
│   │   │   └── EmptyState.vue
│   │   └── auth/
│   │       └── AuthGuard.vue
│   └── views/
│       ├── auth/
│       │   ├── LoginView.vue
│       │   └── RegisterView.vue
│       ├── dashboard/
│       │   └── DashboardView.vue
│       ├── playground/
│       │   └── PlaygroundView.vue
│       ├── tokens/
│       │   └── TokensView.vue
│       ├── admin/
│       │   ├── ModelsView.vue
│       │   ├── ConfigsView.vue
│       │   └── UsersView.vue
│       ├── profile/
│       │   └── ProfileView.vue
│       └── NotFoundView.vue
```

---

## 3. ENV FRONTEND (.env.example)

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_NAME=AI Generate Dashboard
```

Semua request ke backend harus melalui `VITE_API_BASE_URL`.

---

## 4. KONFIGURASI TAILWIND V4 (WAJIB — JANGAN PAKAI v3 STYLE)

### 📄 `vite.config.ts`
```ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),   // ← WAJIB, bukan postcss
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
  },
})
```

### 📄 `src/style.css` (INI PENGGANTI tailwind.config.js)
```css
@import "tailwindcss";

/* Dark mode manual pakai class */
@variant dark (&:where(.dark, .dark *));

/* Custom design tokens */
@theme {
  /* Warna primary (sesuaikan dengan shadcn theme) */
  --color-primary: oklch(0.55 0.2 250);
  --color-primary-foreground: oklch(0.98 0 0);
  
  /* Custom breakpoint untuk TV */
  --breakpoint-tv: 1920px;
  --breakpoint-3xl: 2560px;
  
  /* Font */
  --font-sans: 'Inter', ui-sans-serif, system-ui, sans-serif;
}

/* Base styles */
@layer base {
  html {
    scroll-behavior: smooth;
  }
  body {
    @apply bg-background text-foreground antialiased;
  }
}
```

### ❌ JANGAN BUAT FILE INI:
- `tailwind.config.js`
- `postcss.config.js`
- `.postcssrc`

---

## 5. RESPONSIVE BREAKPOINT (WAJIB)

Gunakan **desktop-first design**, dengan breakpoint Tailwind v4:

| Device       | Breakpoint    | Layout                                              |
|--------------|---------------|-----------------------------------------------------|
| **TV / 4K**  | `tv:` (1920px) | Sidebar lebar 280px, konten max 1600px, font besar  |
| **Ultra**    | `3xl:` (2560px) | Konten max 1800px                                  |
| **Desktop**  | `xl:` (1280px) | Sidebar fixed kiri 260px, konten 2 kolom            |
| **Laptop**   | `lg:` (1024px) | Sidebar fixed kiri 260px, konten 1 kolom            |
| **Tablet**   | `md:` (768px)  | Sidebar collapsible icon-only 64px                  |
| **HP**       | default (<768px)| Sidebar → drawer / bottom nav, konten full         |

**Aturan**:
- Semua breakpoint custom (`tv:`, `3xl:`) sudah didefinisikan di `@theme` — JANGAN pakai config screens
- Font scaling di TV: `tv:text-lg`
- Touch target minimum 44px di mobile
- Gunakan utility Tailwind v4 standar: `sm: md: lg: xl: 2xl: tv: 3xl:`

---

## 6. FITUR UI YANG DIMINTA

### 🔹 FITUR 1 — AUTHENTIKASI

**Halaman Register** (`/register`)
- Form: full_name, username, password
- Validasi real-time dengan vee-validate + zod
- Loading state saat submit
- Redirect ke /login setelah sukses
- Toast sukses/gagal

**Halaman Login** (`/login`)
- Form: username, password
- Simpan JWT + API token di Pinia + localStorage
- Redirect ke /dashboard
- "Remember me" (opsional)

**Logout**
- Hapus token dari store & localStorage
- Redirect ke /login

**Route Guard**:
- Route dengan `meta.requiresAuth` → cek JWT di store
- Route dengan `meta.requiresAdmin` → cek `role === 'admin'`
- Redirect otomatis jika tidak berhak

---

### 🔹 FITUR 2 — DASHBOARD (`/dashboard`)

- **Stat cards** (grid responsive):
  - Total Token
  - Token Aktif
  - Total Model (kalau admin)
  - Request Hari Ini (placeholder)
- **Chart** sederhana (opsional pakai chart.js / vue-chartjs)
- **Quick actions**: Generate AI, Buat Token, Tambah Model (admin)
- **Recent activity** list (opsional)

---

### 🔹 FITUR 3 — AI PLAYGROUND (`/playground`)

UI utama untuk test `/v1/messages`:
- **Layout 2 kolom** (desktop) / stacked (mobile):
  - **Kiri**: Form input
    - Dropdown pilih model (fetch dari `/api/admin/models`)
    - Textarea prompt (auto resize)
    - Parameter: temperature, max_tokens (slider)
    - Tombol "Generate"
  - **Kanan**: Hasil response
    - Markdown render untuk hasil AI
    - Copy button
    - Tombol "Regenerate"
- **Streaming response** (opsional, kalau backend support SSE)
- Simpan history di localStorage

---

### 🔹 FITUR 4 — MANAJEMEN TOKEN (`/tokens`)

- **Data table** dengan kolom: Token (masked), Name, Status, Last Used, Created
- **Actions**: Copy, Revoke, Delete
- **Button** "Buat Token Baru" → buka dialog:
  - Input name
  - Input expires_at (opsional, date picker)
- Setelah dibuat → tampilkan token **sekali saja** dengan tombol copy
- **Filter & search** by name / status

---

### 🔹 FITUR 5 — MANAJEMEN MODEL (ADMIN ONLY) (`/admin/models`)

- **Data table**: Name, Slug, Description, Status (aktif/nonaktif), Actions
- **Create / Edit dialog** dengan shadcn `Dialog`
- **Toggle active** switch inline
- **Delete** dengan `ConfirmDialog`
- **Search & pagination**

---

### 🔹 FITUR 6 — MANAJEMEN CONFIG (ADMIN ONLY) (`/admin/configs`)

- **Data table**: Key, Value (masked untuk API key), Description, Actions
- **Edit dialog** inline
- **Highlight** key sensitif (`HF_API_KEY`) → value disembunyikan (`•••••`) + tombol show/hide
- **Confirm dialog** untuk delete

---

### 🔹 FITUR 7 — PROFILE (`/profile`)

- Lihat & edit `full_name`
- Ganti password (form terpisah)
- Info role
- Lihat daftar token aktif

---

### 🔹 FITUR 8 — 404 PAGE

Halaman NotFound yang rapi + tombol "Kembali ke Dashboard".

---

## 7. LAYOUT UTAMA (AppLayout)

### Struktur Layout:
```
┌─────────────────────────────────────────────────────┐
│  Sidebar  │  Topbar (search, notif, theme, user)     │
│           ├─────────────────────────────────────────┤
│  (nav)    │                                          │
│           │           <router-view />                │
│           │                                          │
└───────────┴──────────────────────────────────────────┘
```

### Sidebar:
- **TV/Desktop**: fixed, lebar 260px, show label + icon
- **Tablet**: collapsible, lebar 64px, icon only + tooltip
- **Mobile**: hidden, muncul sebagai drawer (shadcn `Sheet`)

### Topbar:
- Hamburger menu (mobile/tablet)
- Search bar (cmd+k optional)
- Notification bell
- Theme toggle (dark/light)
- User avatar dropdown (profile, logout)

### Menu Sidebar (role-based):
```
Dashboard         (semua)
Playground        (semua)
API Tokens        (semua)
────────────────
ADMIN
Models            (admin)
Configs           (admin)
Users             (admin, opsional)
────────────────
Profile           (semua)
Logout            (semua)
```

---

## 8. INTEGRASI API (Axios)

### `services/api.ts` — Setup Axios
```ts
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'
import { toast } from 'vue-sonner'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 60000,
})

// Request interceptor → attach JWT
api.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

// Response interceptor → handle 401 & error toast
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      const auth = useAuthStore()
      auth.logout()
      router.push('/login')
    }
    toast.error(err.response?.data?.message || 'Terjadi kesalahan')
    return Promise.reject(err)
  }
)

export default api
```

### `services/auth.service.ts`
```ts
import api from './api'
import type { LoginRequest, RegisterRequest, AuthResponse } from '@/types/auth'

export const authService = {
  register: (data: RegisterRequest) => api.post<AuthResponse>('/api/register', data),
  login: (data: LoginRequest) => api.post<AuthResponse>('/api/login', data),
  me: () => api.get('/api/me'),
}
```

### `services/token.service.ts`
```ts
import api from './api'
export const tokenService = {
  list: () => api.get('/api/tokens'),
  create: (data: { name: string; expires_at?: string }) => api.post('/api/tokens', data),
  detail: (id: number) => api.get(`/api/tokens/${id}`),
  toggleStatus: (id: number, isActive: boolean) =>
    api.put(`/api/tokens/${id}/status`, { is_active: isActive }),
  remove: (id: number) => api.delete(`/api/tokens/${id}`),
}
```

### `services/model.service.ts` (Admin)
```ts
import api from './api'
export const modelService = {
  list: () => api.get('/api/admin/models'),
  create: (data: any) => api.post('/api/admin/models', data),
  update: (id: number, data: any) => api.put(`/api/admin/models/${id}`, data),
  remove: (id: number) => api.delete(`/api/admin/models/${id}`),
}
```

### `services/config.service.ts` (Admin)
```ts
import api from './api'
export const configService = {
  list: () => api.get('/api/admin/configs'),
  create: (data: any) => api.post('/api/admin/configs', data),
  update: (id: number, data: any) => api.put(`/api/admin/configs/${id}`, data),
  remove: (id: number) => api.delete(`/api/admin/configs/${id}`),
}
```

### `services/ai.service.ts`
```ts
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

export const aiService = {
  generate: async (data: { model: string; prompt: string }) => {
    const auth = useAuthStore()
    return axios.post(
      `${import.meta.env.VITE_API_BASE_URL}/v1/messages`,
      data,
      { headers: { 'X-API-Key': auth.apiToken } }
    )
  },
}
```

---

## 9. STATE MANAGEMENT (Pinia)

### `stores/auth.ts`
```ts
import { defineStore } from 'pinia'

interface User {
  id: number
  full_name: string
  username: string
  role: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('jwt') || '',
    apiToken: localStorage.getItem('api_token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null') as User | null,
  }),
  getters: {
    isLoggedIn: (s) => !!s.token,
    isAdmin: (s) => s.user?.role === 'admin',
  },
  actions: {
    setAuth({ token, apiToken, user }: any) {
      this.token = token
      this.apiToken = apiToken
      this.user = user
      localStorage.setItem('jwt', token)
      localStorage.setItem('api_token', apiToken)
      localStorage.setItem('user', JSON.stringify(user))
    },
    logout() {
      this.token = ''
      this.apiToken = ''
      this.user = null
      localStorage.clear()
    },
  },
})
```

### `stores/theme.ts`
```ts
import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', {
  state: () => ({
    isDark: localStorage.getItem('theme') === 'dark',
  }),
  actions: {
    toggle() {
      this.isDark = !this.isDark
      document.documentElement.classList.toggle('dark', this.isDark)
      localStorage.setItem('theme', this.isDark ? 'dark' : 'light')
    },
    init() {
      document.documentElement.classList.toggle('dark', this.isDark)
    },
  },
})
```

---

## 10. SHADCN-VUE COMPONENTS (WAJIB PAKAI)

Jalankan init shadcn-vue:
```bash
npx shadcn-vue@latest init
```

Install komponen yang dibutuhkan:
```bash
npx shadcn-vue@latest add button
npx shadcn-vue@latest add input
npx shadcn-vue@latest add label
npx shadcn-vue@latest add card
npx shadcn-vue@latest add dialog
npx shadcn-vue@latest add alert-dialog
npx shadcn-vue@latest add dropdown-menu
npx shadcn-vue@latest add table
npx shadcn-vue@latest add form
npx shadcn-vue@latest add select
npx shadcn-vue@latest add switch
npx shadcn-vue@latest add badge
npx shadcn-vue@latest add avatar
npx shadcn-vue@latest add sheet
npx shadcn-vue@latest add sonner
npx shadcn-vue@latest add skeleton
npx shadcn-vue@latest add tooltip
npx shadcn-vue@latest add tabs
npx shadcn-vue@latest add textarea
npx shadcn-vue@latest add separator
npx shadcn-vue@latest add popover
npx shadcn-vue@latest add calendar
```

---

## 11. UX / UI RULES (WAJIB)

### 🎨 Visual
- **Tema**: gunakan shadcn default (neutral) — support dark & light mode
- **Warna aksen**: biru/indigo untuk primary
- **Spacing**: konsisten pakai skala Tailwind (gap-4, p-6, dll)
- **Font**: Inter (default shadcn)
- **Border radius**: `rounded-lg` (konsisten)
- **Shadow**: subtle saja — perhatikan **Tailwind v4 rename**: `shadow-sm` → `shadow-xs`, `rounded-sm` → `rounded-xs`

### ♿ Accessibility
- Semua tombol punya `aria-label` jika icon-only
- Keyboard navigable (Tab, Enter, Esc)
- Focus ring selalu terlihat
- Color contrast minimum AA

### 🚦 Loading & Error State
- **Setiap** fetch data harus ada: loading skeleton / spinner
- **Setiap** error → toast + fallback UI
- **Setiap** form submit → disable button + loading indicator
- Gunakan `Skeleton` dari shadcn untuk table/list loading

### 📱 Responsive
- Tabel di mobile → convert ke card list (stacked)
- Dialog di mobile → full screen sheet
- Form di mobile → full width, label di atas
- Test di ukuran: 360px, 768px, 1280px, 1920px, 2560px

### 🎯 Micro-interactions
- Hover effects di tombol & card
- Transition 150-200ms
- Toast untuk setiap aksi CRUD
- Confirm dialog untuk aksi destruktif (delete, revoke)

---

## 12. OUTPUT YANG DIHARAPKAN

Tolong hasilkan **SEMUA FILE LENGKAP** dengan urutan:

### A. Konfigurasi Project (Tailwind v4 style)
1. `package.json` (dependencies v4: `tailwindcss`, `@tailwindcss/vite`)
2. `vite.config.ts` (dengan `@tailwindcss/vite` plugin)
3. `tsconfig.json`
4. `tsconfig.node.json`
5. `components.json`
6. `.env.example`
7. `index.html`
8. `src/env.d.ts`

### B. Core Files
9. `src/style.css` (`@import "tailwindcss"` + `@theme` + `@variant dark`)
10. `src/main.ts`
11. `src/App.vue`
12. `src/lib/utils.ts`

### C. Types & Services
13. `src/types/auth.ts`, `token.ts`, `model.ts`, `config.ts`, `ai.ts`
14. `src/services/api.ts`
15. `src/services/auth.service.ts`
16. `src/services/token.service.ts`
17. `src/services/model.service.ts`
18. `src/services/config.service.ts`
19. `src/services/ai.service.ts`

### D. Store & Router
20. `src/stores/auth.ts`
21. `src/stores/theme.ts`
22. `src/router/index.ts`

### E. Layout & Shared Components
23. `src/components/layout/AppLayout.vue`
24. `src/components/layout/Sidebar.vue`
25. `src/components/layout/Topbar.vue`
26. `src/components/layout/MobileNav.vue`
27. `src/components/layout/ThemeToggle.vue`
28. `src/components/shared/DataTable.vue`
29. `src/components/shared/ConfirmDialog.vue`
30. `src/components/shared/LoadingSpinner.vue`
31. `src/components/shared/EmptyState.vue`

### F. Views
32. `src/views/auth/LoginView.vue`
33. `src/views/auth/RegisterView.vue`
34. `src/views/dashboard/DashboardView.vue`
35. `src/views/playground/PlaygroundView.vue`
36. `src/views/tokens/TokensView.vue`
37. `src/views/admin/ModelsView.vue`
38. `src/views/admin/ConfigsView.vue`
39. `src/views/profile/ProfileView.vue`
40. `src/views/NotFoundView.vue`

### G. Dokumentasi
41. `README.md` — cara install & run (Tailwind v4)
42. Sertakan **perintah lengkap** untuk setup

---

## 13. CONTOH CARA PAKAI (SERTAKAN DI README)

```bash
# 1. Buat project Vite
npm create vite@latest frontend -- --template vue-ts
cd frontend

# 2. Install Tailwind v4 (JANGAN install autoprefixer/postcss)
npm install tailwindcss @tailwindcss/vite

# 3. Tambahkan @tailwindcss/vite ke vite.config.ts (lihat section 4)

# 4. Ganti src/style.css dengan @import "tailwindcss"

# 5. Init shadcn-vue
npx shadcn-vue@latest init

# 6. Install komponen shadcn (lihat section 10)

# 7. Install dependencies lain
npm install vue-router@4 pinia axios vee-validate @vee-validate/zod zod
npm install lucide-vue-next vue-sonner

# 8. Setup .env
cp .env.example .env

# 9. Jalankan
npm run dev
```

---

## 14. CHECKLIST AKHIR

- [ ] Semua 42 file dihasilkan lengkap
- [ ] Setup Vite + Vue 3 + TS + **Tailwind v4** + shadcn-vue benar
- [ ] **TIDAK ADA** `tailwind.config.js` / `postcss.config.js`
- [ ] `@tailwindcss/vite` plugin terpasang di `vite.config.ts`
- [ ] `style.css` pakai `@import "tailwindcss"` + `@theme` + `@variant dark`
- [ ] Breakpoint `tv:` dan `3xl:` didefinisikan di `@theme`
- [ ] Dark mode berjalan
- [ ] Router guard (auth & admin) berjalan
- [ ] Axios interceptor attach JWT + handle 401
- [ ] Login/Register berjalan ke backend
- [ ] Playground bisa call `/v1/messages` dengan X-API-Key
- [ ] CRUD Tokens (user)
- [ ] CRUD Models (admin)
- [ ] CRUD Configs (admin)
- [ ] Layout responsive: TV, Desktop, Tablet, Mobile
- [ ] Loading skeleton di semua fetch
- [ ] Toast untuk semua aksi CRUD
- [ ] Confirm dialog untuk delete
- [ ] Semua form pakai vee-validate + zod
- [ ] Dark/light toggle berjalan
- [ ] 404 page rapi
- [ ] README lengkap dengan langkah setup Tailwind v4

---

## 15. ATURAN PENTING

1. **JANGAN** buat backend baru — backend sudah ada di `http://localhost:8080`
2. **JANGAN** buat `tailwind.config.js` — Tailwind v4 pakai `@theme` di CSS
3. **JANGAN** pakai `@tailwind base/components/utilities` — pakai `@import "tailwindcss"`
4. **JANGAN** install `autoprefixer` — sudah built-in di v4
5. **GUNAKAN** shadcn-vue untuk semua komponen
6. **PERHATIKAN** rename utility v4: `shadow-sm`→`shadow-xs`, `rounded-sm`→`rounded-xs`, `blur-sm`→`blur-xs`
7. **CSS variable syntax**: `bg-(--var)` bukan `bg-[--var]`
8. **TYPE-SAFE**: semua response API harus punya TypeScript type
9. **KOMPONEN REUSABLE**: DataTable, ConfirmDialog, dll
10. **KONSISTEN**: gunakan `<script setup lang="ts">` di semua `.vue`
11. **JANGAN** pakai `any` kecuali terpaksa
12. **RESPONSIVE**: uji minimal di 360px, 768px, 1280px, 1920px, 2560px
13. **ACCESSIBLE**: label form, aria-attribute, keyboard nav
14. **TOAST** untuk setiap aksi yang memodifikasi data
```

---

## 🎯 PERBANDINGAN LENGKAP: v3 vs v4 (yang ada di prompt ini)

| Aspek | v3 (lama) | v4 (prompt ini) |
|---|---|---|
| **File config** | `tailwind.config.js` | `@theme {}` di `style.css` |
| **Vite plugin** | `postcss.config.js` + `tailwindcss` | `@tailwindcss/vite` |
| **CSS import** | `@tailwind base; @tailwind components; @tailwind utilities;` | `@import "tailwindcss";` |
| **Content detection** | Manual `content: [...]` | Otomatis |
| **Autoprefixer** | Install terpisah | Built-in |
| **Custom breakpoint** | `screens: { tv: '1920px' }` di config | `--breakpoint-tv: 1920px` di `@theme` |
| **Dark mode config** | `darkMode: 'class'` di config | `@variant dark (&:where(.dark, .dark *));` |
| **Utility rename** | `shadow-sm`, `rounded-sm` | `shadow-xs`, `rounded-xs` |
| **CSS var syntax** | `bg-[--var]` | `bg-(--var)` |

---

Prompt ini **sudah LENGKAP** (42 file, 15 section, checklist, aturan, contoh setup) — sama lengkapnya dengan prompt sebelumnya, hanya disesuaikan ke **Tailwind v4**. Tinggal copy-paste. 🚀

Mau saya buatkan **prompt terpisah untuk 1 halaman spesifik** (misal Playground dengan streaming) atau **prompt deployment** (Vercel/Netlify/Docker)?