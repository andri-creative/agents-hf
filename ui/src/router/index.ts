// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/RegisterView.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/components/layout/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/dashboard',
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/DashboardView.vue'),
        },
        {
          path: 'playground',
          name: 'playground',
          component: () => import('@/views/playground/PlaygroundView.vue'),
        },
        {
          path: 'tokens',
          name: 'tokens',
          component: () => import('@/views/tokens/TokensView.vue'),
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('@/views/profile/ProfileView.vue'),
        },
        {
          path: 'admin/models',
          name: 'admin-models',
          component: () => import('@/views/admin/ModelsView.vue'),
          meta: { requiresAdmin: true },
        },
        {
          path: 'admin/configs',
          name: 'admin-configs',
          component: () => import('@/views/admin/ConfigsView.vue'),
          meta: { requiresAdmin: true },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
    },
  ],
})

// Route guard
router.beforeEach((to, from, next) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    next('/login')
  } else if (to.meta.requiresAdmin && !auth.isAdmin) {
    next('/dashboard')
  } else if ((to.name === 'login' || to.name === 'register') && auth.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
