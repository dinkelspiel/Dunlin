<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthUser } from './AuthUserProvider'
import type { AuthUser } from './AuthUserProvider'
import { useRouter } from 'vue-router'

const { setUser } = useAuthUser()
const router = useRouter()

onMounted(async () => {
  const res = await fetch('http://localhost:8080/api/v1/auth/check-session', {
    credentials: 'include',
  })
  if (!res.ok) router.push('/auth/login')

  const user: AuthUser = await res.json()
  setUser(user)
})
</script>

<template>
  <slot></slot>
</template>
