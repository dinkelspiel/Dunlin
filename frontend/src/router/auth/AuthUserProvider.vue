<script setup lang="ts">
import { createAuthUserProvider } from './AuthUserProvider'
import { ref } from 'vue'
import AuthGuard from './AuthGuard.vue'

const user = ref(null)
createAuthUserProvider(user.value)
</script>

<template>
  <AuthGuard>
    <router-view></router-view>
  </AuthGuard>
</template>
