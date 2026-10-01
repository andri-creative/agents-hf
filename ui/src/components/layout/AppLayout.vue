<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { Menu, X, Sun, Moon, LogOut, User } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()

const sidebarOpen = ref(false)

const menuItems = computed(() => {
  const items = [
    { name: 'Dashboard', path: '/dashboard', icon: 'LayoutDashboard' },
    { name: 'Playground', path: '/playground', icon: 'Sparkles' },
    { name: 'API Tokens', path: '/tokens', icon: 'Key' },
  ]

  if (authStore.isAdmin) {
    items.push(
      { name: 'Models', path: '/admin/models', icon: 'Bot' },
      { name: 'Configs', path: '/admin/configs', icon: 'Settings' }
    )
  }

  items.push({ name: 'Profile', path: '/profile', icon: 'User' })

  return items
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Sidebar Desktop -->
    <aside class="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700">
      <div class="flex flex-col flex-1 min-h-0">
        <div class="flex items-center h-16 px-6 border-b border-gray-200 dark:border-gray-700">
          <h1 class="text-xl font-bold">AI Dashboard</h1>
        </div>
        
        <nav class="flex-1 px-3 py-4 space-y-1">
          <router-link
            v-for="item in menuItems"
            :key="item.path"
            :to="item.path"
            class="flex items-center px-3 py-2 text-sm font-medium rounded-md hover:bg-gray-100 dark:hover:bg-gray-700"
            active-class="bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400"
          >
            {{ item.name }}
          </router-link>
        </nav>

        <div class="p-4 border-t border-gray-200 dark:border-gray-700">
          <div class="flex items-center justify-between mb-3">
            <span class="text-sm font-medium">{{ authStore.user?.username }}</span>
            <button @click="themeStore.toggle()" class="p-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700">
              <Sun v-if="themeStore.isDark" class="w-5 h-5" />
              <Moon v-else class="w-5 h-5" />
            </button>
          </div>
          <button
            @click="handleLogout"
            class="flex items-center w-full px-3 py-2 text-sm font-medium text-red-600 dark:text-red-400 rounded-md hover:bg-red-50 dark:hover:bg-red-900/20"
          >
            <LogOut class="w-4 h-4 mr-2" />
            Logout
          </button>
        </div>
      </div>
    </aside>

    <!-- Mobile Header -->
    <div class="lg:hidden flex items-center justify-between h-16 px-4 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <button @click="sidebarOpen = !sidebarOpen" class="p-2">
        <Menu class="w-6 h-6" />
      </button>
      <h1 class="text-lg font-bold">AI Dashboard</h1>
      <button @click="themeStore.toggle()" class="p-2">
        <Sun v-if="themeStore.isDark" class="w-5 h-5" />
        <Moon v-else class="w-5 h-5" />
      </button>
    </div>

    <!-- Mobile Sidebar -->
    <div v-if="sidebarOpen" class="lg:hidden fixed inset-0 z-50 bg-black/50" @click="sidebarOpen = false">
      <aside class="fixed inset-y-0 left-0 w-64 bg-white dark:bg-gray-800" @click.stop>
        <div class="flex items-center justify-between h-16 px-6 border-b border-gray-200 dark:border-gray-700">
          <h1 class="text-xl font-bold">AI Dashboard</h1>
          <button @click="sidebarOpen = false">
            <X class="w-6 h-6" />
          </button>
        </div>
        
        <nav class="px-3 py-4 space-y-1">
          <router-link
            v-for="item in menuItems"
            :key="item.path"
            :to="item.path"
            @click="sidebarOpen = false"
            class="flex items-center px-3 py-2 text-sm font-medium rounded-md hover:bg-gray-100 dark:hover:bg-gray-700"
            active-class="bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400"
          >
            {{ item.name }}
          </router-link>
        </nav>

        <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200 dark:border-gray-700">
          <div class="text-sm font-medium mb-3">{{ authStore.user?.username }}</div>
          <button
            @click="handleLogout"
            class="flex items-center w-full px-3 py-2 text-sm font-medium text-red-600 dark:text-red-400 rounded-md hover:bg-red-50 dark:hover:bg-red-900/20"
          >
            <LogOut class="w-4 h-4 mr-2" />
            Logout
          </button>
        </div>
      </aside>
    </div>

    <!-- Main Content -->
    <main class="lg:pl-64">
      <div class="py-6 px-4 sm:px-6 lg:px-8">
        <router-view />
      </div>
    </main>
  </div>
</template>
