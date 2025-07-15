<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthUser } from './AuthUserProvider'
import type { AuthUser } from './AuthUserProvider'
import { useRoute, useRouter } from 'vue-router'

const { setAuthUser } = useAuthUser()
const router = useRouter()
const route = useRoute()

onMounted(async () => {
  const res = await fetch(`${import.meta.env.VITE_API_URL}/api/v1/auth/check-session`, {
    credentials: 'include',
  })
  // if (!res.ok) router.push('/auth/login')
  if (!res.ok) {
    setAuthUser(null)

    if (route.meta.redirectOnAuthFail) {
      router.push('/auth/login')
    }
  }

  const user: AuthUser | { error: string } = await res.json()
  if ('error' in user) {
    setAuthUser(null)
    return
  }
  setAuthUser(user)
})
</script>

<template>
  <slot></slot>
</template>
