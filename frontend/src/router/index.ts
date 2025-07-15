import { createRouter, createWebHistory } from 'vue-router'
import AuthUserProvider from './auth/AuthUserProvider.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/auth/login',
      name: 'login',
      component: () => import('../views/auth/Login.vue'),
    },
    {
      path: '/setup',
      name: 'setup',
      component: () => import('../views/Setup.vue'),
    },
    {
      path: '/-',
      component: AuthUserProvider,
      children: [
        {
          path: '',
          name: 'teams',
          component: () => import('../views/Teams.vue'),
          meta: { redirectOnAuthFail: true },
        },
        {
          path: ':team',
          name: 'team',
          component: () => import('../views/Team.vue'),
          meta: { redirectOnAuthFail: true },
        },
        {
          path: ':team/:project/:filepath(.*)*',
          name: 'project',
          component: () => import('../views/Dashboard.vue'),
          meta: { redirectOnAuthFail: false },
        },
      ],
    },
  ],
})

export default router
