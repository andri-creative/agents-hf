// src/stores/auth.ts
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
    setAuth({ token, apiToken, user }: { token: string; apiToken: string; user: User }) {
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
